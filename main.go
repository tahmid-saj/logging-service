package main

import (
	"logging-service/bucket"
	"logging-service/conf"
)

// ListBuckets
// PrintBuckets
// HeadBucket (BucketExists)
// CreateBucket
// DeleteBucket

// ListObjectsV2
// GetObject
// ListObjectVersions
// DeleteObject
// DeleteObjects
// PubObject

func main() {
	s3Client, err := conf.ConfigureS3()
	if err != nil {
		return
	}

	// ListBuckets
	buckets, err := bucket.ListBuckets(s3Client)
	if err != nil {
		return
	}
	bucket.PrintBuckets(buckets)

	// BucketExists
	_, err = bucket.BucketExists(s3Client, "logging-service-chat-sigma-api-logs")
	if err != nil {
		return
	}
}
