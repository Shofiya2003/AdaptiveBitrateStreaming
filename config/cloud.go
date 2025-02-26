package config

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CloudSession struct {
	AWS *s3.Client
}

var cloudSession *CloudSession

func InitCloudSession() {
	cloudSession = new(CloudSession)
	awsClient, err := getAWSSession()
	if err != nil {
		log.Panicln("Failed to initialize aws client")
	}

	cloudSession.AWS = awsClient
}

func GetSession() *CloudSession {

	return cloudSession

}

func getAWSSession() (*s3.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		fmt.Println("Error loading the config:", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
