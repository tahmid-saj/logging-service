package routes

import (
	"logging-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// bucket
func getBuckets(context *gin.Context) {
	buckets, err := models.GetListBuckets()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch buckets"})
		return
	}

	context.JSON(http.StatusOK, buckets)
}

func getBucketExists(context *gin.Context) {
	
}