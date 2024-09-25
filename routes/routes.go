package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// bucket

	// ListBuckets
	server.GET("/bucket", getBuckets)
	// BucketExists
	server.GET("/bucket/:bucketName", getBucketExists)
	// CreateBucket
	server.POST("/bucket", postCreateBucket)
	// DeleteBucket
	server.DELETE("/bucket/:bucketName", deleteBucket)

	// object

	// ListObjects
	server.GET("/object/:bucketName")
	// DownloadObject
	server.POST("/object/:bucketName/:filename/download")
	// ListObjectVersions
	server.GET("/object/:bucketName/versions")
	// UploadObject
	server.POST("/object")
	// DeleteObjects
	server.POST("/objects/:bucketName/delete")
	// DeleteObject
	server.DELETE("/object/:bucketName/delete")

	// dashboard

	// paginated dashboard view of logs
	server.GET("/dashboard-logs") // query params: start_time, end_time, hostname, path, ok
	// aggregated dashboard summary (counts / metrics) of logs
	server.GET("/dashboard-logs/aggregated") // query params: start_time, end_time, hostname, path, ok
}