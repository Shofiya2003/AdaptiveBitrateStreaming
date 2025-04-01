package utils

import (
	"abr_backend/config"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

// channel to extract files from the folder
type fileWalk chan string

type Uploader interface {
	Upload(walker fileWalk, key_prefix string) error
}

func UploadtoCloudStorage(uploader Uploader, path, key_prefix string) error {
	fw := make(fileWalk)

	go func() {
		if err := filepath.Walk(path, fw.WalkFunc); err != nil {
			fmt.Println("Error walking directory:", err)
		}

		close(fw)

	}()

	return uploader.Upload(fw, key_prefix)
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
	S3Client *s3.Client
}

// Create CloudSession interface
// Upload API -> pushes the task to upload into rabbitmq
// rabbitmq task -> creates Uploader object, calls for upload

func (a AwsUploader) Upload(walker fileWalk, key_prefix string) error {

	s3Client := a.S3Client

	uploader := manager.NewUploader(s3Client)

	bucket := config.ConfigValues[config.AWS_S3_TRANSCODED_BUCKET_NAME]

	for pathName := range walker {
		fmt.Printf("Uploading %s", pathName)

		filename := filepath.Base(pathName)

		file, err := os.Open(pathName)
		if err != nil {
			return fmt.Errorf("Failed opening file", pathName, err)
		}

		key := fmt.Sprintf("%s/%s", key_prefix, filename)
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   file,
			ACL:    types.ObjectCannedACL("public-read"),
		})

		if err != nil {
			file.Close()
			return fmt.Errorf("Failed to upload", pathName, err)
		}

		log.Println("Uploaded", pathName, result.Location)

		if err := file.Close(); err != nil {
			log.Println("Unable to close the file")
		}

		os.Remove(pathName)

	}

	return nil

}
