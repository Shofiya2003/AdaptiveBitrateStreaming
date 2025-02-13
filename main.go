package main

import "abr_backend/utils"

func main() {
	uploder := utils.AwsUploader{}
	utils.UploadtoCloudStorage(uploder, "test")
}
