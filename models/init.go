package models

import (
	"logging-service/conf"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Response struct {
	Ok       bool        `json:"ok"`
	Response interface{} `json:"response"`
}

func InitS3() *s3.Client {
	var s3Client, err = conf.ConfigureS3()
	if err != nil {
		os.Exit(1)
	}

	return s3Client
}