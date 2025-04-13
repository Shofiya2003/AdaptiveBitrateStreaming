package utils

import (
	"abr_backend/config"
	"abr_backend/data"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func Fetch(url string, fileName string) (string, error) {
	log.Println("Downloading Video: ")

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(response.StatusCode)
		return "", errors.New(strconv.Itoa(response.StatusCode))
	}

	// create an empty file
	filePath := "tmp/" + fileName

	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer file.Close()

	// write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("Download Complete: ")

	return filePath, nil
}

func GetVideoUrl(bucket, key string) (string, error) {

	cloudSession, err := config.GetSession()

	if err != nil {
		return "", fmt.Errorf("failed to get session: %v", err)
	}

	s3Client := cloudSession.AWS
	presignClient := s3.NewPresignClient(s3Client)

	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(time.Minute*15))
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL: %v", err)
	}

	return presignedUrl.URL, nil
}

func TranscodeVideo(bucket, key string) error {
	clientID := strings.Split(key, "/")[0]
	err := UpdateVideoStatus(config.GetDB(), clientID, key, bucket, "transcoding")
	if err != nil {
		log.Println("Error updating video record status", err)
	}
	url, err := GetVideoUrl(bucket, key)
	fmt.Printf("File uploaded - Bucket: %s, Key: %s\n", bucket, key)
	if err != nil {
		UpdateVideoStatus(config.GetDB(), clientID, key, bucket, "failed")
		return fmt.Errorf("error getting video URL: %v", err)
	}

	filePath, err := Fetch(url, key)
	if err != nil {
		UpdateVideoStatus(config.GetDB(), clientID, key, bucket, "failed")
		return fmt.Errorf("error fetching video: %v", err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- Transcode(filePath, "tmp/vmaf", "tmp/transcoded")
	}()

	if err := <-errChan; err != nil {
		return fmt.Errorf("error transcoding video: %v", err)
	}

	cloudSession, err := config.GetSession()
	if err != nil {
		UpdateVideoStatus(config.GetDB(), clientID, key, bucket, "failed")
		return fmt.Errorf("error getting session: %v", err)
	}
	uploder := AwsUploader{
		S3Client: cloudSession.AWS,
	}

	err = UploadtoCloudStorage(uploder, "tmp/transcoded", key)
	if err != nil {
		UpdateVideoStatus(config.GetDB(), clientID, key, bucket, "failed")
		return fmt.Errorf("failed to upload video")
	}
	bucket = config.ConfigValues[config.AWS_S3_TRANSCODED_BUCKET_NAME]
	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, "ap-south-1", fmt.Sprintf("%s/index.m3u8", key))

	fmt.Println(publicURL)
	UpdateVideoStatus(config.GetDB(), clientID, key, bucket, "completed")

	return nil
}

func Transcode(inputFilePath, workingDir, outputDir string) error {

	cmd := exec.Command("/home/shofiya/abr/transcoder/smart_transcode.sh", inputFilePath, workingDir, outputDir)
	cmd.Stdout = log.Writer() // Log output to console
	cmd.Stderr = log.Writer() // Log errors to console

	if err := cmd.Run(); err != nil {
		log.Println("Error running transcoder:", err)
		return err
	}

	return nil

}

type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func GetRecords(body []byte) (string, string, error) {
	var event S3Event
	if err := json.Unmarshal(body, &event); err != nil {
		log.Println("Failed to parse SNS event:", err)

		return "", "", err
	}

	// Step 3: Check if Records Exist
	if len(event.Records) == 0 {
		log.Println("No records found in the event")
		return "", "", errors.New("no records found")
	}

	var bucket string
	var key string
	// Extract bucket and key from the first record
	if len(event.Records) > 0 {
		record := event.Records[0]
		bucket = record.S3.Bucket.Name
		key = record.S3.Object.Key
	} else {
		log.Println("No records found")
		return "", "", errors.New("no records found")
	}

	return bucket, key, nil
}

func GenerateVideoID() string {
	return uuid.New().String()
}

func GetVideos(db *sql.DB, clientID string, page int) ([]data.Video, error) {
	const pageSize = 10
	offset := (page - 1) * pageSize

	query := `
	SELECT * 
	FROM videos 
	WHERE client_id = ? 
	ORDER BY upload_time DESC
	LIMIT ? OFFSET ?;
	`

	rows, err := db.Query(query, clientID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching videos: %v", err)
	}
	defer rows.Close()

	var videos []data.Video

	for rows.Next() {
		var v data.Video
		if err := rows.Scan(&v.VideoID, &v.ClientID, &v.UploadTime, &v.Status, &v.FileKey, &v.Bucket, &v.Strategy); err != nil {
			return nil, fmt.Errorf("error scanning video: %v", err)
		}
		publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", v.Bucket, "ap-south-1", fmt.Sprintf("%s/index.m3u8", v.FileKey))
		v.Url = publicURL
		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %v", err)
	}

	return videos, nil
}

func AddVideo(db *sql.DB, videoID, clientID, status, fileKey, bucket, strategy string) error {
	query := `
	INSERT INTO videos (video_id, client_id, status, file_key, bucket, strategy)
	VALUES (?, ?, ?, ?, ?, ?);
	`

	_, err := db.Exec(query, videoID, clientID, status, fileKey, bucket, strategy)
	if err != nil {
		log.Printf("Error inserting video: %v", err)
		return err
	}

	log.Println("Video added successfully!")
	return nil
}

func UpdateVideoStatus(db *sql.DB, clientID, fileKey, bucket, newStatus string) error {
	query := `
	UPDATE videos 
	SET status = ?
	WHERE client_id = ? AND file_key = ?;
	`

	result, err := db.Exec(query, newStatus, clientID, fileKey)
	if err != nil {
		log.Printf("Error updating video status: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no video found with client_id=%s, file_key=%s, bucket=%s", clientID, fileKey, bucket)
	}

	log.Printf("Successfully updated video status to %s for client_id=%s, file_key=%s", newStatus, clientID, fileKey)
	return nil
}
