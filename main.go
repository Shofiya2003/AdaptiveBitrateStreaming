package main

import (
	"abr_backend/config"
	"abr_backend/routes"
)

func main() {

	config.LoadEnv()
	r := routes.InitGin()
	config.InitCloudSession()
	config.InitDB()
	r.Run()

	// config.InitRabbitMq()

	// cloudSession := config.GetSession()
	// uploder := utils.AwsUploader{
	// 	S3Client: cloudSession.AWS,
	// }
	// // utils.UploadtoCloudStorage(uploder, "test")
	// utils.UploadtoCloudStorage(uploder, "output")

	// generate a presigned url
	// download the video
	// ffmpeg
	// upload
	// url := "https://abr-raw.s3.ap-south-1.amazonaws.com/4114797-sd_426_240_25fps.mp4?response-content-disposition=inline&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEBQaCmFwLXNvdXRoLTEiRzBFAiBw%2BLLlj9RQgLcMPJdRPNu7RKXCTgYoSIl3CLNWCLXd1AIhANfV1K4zOuvq76nHmMwDJbrWP0sJgiOg8PZTCHGju%2FaNKscDCE0QABoMMTIzMjE1MzM1MTQ5IgxuMWoVtEZu32tt6%2BgqpAP4fZr%2B3YSy%2BVVpGoTxf%2FHx%2F2LDn%2F0Np9nMDIq%2BRhtE6V8gps4cBzjA24T7oFew2p9YP%2Fleb8XN8CiCRtFy0aX%2BtEY0QzIyV%2B8GT4jqLarev0NoqwKuk%2FBifFpUgIh527kAeV16gWAlLN%2B5frRqZ0yuX7221hTvS2RBI4vbXg4Ja6pV3fqR70HXoBraDeLcpJKR9L8wAM%2FXIDbwXKaH5U9sfdkb5PDTDapwe52PSYofHsqe%2FBR3wTJePdH%2Fb6PLHCL76sMONLavRECYdl%2BNtKt5%2Foz4dR5VSoODZVzPFuZnyP%2F5Aq%2F5mi3n542KI97AwGBYt4eAXbXJ2obHkeE%2Bv88d6%2FAv7FQEk1ybriRjFLKXWLcWC4%2B6MIOMoj2RVdIUJUvzI05Aw97GtZ23ELlQGFJGBfNeq%2BSkeqqt9R7y0y8O8AgrWgrzKm0FzeJMbsoKiXRh9%2FUC5pG6bFCvOM0wQ5Uk0YtwIbceMfllFl20pCXKPKp8vW8q0JvRxF3nCZaNk%2BplFGEpXJYXdGXqKEvO4V9MbHo7akSAaDsOM%2FMX8uyxndH%2FYikwx5f4vQY65AIpWNWSCPw0o1yowq5wX9nSAzbAdQppO%2FrdLuAfOuq%2Buk6g8fPG%2FA0QBbrjuTqdgy7EUMX2GbvWuP21cZNU%2BKH%2FhS%2F%2Fzr1oQ658uh%2BURoIQuvtKv2KPKFrlQpS590tWTR0%2B5KQ1vH%2FCICYhcnh1GaxYABnLDUjtzqtvWMJy%2FwmK9S5yGN8RIdGHPZGu0QTcombCkGnHp7AtmI%2Feodhxt4%2FCUwqhu8rMfIPjV%2FC4nH8QdIbHRGcjrQ3kRppS81e9A2gV%2BawsGdGQdrBoWWQzONWaYG0ITtYfdRfa4e0j9PexPmIK95mhNmayPBejOGt5bJ1VMwIuUCvEZNZg%2FP3oId5ljcxSQSjl6Iw8wdMhgZEEOUbOmWGxz7DLxjxWo5FYomUBjcKrHg%2FQf6W1nS0mD1lTAJMvslL5uOy6yAkEkMNksYriiIEnQg%2FePjbQEgQXHVKrqqHvzvIbbGYudAGFjAs1xMN10A%3D%3D&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIARZMBUZ3WR2SIXYXE%2F20250225%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20250225T194245Z&X-Amz-Expires=10800&X-Amz-SignedHeaders=host&X-Amz-Signature=e86c972a307d865d3be5d7458c6aa22a152344716883a23ea6f1b764308f2f30"
	// utils.Fetch(url, "tmp/raw.mp4")
}
