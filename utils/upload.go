package utils

import (
	"abr_backend/config"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

func ListMultipartUploadParts(bucket, key, uploadId string) ([]types.CompletedPart, error) {

	cloudSession, err := config.GetSession()

	if err != nil {
		return nil, fmt.Errorf("failed to create presigned URL: %v", err)
	}

	s3Client := cloudSession.AWS
	listPartsOutput, err := s3Client.ListParts(context.TODO(), &s3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	})

	var completedParts []types.CompletedPart

	for _, part := range listPartsOutput.Parts {
		completedParts = append(completedParts, types.CompletedPart{
			ETag:       part.ETag,
			PartNumber: part.PartNumber,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %v", err)
	}

	return completedParts, nil

}
