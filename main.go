package main

import (
	"logging-service/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// ListBuckets
// PrintBuckets
// HeadBucket (BucketExists)
// CreateBucket
// DeleteBucket

// ListObjectsV2
// PrintObjects
// DownloadObject
// ListObjectVersions
// PrintObjectVersions
// DeleteObjects
// DeleteObject
// UploadObject

func main() {
	// s3Client, err := conf.ConfigureS3()
	// if err != nil {
	// 	return
	// }

	// // ListBuckets
	// buckets, err := bucket.ListBuckets(s3Client)
	// if err != nil {
	// 	return
	// }
	// // PrintBuckets
	// bucket.PrintBuckets(buckets)

	// // BucketExists
	// _, err = bucket.BucketExists(s3Client, "logging-service-chat-sigma-api-logs")
	// if err != nil {
	// 	return
	// }

	// // ListObjectsV2
	// objects, err := object.ListObjects(*s3Client, "logging-service-chat-sigma-api-logs")
	// if err != nil {
	// 	return
	// }
	// // PrintObjects
	// object.PrintObjects(objects)

	// // DownloadObject
	// err = object.DownloadObject(*s3Client, "logging-service-chat-sigma-api-logs", "test_log.txt", "downloaded_obj.txt")
	// if err != nil {
	// 	return
	// }

	// // ListObjectVersions
	// objectVersions, err := object.ListObjectVersions(s3Client, "logging-service-chat-sigma-api-logs")
	// if err != nil {
	// 	return
	// }
	// object.PrintObjectVersions(objectVersions)

	// // UploadObject
	// err = object.UploadObject(s3Client, "logging-service-chat-sigma-api-logs", "downloaded_obj.txt", "downloaded_obj.txt")
	// if err != nil {
	// 	return
	// }

	godotenv.Load()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(os.Getenv("PORT"))
}
