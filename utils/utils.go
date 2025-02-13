package utils

import (
	"abr-backend/config"
	"abr_backend/config"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

// channel to extract files from the folder
type fileWalk chan string

type Uploader interface {
	Upload(walker fileWalk)
}

func UploadtoCloudStorage(uploader Uploader, path string) {
	fw := make(fileWalk)

	go func() {
		if err := filepath.Walk(path, fw.WalkFunc); err != nil {
			fmt.Println("Error walking directory:", err)
		}

		close(fw)

	}()

	uploader.Upload(fw)
}

func (f fileWalk) WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		f <- path
	}

	return nil
}

type AwsUploader struct {
}

// Create CloudSession interface
// Upload API -> pushes the task to upload into rabbitmq
// rabbitmq task -> creates Uploader object, calls for upload

func (a AwsUploader) Upload(walker fileWalk) {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		fmt.Println("Error loading the config:", err)
		return
	}

	s3Client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(s3Client)

	bucket := config.ConfigValues[config.AWS_S3_RAW_BUCKET_NAME]

	for pathName := range walker {
		fmt.Printf("Uploading %s", pathName)
		filename := filepath.Base(pathName)

		file, err := os.Open(pathName)
		if err != nil {
			log.Println("Failed opening file", pathName, err)
			continue
		}

		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   file,
		})

		if err != nil {
			file.Close()
			log.Println("Failed to upload", pathName, err)
		}

		log.Println("Uploaded", pathName, result.Location)

		if err := file.Close(); err != nil {
			log.Println("Unable to close the file")
		}
	}
}
