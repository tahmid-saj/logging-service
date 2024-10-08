package routes

import (
	"logging-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// object
func getObjects(context *gin.Context) {
	bucketName := context.Param("bucketName")

	resObjects, err := models.GetListObjects(bucketName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch objects"})
		return
	}

	context.JSON(http.StatusOK, resObjects)
}

func postDownloadObject(context *gin.Context) {
	bucketName := context.Param("bucketName")
	objectKey := context.Param("objectKey")

	var objectDownloadInput models.ObjectDownloadInput

	err := context.ShouldBindJSON(&objectDownloadInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	resDownloadObject, err := models.PostDownloadObject(bucketName, objectKey, objectDownloadInput.FileName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not download object"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "object downloaded", "objectDownloadResponse": resDownloadObject})
}

func getListObjectVersions(context *gin.Context) {
	bucketName := context.Param("bucketName")
	
	resObjectVersions, err := models.GetListObjectVersions(bucketName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch objects"})
		return
	}

	context.JSON(http.StatusOK, resObjectVersions)
}

func postUploadObject(context *gin.Context) {
	bucketName := context.Param("bucketName")

	var objectCreateInput models.ObjectCreateInput

	err := context.ShouldBindJSON(&objectCreateInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	resUploadObject, err := models.PostUploadObject(bucketName, objectCreateInput.ObjectKey, objectCreateInput.FileName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not upload object"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "object uploaded", "objectResponse": resUploadObject})
}

func deleteObjects(context *gin.Context) {
	bucketName := context.Param("bucketName")

	var objectVersionsInput models.ObjectVersionInput

	err := context.ShouldBindJSON(&objectVersionsInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body"})
		return
	}

	resDeleteObjects, err := models.DeleteObjects(bucketName, objectVersionsInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete objects"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "objects deleted", "objectResponse": resDeleteObjects})
}

func deleteObject(context *gin.Context) {
	bucketName := context.Param("bucketName")
	objectKey := context.Param("objectKey")
	versionID := context.Param("versionID")

	resDeleteObject, err := models.DeleteObject(bucketName, objectKey, versionID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete object"})
		return
	}
	
	switch resDeleteObject.Response.(type) {
	case bool:
		resDelete := resDeleteObject.Response
		if resDelete == false {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete object"})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{"message": "object deleted", "objectResponse": resDeleteObject})
}