package upload

import (
	"abr_backend/config"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type SinglePartUploadStrategy struct{}

func (s *SinglePartUploadStrategy) InitializeUpload(bucket, name, fileType string) (string, string, error) {

	cloudSession, err := config.GetSession()

	if err != nil {
		return "", "", fmt.Errorf("failed to create presigned URL: %v", err)
	}

	S3Client := cloudSession.AWS
	presignClient := s3.NewPresignClient(S3Client)
	presignedURL, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(name),
		ContentType: aws.String(fileType),
	}, s3.WithPresignExpires(1000*time.Minute))

	if err != nil {
		return "", "", fmt.Errorf("failed to create presigned URL: %v", err)
	}

	return presignedURL.URL, "", nil

}
