package routes

import (
	"logging-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// bucket
func getBuckets(context *gin.Context) {
	resBuckets, err := models.GetListBuckets()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch buckets"})
		return
	}

	context.JSON(http.StatusOK, resBuckets)
}

func getBucketExists(context *gin.Context) {
	bucketName := context.Param("bucketName")
	
	resBucketExists, err := models.GetBucketExists(bucketName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch bucket"})
		return
	}

	context.JSON(http.StatusOK, resBucketExists)
}

func postCreateBucket(context *gin.Context) {
	var bucketInput models.BucketInput

	err := context.ShouldBindJSON(&bucketInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	resBucketCreated, err := models.PostCreateBucket(bucketInput.BucketName, bucketInput.Region)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create bucket"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "bucket created", "bucketResponse": resBucketCreated})
}

func deleteBucket(context *gin.Context) {
	bucketName := context.Param("bucketName")

	resBucketDeleted, err := models.DeleteBucket(bucketName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete bucket"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "bucket deleted", "bucketResponse": resBucketDeleted})
}