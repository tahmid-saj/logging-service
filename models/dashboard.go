package models

import (
	"fmt"
	"log"
	"logging-service/bucket"
	"logging-service/object"
	"logging-service/utils"
	"strconv"
)

type Log struct {
	Timestamp   string `json:"timestamp"`
	HostName    string `json:"hostName"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Ok          bool   `json:"ok"`
	Error       bool   `json:"error"`
	Description string `json:"description"`
}

type LogAggregated struct {
	Requests int `json:"requests"`
	Ok int `json:"ok"`
	Errors int `json:"errors"`
	RequestsView RequestView `json:"requestsView"`
	AggregatedTable []string `json:"aggregatedTable"`
}

// hostname maps to method, then maps to path, then maps to the number of requests (int)
type RequestView map[string]map[string]map[string]RequestTotals

type RequestTotals struct {
	Requests int
	Ok int
	Errors int
}

type DashboardInput struct {
	Skip int64 `json:"skip"`
	Limit int64 `json:"limit"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	HostName  string `json:"hostName"`
	Method 		 string `json:"method"`
	Path      string `json:"path"`
	Ok        string   `json:"ok"`
}

type DashboardInputByBucket struct {
	Skip int `json:"skip"`
	Limit int `json:"limit"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	HostName   string `json:"hostName"`
	Method 		 string `json:"method"`
	Path       string `json:"path"`
	Ok         string   `json:"ok"`
	BucketName string `json:"bucketName"`
}

// dashboard
func GetDashboardLogs(dashboardInput *DashboardInput) (*Response, error) {
	logs, err := getAllLogs(dashboardInput)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	return &Response{
		Ok: true,
		Response: logs,
	}, nil
}

func GetDashboardLogsAggregated(dashboardInput *DashboardInput) (*Response, error) {
	logs, err := getAllLogs(dashboardInput)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	logsAggregated, err := getLogsAggregated(logs)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, nil
	}

	return &Response{
		Ok: true,
		Response: logsAggregated,
	}, nil
}

// func GetDashboardLogsByBucket(dashboardInputByBucket *DashboardInputByBucket) (*Response, error) {

// }

// func GetDashboardLogsAggregatedByBucket(dashboardInputByBucket *DashboardInputByBucket) (*Response, error) {

// }

func getAllLogs(dashboardInput *DashboardInput) ([]Log, error) {
	s3Client := InitS3()

	buckets, err := bucket.ListBuckets(s3Client)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	var logs = []Log{}

	// filtering
	for _, bucket := range buckets {
		objects, err := object.ListObjects(s3Client, *bucket.Name)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}

		for _, obj := range objects {
			// read the object rows
			resRows, err := object.ReadObject(s3Client, *bucket.Name, *obj.Key)
			if err != nil {
				log.Print(err.Error())
				return nil, err
			}

			for rowIndex, row := range resRows {
				if rowIndex == 0 { continue }

				timestamp := row[0]
				hostName := row[1]
				method := row[2]
				path := row[3]

				okParsed, err := strconv.ParseBool(row[4])
				if err != nil {
					return nil, err
				}
				
				errorParsed, err := strconv.ParseBool(row[5])
				if err != nil {
					return nil, err
				}

				description := row[6]

				isTimeBetween, err := utils.IsTimeBetween(timestamp, dashboardInput.StartTime, dashboardInput.EndTime)
				if err != nil {
					log.Print(err.Error())
					return nil, err
				}

				if !isTimeBetween { continue }
				if dashboardInput.HostName != "" && hostName != dashboardInput.HostName { continue }
				if dashboardInput.Method != "" && method != dashboardInput.Method { continue }
				if dashboardInput.Path != "" && path != dashboardInput.Path { continue }
				if dashboardInput.Ok != "" && row[4] != dashboardInput.Ok { continue }
				
				logs = append(logs, Log{
					Timestamp: timestamp,
					HostName: hostName,
					Method: method,
					Path: path,
					Ok: okParsed,
					Error: errorParsed,
					Description: description,
				})
			}
		}
	}

	// pagination
	if dashboardInput.Skip == -1 || dashboardInput.Limit == -1 {
		return logs, nil
	}

	var startIndex int64
	var endIndex int64
	
	if dashboardInput.Skip < int64(len(logs)) && dashboardInput.Skip >= 0 {
		startIndex = dashboardInput.Skip
	} else {
		startIndex = 0
	}

	if dashboardInput.Skip + dashboardInput.Limit < int64(len(logs)) {
		endIndex = dashboardInput.Skip + dashboardInput.Limit
	} else {
		endIndex = int64(len(logs))
	}
	paginatedLogs := logs[startIndex: endIndex + 1]

	return paginatedLogs, nil
}

func getLogsAggregated(logs []Log) (LogAggregated, error) {
	requests := len(logs)
	var ok int
	var errors int

	// RequestsView type
	requestsView := make(map[string]map[string]map[string]RequestTotals)

	for _, log := range logs {
		if log.Ok { ok++ }
		if log.Error { errors++ }

		// initialize the hostname map if it's nil
		if _, exists := requestsView[log.HostName]; !exists {
			requestsView[log.HostName] = make(map[string]map[string]RequestTotals)
		}

		// initialize the method map if it's nil
		if _, exists := requestsView[log.HostName][log.Method]; !exists {
			requestsView[log.HostName][log.Method] = make(map[string]RequestTotals)
		}

		// increment the count for the given path and request total
		requestTotals := requestsView[log.HostName][log.Method][log.Path]

		requestTotals.Requests++
		if log.Ok { requestTotals.Ok++ }
		if log.Error { requestTotals.Errors++ }

		requestsView[log.HostName][log.Method][log.Path] = requestTotals
	}

	aggregatedTable := getLogsAggregatedTable(requestsView)

	return LogAggregated{
		Requests: requests,
		Ok: ok,
		Errors: errors,
		RequestsView: requestsView,
		AggregatedTable: aggregatedTable,
	}, nil
}

func getLogsAggregatedTable(requestsView RequestView) (aggregatedTable []string) {
	for hostName, methodMap := range requestsView {
		for method, pathMap := range methodMap {
			for path, totals := range pathMap {
				aggregatedTable = append(aggregatedTable, fmt.Sprintf("Host: %s, Method: %s, Path: %s, Requests: %d, Ok: %d, Errors: %d", 
					hostName, method, path, totals.Requests, totals.Ok, totals.Errors))
			}
		}
	}

	return aggregatedTable
}