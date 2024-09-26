package models

import (
	"fmt"
	"log"
	"logging-service/bucket"
	"logging-service/object"
	"logging-service/utils"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

type LogsPaginated struct {
	Logs []*Log
	PaginationCursor PaginationCursor
}

type PaginationCursor struct {
	Next *NextCursor
	Previous *PreviousCursor
}

type NextCursor struct {
	NextSkip int `json:"nextSkip"`
	NextLimit int `json:"nextLimit"`
	NextURL string `json:"nextURL"`
}

type PreviousCursor struct {
	PreviousSkip int `json:"previousSkip"`
	PreviousLimit int `json:"previousLimit"`
	PreviousURL string `json:"previousURL"`
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
	Skip int64 `json:"skip"`
	Limit int64 `json:"limit"`
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
	logs, paginationCursor, err := getAllLogs(dashboardInput)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	// error checking paginationCursor
	var resPaginationCursor PaginationCursor
	if paginationCursor == nil { 
		resPaginationCursor = PaginationCursor{} 
	} else { 
		resPaginationCursor = *paginationCursor 
	}

	return &Response{
		Ok: true,
		Response: &LogsPaginated{
			Logs: logs,
			PaginationCursor: resPaginationCursor,
		},
	}, nil
}

func GetDashboardLogsAggregated(dashboardInput *DashboardInput) (*Response, error) {
	logs, _, err := getAllLogs(dashboardInput)
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

func GetDashboardLogsByBucket(dashboardInputByBucket *DashboardInputByBucket) (*Response, error) {
	s3Client := InitS3()

	dashboardInput := &DashboardInput{
		Skip: int64(dashboardInputByBucket.Skip),
		Limit: int64(dashboardInputByBucket.Limit),
		StartTime: dashboardInputByBucket.StartTime,
		EndTime: dashboardInputByBucket.EndTime,
		HostName: dashboardInputByBucket.HostName,
		Method: dashboardInputByBucket.Method,
		Path: dashboardInputByBucket.Path,
		Ok: dashboardInputByBucket.Ok,
	}

	objects, err := object.ListObjects(s3Client, dashboardInputByBucket.BucketName)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	logs, err := getAllLogsInBucket(s3Client, dashboardInput, dashboardInputByBucket.BucketName, objects)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, nil
	}

	paginatedLogs, paginationCursor, err := getLogsPaginated(logs, dashboardInput)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, nil
	}

	// error checking paginationCursor
	var resPaginationCursor PaginationCursor
	if paginationCursor == nil { 
		resPaginationCursor = PaginationCursor{} 
	} else { 
		resPaginationCursor = *paginationCursor 
	}

	return &Response{
		Ok: true,
		Response: &LogsPaginated{
			Logs: paginatedLogs,
			PaginationCursor: resPaginationCursor,
		},
	}, nil
}

// func GetDashboardLogsAggregatedByBucket(dashboardInputByBucket *DashboardInputByBucket) (*Response, error) {

// }

func getAllLogs(dashboardInput *DashboardInput) ([]*Log, *PaginationCursor, error) {
	s3Client := InitS3()

	buckets, err := bucket.ListBuckets(s3Client)
	if err != nil {
		log.Print(err.Error())
		return nil, nil, err
	}

	allBucketsLogs := []*Log{}

	for _, bucket := range buckets {
		objects, err := object.ListObjects(s3Client, *bucket.Name)
		if err != nil {
			log.Print(err.Error())
			return nil, nil, err
		}

		currentBucketLogs, err := getAllLogsInBucket(s3Client, dashboardInput, *bucket.Name, objects)
		if err != nil {
			return nil, nil, err
		}

		allBucketsLogs = append(allBucketsLogs, currentBucketLogs...)

		// for _, obj := range objects {
		// 	// read the object rows
		// 	resRows, err := object.ReadObject(s3Client, *bucket.Name, *obj.Key)
		// 	if err != nil {
		// 		log.Print(err.Error())
		// 		return nil, nil, err
		// 	}

		// 	for rowIndex, row := range resRows {
		// 		if rowIndex == 0 { continue }

		// 		timestamp := row[0]
		// 		hostName := row[1]
		// 		method := row[2]
		// 		path := row[3]

		// 		okParsed, err := strconv.ParseBool(row[4])
		// 		if err != nil {
		// 			return nil, nil, err
		// 		}
				
		// 		errorParsed, err := strconv.ParseBool(row[5])
		// 		if err != nil {
		// 			return nil, nil, err
		// 		}

		// 		description := row[6]

		// 		isTimeBetween, err := utils.IsTimeBetween(timestamp, dashboardInput.StartTime, dashboardInput.EndTime)
		// 		if err != nil {
		// 			log.Print(err.Error())
		// 			return nil, nil, err
		// 		}

		// 		if !isTimeBetween { continue }
		// 		if dashboardInput.HostName != "" && hostName != dashboardInput.HostName { continue }
		// 		if dashboardInput.Method != "" && method != dashboardInput.Method { continue }
		// 		if dashboardInput.Path != "" && path != dashboardInput.Path { continue }
		// 		if dashboardInput.Ok != "" && row[4] != dashboardInput.Ok { continue }
				
		// 		logs = append(logs, Log{
		// 			Timestamp: timestamp,
		// 			HostName: hostName,
		// 			Method: method,
		// 			Path: path,
		// 			Ok: okParsed,
		// 			Error: errorParsed,
		// 			Description: description,
		// 		})
		// 	}
		// }
	}

	// // pagination
	// if dashboardInput.Skip == -1 || dashboardInput.Limit == -1 {
	// 	return logs, nil, nil
	// }

	// var startIndex int64
	// var endIndex int64
	
	// if dashboardInput.Skip < int64(len(logs)) && dashboardInput.Skip >= 0 {
	// 	startIndex = dashboardInput.Skip
	// } else {
	// 	startIndex = 0
	// }

	// if dashboardInput.Skip + dashboardInput.Limit < int64(len(logs)) {
	// 	endIndex = dashboardInput.Skip + dashboardInput.Limit
	// } else {
	// 	endIndex = int64(len(logs))
	// }
	// paginatedLogs := logs[startIndex: endIndex]

	// // next and previous cursor
	// var next *NextCursor
	// var previous *PreviousCursor
	
	// if dashboardInput.Skip + dashboardInput.Limit < int64(len(logs)) {
	// 	next = &NextCursor{
	// 		NextSkip: int(dashboardInput.Skip) + int(dashboardInput.Limit),
	// 		NextLimit: int(dashboardInput.Limit),
	// 		NextURL: fmt.Sprintf("%v%v?skip=%v&limit=%v&startTime=%v&endTime=%v&hostName=%v&method=%v&path=%v&ok=%v", 
	// 			os.Getenv("ROOT_API_URL"), 
	// 			os.Getenv("DASHBOARD_LOGS_PATH"),
	// 			int(dashboardInput.Skip) + int(dashboardInput.Limit),
	// 			int(dashboardInput.Limit),
	// 			utils.ReplaceSpacesURL(dashboardInput.StartTime),
	// 			utils.ReplaceSpacesURL(dashboardInput.EndTime),
	// 			dashboardInput.HostName,
	// 			dashboardInput.Method,
	// 			dashboardInput.Path,
	// 			dashboardInput.Ok,
	// 		),
	// 	}
	// } else {
	// 	next = nil
	// }

	// if dashboardInput.Skip - dashboardInput.Limit >= 0 {
	// 	previous = &PreviousCursor{
	// 		PreviousSkip: int(dashboardInput.Skip) - int(dashboardInput.Limit),
	// 		PreviousLimit: int(dashboardInput.Limit),
	// 		PreviousURL: fmt.Sprintf("%v%v?skip=%v&limit=%v&startTime=%v&endTime=%v&hostName=%v&method=%v&path=%v&ok=%v", 
	// 			os.Getenv("ROOT_API_URL"), 
	// 			os.Getenv("DASHBOARD_LOGS_PATH"),
	// 			int(dashboardInput.Skip) - int(dashboardInput.Limit),
	// 			int(dashboardInput.Limit),
	// 			utils.ReplaceSpacesURL(dashboardInput.StartTime),
	// 			utils.ReplaceSpacesURL(dashboardInput.EndTime),
	// 			dashboardInput.HostName,
	// 			dashboardInput.Method,
	// 			dashboardInput.Path,
	// 			dashboardInput.Ok,
	// 		),
	// 	}
	// } else {
	// 	previous = nil
	// }

	paginatedLogs, paginationCursor, err := getLogsPaginated(allBucketsLogs, dashboardInput)
	if err != nil {
		return nil, nil, err
	}

	return paginatedLogs, paginationCursor, nil
}

func getAllLogsInBucket(s3Client *s3.Client, dashboardInput *DashboardInput, bucketName string, objects []types.Object) ([]*Log, error) {
	logs := []*Log{}

	for _, obj := range objects {
		// read the object rows
		resRows, err := object.ReadObject(s3Client, bucketName, *obj.Key)
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

			// filtering
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
			
			logs = append(logs, &Log{
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

	return logs, nil
}

func getLogsPaginated(logs []*Log, dashboardInput *DashboardInput) ([]*Log, *PaginationCursor, error) {
	// pagination
	if dashboardInput.Skip == -1 || dashboardInput.Limit == -1 {
		return logs, nil, nil
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
	paginatedLogs := logs[startIndex: endIndex]

	// next and previous cursor
	var next *NextCursor
	var previous *PreviousCursor
	
	if dashboardInput.Skip + dashboardInput.Limit < int64(len(logs)) {
		next = &NextCursor{
			NextSkip: int(dashboardInput.Skip) + int(dashboardInput.Limit),
			NextLimit: int(dashboardInput.Limit),
			NextURL: fmt.Sprintf("%v%v?skip=%v&limit=%v&startTime=%v&endTime=%v&hostName=%v&method=%v&path=%v&ok=%v", 
				os.Getenv("ROOT_API_URL"), 
				os.Getenv("DASHBOARD_LOGS_PATH"),
				int(dashboardInput.Skip) + int(dashboardInput.Limit),
				int(dashboardInput.Limit),
				utils.ReplaceSpacesURL(dashboardInput.StartTime),
				utils.ReplaceSpacesURL(dashboardInput.EndTime),
				dashboardInput.HostName,
				dashboardInput.Method,
				dashboardInput.Path,
				dashboardInput.Ok,
			),
		}
	} else {
		next = nil
	}

	if dashboardInput.Skip - dashboardInput.Limit >= 0 {
		previous = &PreviousCursor{
			PreviousSkip: int(dashboardInput.Skip) - int(dashboardInput.Limit),
			PreviousLimit: int(dashboardInput.Limit),
			PreviousURL: fmt.Sprintf("%v%v?skip=%v&limit=%v&startTime=%v&endTime=%v&hostName=%v&method=%v&path=%v&ok=%v", 
				os.Getenv("ROOT_API_URL"), 
				os.Getenv("DASHBOARD_LOGS_PATH"),
				int(dashboardInput.Skip) - int(dashboardInput.Limit),
				int(dashboardInput.Limit),
				utils.ReplaceSpacesURL(dashboardInput.StartTime),
				utils.ReplaceSpacesURL(dashboardInput.EndTime),
				dashboardInput.HostName,
				dashboardInput.Method,
				dashboardInput.Path,
				dashboardInput.Ok,
			),
		}
	} else {
		previous = nil
	}

	return paginatedLogs, &PaginationCursor{
		Next: next,
		Previous: previous,
	}, nil
}

func getLogsAggregated(logs []*Log) (LogAggregated, error) {
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