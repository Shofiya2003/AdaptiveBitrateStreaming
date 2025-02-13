package main

import (
	"abr_backend/config"
	"abr_backend/utils"
)

func main() {
	config.LoadEnv()
	uploder := utils.AwsUploader{}
	utils.UploadtoCloudStorage(uploder, "test")
}
