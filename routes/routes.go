package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// bucket

	// ListBuckets
	server.GET("/buckets", getBuckets)
	// BucketExists
	server.GET("/buckets/:bucketName", getBucketExists)
	// CreateBucket
	server.POST("/buckets", postCreateBucket)
	// DeleteBucket
	server.DELETE("/buckets/:bucketName", deleteBucket)

	// object

	// ListObjects
	server.GET("/buckets/:bucketName/objects", getObjects)
	// DownloadObject
	server.POST("/buckets/:bucketName/objects/:objectKey/download", postDownloadObject)
	// ListObjectVersions
	server.GET("/buckets/:bucketName/objects/versions", getListObjectVersions)
	// UploadObject
	server.POST("/buckets/:bucketName/objects", postUploadObject)
	// DeleteObjects
	server.POST("/buckets/:bucketName/objects/delete", deleteObjects)
	// DeleteObject
	server.DELETE("/buckets/:bucketName/objects/:objectKey/versions/:versionID", deleteObject)

	// dashboard

	// paginated dashboard view of logs
	server.GET("/dashboard-logs") // query params: start_time, end_time, hostname, path, ok
	// aggregated dashboard summary (counts / metrics) of logs
	server.GET("/dashboard-logs/aggregated") // query params: start_time, end_time, hostname, path, ok
}