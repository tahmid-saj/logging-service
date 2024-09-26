package routes

import (
	"logging-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// dashboard
func getDashboardLogs(context *gin.Context) {
	skip := context.Query("skip")
	limit := context.Query("limit")
	startTime := context.Query("startTime")
	endTime := context.Query("endTime")
	hostName := context.Query("hostName")
	method := context.Query("method")
	path := context.Query("path")
	ok := context.Query("ok")

	var skipParsed int64
	var limitParsed int64
	var err error
	if skip != "" {
		skipParsed, err = strconv.ParseInt(skip, 10, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request query"})
			return
		}
	} else if skip == "" { skipParsed = -1 }

	if limit != "" {
		limitParsed, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request query"})
			return
		}
	} else if limit == "" { limitParsed = -1 }

	dashboardInput := &models.DashboardInput{
		Skip: skipParsed,
		Limit: limitParsed,
		StartTime: startTime,
		EndTime: endTime,
		HostName: hostName,
		Method: method,
		Path: path,
		Ok: ok,
	}

	resGetDashboardLogs, err := models.GetDashboardLogs(dashboardInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch logs"})
		return
	}

	context.JSON(http.StatusOK, resGetDashboardLogs)
}

func getDashboardLogsAggregated(context *gin.Context) {
	startTime := context.Query("startTime")
	endTime := context.Query("endTime")
	hostName := context.Query("hostName")
	method := context.Query("method")
	path := context.Query("path")
	ok := context.Query("ok")

	dashboardInput := &models.DashboardInput{
		Skip: -1,
		Limit: -1,
		StartTime: startTime,
		EndTime: endTime,
		HostName: hostName,
		Method: method,
		Path: path,
		Ok: ok,
	}

	resGetDashboardLogsAggregated, err := models.GetDashboardLogsAggregated(dashboardInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch logs"})
		return
	}

	context.JSON(http.StatusOK, resGetDashboardLogsAggregated)
}

func getDashboardLogsByBucket(context *gin.Context) {
	bucketName := context.Param("bucketName")

	skip := context.Query("skip")
	limit := context.Query("limit")
	startTime := context.Query("startTime")
	endTime := context.Query("endTime")
	hostName := context.Query("hostName")
	method := context.Query("method")
	path := context.Query("path")
	ok := context.Query("ok")

	var skipParsed int64
	var limitParsed int64

	var err error
	if skip != "" {
		skipParsed, err = strconv.ParseInt(skip, 10, 64)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request query"})
			return
		}
	} else { skipParsed = -1 }

	if limit != "" {
		limitParsed, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request query"})
			return
		}
	} else { limitParsed = -1 }

	dashboardInputByBucket := &models.DashboardInputByBucket{
		Skip: skipParsed,
		Limit: limitParsed,
		StartTime: startTime,
		EndTime: endTime,
		HostName: hostName,
		Method: method,
		Path: path,
		Ok: ok,
		BucketName: bucketName,
	}

	resGetDashboardLogsByBucket, err := models.GetDashboardLogsByBucket(dashboardInputByBucket)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch logs"})
		return
	}

	context.JSON(http.StatusOK, resGetDashboardLogsByBucket)
}

func getDashboardLogsAggregatedByBucket(context *gin.Context) {
	bucketName := context.Param("bucketName")

	startTime := context.Query("startTime")
	endTime := context.Query("endTime")
	hostName := context.Query("hostName")
	method := context.Query("method")
	path := context.Query("path")
	ok := context.Query("ok")

	dashboardInputByBucket := &models.DashboardInputByBucket{
		Skip: -1,
		Limit: -1,
		StartTime: startTime,
		EndTime: endTime,
		HostName: hostName,
		Method: method,
		Path: path,
		Ok: ok,
		BucketName: bucketName,
	}

	resGetDashboardLogsAggregatedByBucket, err := models.GetDashboardLogsAggregatedByBucket(dashboardInputByBucket)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch logs"})
		return
	}

	context.JSON(http.StatusOK, resGetDashboardLogsAggregatedByBucket)
}