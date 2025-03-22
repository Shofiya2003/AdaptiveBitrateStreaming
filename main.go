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

}
