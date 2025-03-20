package upload

import "errors"

type UploadInitializer interface {
	InitializeUpload(bucket, name, fileType string) (string, string, error)
}

// Specific to Multipart Upload
type MultipartUploader interface {
	GeneratePresignedURL(bucket, name, uploadID string, partNumber int32) (string, error)
	CompleteUpload(bucket, name, uploadID string) error
}

func GetUploadStrategy(strategy string) (UploadInitializer, error) {
	switch strategy {
	case "single":
		return &SinglePartUploadStrategy{}, nil
	case "multipart":
		return &multiPartUploadStrategy{}, nil
	default:
		return nil, errors.New("invalid upload strategy")
	}

}
