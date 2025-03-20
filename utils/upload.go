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
	fmt.Printf("Listing parts for bucket: %s, key: %s, uploadId: %s\n", bucket, key, uploadId)

	listPartsOutput, err := s3Client.ListParts(context.TODO(), &s3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %v", err)
	}

	if listPartsOutput == nil {
		return nil, fmt.Errorf("list parts output is nil")
	}

	fmt.Printf("Found %d parts in the multipart upload\n", len(listPartsOutput.Parts))

	var completedParts []types.CompletedPart
	for _, part := range listPartsOutput.Parts {
		fmt.Printf("Part %d: ETag %s\n", part.PartNumber, *part.ETag)
		completedParts = append(completedParts, types.CompletedPart{
			ETag:       part.ETag,
			PartNumber: part.PartNumber,
		})
	}

	return completedParts, nil
}
