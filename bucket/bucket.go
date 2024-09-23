package bucket

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ListBuckets(s3Client *s3.Client) ([]types.Bucket, error) {
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	var buckets []types.Bucket
	if err != nil {
		log.Printf("Couldn't list buckets for this account. %v.\n", err)
	} else {
		buckets = result.Buckets
	}

	return buckets, err
}

func PrintBuckets(buckets []types.Bucket) {
	for _, bucket := range buckets {
		fmt.Printf("\t%v\n", *bucket.Name)
	}
}

func BucketExists(s3Client *s3.Client, bucketName string) (bool, error) {
	_, err := s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	exists := true
	if err != nil {
		var apiError error
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is unavailable.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to %v or here is what happened: %v.\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists.\n", bucketName)
	}

	return exists, err
}

func CreateBucket(s3Client *s3.Client, bucketName string, region string) error {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		log.Printf("Couldn't create bucket %v in region %v. %v.\n", bucketName, region, err)
	}

	return err
}

func DeleteBucket(s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Couldn't delete bucket %v. %v.\n", bucketName, err)
	}

	return err
}