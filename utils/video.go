package utils

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Fetch(url string, fileName string) error {
	log.Println("Downloading Video: ")

	// newFileName, err := GetDownloadFilePathName(fileName)

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(response.StatusCode)
		return errors.New(strconv.Itoa(response.StatusCode))
	}

	// create an empty file
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return err
	}

	defer file.Close()

	// write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Download Complete: ")

	return nil
}
