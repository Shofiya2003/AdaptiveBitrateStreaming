package constants

import "time"

const S3_BUCKET_NAME = "zstream-bucket"
const S3_REGION = "us-east-1"
const AWS_ENDPOINT = "http://localhost:4566"
const PRESIGNED_URL_EXPIRATION = 60 * time.Minute
const OUTPUT_FILE_PATH_PREFIX = "output"
