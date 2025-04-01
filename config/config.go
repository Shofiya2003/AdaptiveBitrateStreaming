package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type CONFIG_KEY int

const (
	PORT CONFIG_KEY = iota
	AWS_S3_RAW_BUCKET_NAME
	AWS_S3_TRANSCODED_BUCKET_NAME
	JWT_SECRET
)

var configVarNames = map[CONFIG_KEY]string{
	PORT:                          "PORT",
	AWS_S3_RAW_BUCKET_NAME:        "AWS_S3_RAW_BUCKET_NAME",
	AWS_S3_TRANSCODED_BUCKET_NAME: "AWS_S3_TRANSCODED_BUCKET_NAME",
	JWT_SECRET:                    "JWT_SECRET",
}

var ConfigValues = map[CONFIG_KEY]string{
	PORT:                          "8000",
	AWS_S3_RAW_BUCKET_NAME:        "",
	AWS_S3_TRANSCODED_BUCKET_NAME: "",
	JWT_SECRET:                    "",
}

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}

	for key, val := range configVarNames {
		envValue, found := os.LookupEnv(val)
		if !found {
			envValue = val
		}
		ConfigValues[key] = envValue
	}
}
