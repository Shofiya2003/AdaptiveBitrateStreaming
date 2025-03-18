package upload

import "errors"

type UploadStrategy interface {
	InitializeUpload(bucket, name, fileType string) (string, error)
}

func GetUploadStrategy(strategy string) (UploadStrategy, error) {
	switch strategy {
	case "single":
		return &SinglePartUploadStrategy{}, nil
	default:
		return nil, errors.New("invalid upload strategy")
	}

}
