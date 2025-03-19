package upload

import (
	"abr_backend/config"
	"abr_backend/utils"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type multiPartUploadStrategy struct{}

func (m *multiPartUploadStrategy) InitializeUpload(bucket, name, fileType string) (string, error) {

	cloudSession, err := config.GetSession()

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %v", err)
	}

	S3Client := cloudSession.AWS

	multiPartUploadOutput, err := S3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(name),
		ContentType: aws.String(fileType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to get upload ID %v", err)
	}

	uploadId := *multiPartUploadOutput.UploadId

	return GetPresignedUrl(bucket, name, uploadId, 1)

}

func GetPresignedUrl(bucket, name, uploadId string, partNumber int32) (string, error) {
	cloudSession, err := config.GetSession()

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %v", err)
	}

	S3Client := cloudSession.AWS
	presignClient := s3.NewPresignClient(S3Client)
	presignedURL, err := presignClient.PresignUploadPart(context.TODO(), &s3.UploadPartInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(name),
		UploadId:   aws.String(uploadId),
		PartNumber: aws.Int32(partNumber),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %v", err)
	}

	return presignedURL.URL, nil
}

func CompleteUpload(bucket, key, uploadId string) error {

	cloudSession, err := config.GetSession()

	if err != nil {
		return fmt.Errorf("failed to create presigned URL: %v", err)
	}

	s3Client := cloudSession.AWS

	parts, err := utils.ListMultipartUploadParts(bucket, key, uploadId)

	if err != nil {
		return fmt.Errorf("failed to list multipart upload parts: %v", err)
	}

	_, err = s3Client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to complete multipart upload: %v", err)
	}

	return nil
}
