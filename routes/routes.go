package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// bucket

	// ListBuckets
	server.GET("/bucket", getBuckets)
	// BucketExists
	server.GET("/bucket/:bucketName")
	// CreateBucket
	server.POST("/bucket")
	// DeleteBucket
	server.DELETE("/bucket/:bucketName")

	// object

	// ListObjects
	server.GET("/object/:bucketName")
	// DownloadObject
	server.POST("/object/download")
	// ListObjectVersions
	server.GET("/object/:bucketName/versions")
	// DeleteObjects
	server.POST("/objects")
	// DeleteObject
	server.POST("/object")

	// dashboard

	// paginated dashboard view of logs
	server.GET("/dashboard-logs") // query params: start_time, end_time, hostname, path, ok
	// aggregated dashboard summary (counts / metrics) of logs
	server.GET("/dashboard-logs/aggregated") // query params: start_time, end_time, hostname, path, ok
}