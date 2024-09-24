package models

import (
	"logging-service/bucket"
	"logging-service/conf"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	BucketName   string `json:"bucketName"`
	CreationDate string `json:"creationDate"`
}

func InitS3() *s3.Client {
	var s3Client, err = conf.ConfigureS3()
	if err != nil {
		os.Exit(1)
	}

	return s3Client
}

// bucket
func GetListBuckets() ([]*Bucket, error) {
	s3Client := InitS3()

	buckets, err := bucket.ListBuckets(s3Client)
	if err != nil {
		return nil, err
	}

	var resBuckets []*Bucket
	for _, bucket := range buckets {
		resBuckets = append(resBuckets, &Bucket{
			BucketName: *bucket.Name,
			CreationDate: bucket.CreationDate.GoString(),
		})
	}

	return resBuckets, nil
}