package models

import (
	"logging-service/bucket"
	"time"
)

type Bucket struct {
	BucketName   string `json:"bucketName"`
	CreationDate string `json:"creationDate"`
}

type BucketInput struct {
	BucketName string `json:"bucketName"`
	Region string `json:"region"`
}

// bucket
func GetListBuckets() (*Response, error) {
	s3Client := InitS3()

	buckets, err := bucket.ListBuckets(s3Client)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	var resBuckets []*Bucket
	for _, bucket := range buckets {
		resBuckets = append(resBuckets, &Bucket{
			BucketName: *bucket.Name,
			CreationDate: bucket.CreationDate.GoString(),
		})
	}

	return &Response{
		Ok: true,
		Response: resBuckets,
	}, nil
}

func GetBucketExists(bucketName string) (*Response, error) {
	s3Client := InitS3()

	bucketExists, err := bucket.BucketExists(s3Client, bucketName)
	if err != nil {
		return &Response{Ok: false, Response: nil}, err
	}

	return &Response{Ok: true, Response: bucketExists}, nil
}

func PostCreateBucket(bucketName string, region string) (*Response, error) {
	s3Client := InitS3()

	err := bucket.CreateBucket(s3Client, bucketName, region)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	bucket := &Bucket{
		BucketName: bucketName,
		CreationDate: time.Now().String(),
	}

	return &Response{
		Ok: true,
		Response: bucket,
	}, nil
}

func DeleteBucket(bucketName string) (*Response, error) {
	s3Client := InitS3()

	err := bucket.DeleteBucket(s3Client, bucketName)
	if err != nil {
		return &Response{Ok: false, Response: false}, err
	}

	return &Response{Ok: true, Response: true}, nil
}