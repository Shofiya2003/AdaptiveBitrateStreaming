package utils

import (
	"abr_backend/config"
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
	url, err := GetVideoUrl(bucket, key)
	fmt.Printf("File uploaded - Bucket: %s, Key: %s\n", bucket, key)
	if err != nil {
		return fmt.Errorf("error getting video URL: %v", err)
	}

	filePath, err := Fetch(url, key)
	if err != nil {
		return fmt.Errorf("error fetching video: %v", err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- Transcode(filePath, "tmp/transcoded")
	}()

	if err := <-errChan; err != nil {
		return fmt.Errorf("error transcoding video: %v", err)
	}

	return nil
}

func Transcode(inputFilePath, outputDir string) error {

	cmd := exec.Command("/home/shofiya/abr/utils/transcode.sh", inputFilePath, outputDir)
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
