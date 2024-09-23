package conf

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Conf struct {
	S3Client *s3.Client
}

func ConfigureS3() (*s3.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration")
		fmt.Println(err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)
	
	return s3Client, nil
}