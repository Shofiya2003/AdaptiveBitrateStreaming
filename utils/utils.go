package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func (a AwsUploader) Upload(walker fileWalk) {

	region := "ap-south-1"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	bucket := "abr-raw"

	uploader := s3manager.NewUploader(sess)

	for pathName := range walker {
		fmt.Printf("Uploading %s", pathName)
		filename := filepath.Base(pathName)

		file, err := os.Open(pathName)
		if err != nil {
			log.Println("Failed opening file", pathName, err)
			continue
		}

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(path.Join(filename)),
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
