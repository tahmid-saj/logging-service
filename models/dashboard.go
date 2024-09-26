package models

type Log struct {
	Timestamp string `json:"timestamp"`
	HostName string `json:"hostName"`
	Method string	`json:"method"`
	Path string `json:"path"`
	Ok bool `json:"ok"`
	Error bool `json:"error"`
	Description string `json:"description"`
}

type DashboardInput struct {
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	HostName string `json:"hostName"`
	Path string `json:"path"`
	Ok bool `json:"ok"`
}

type DashboardInputByBucket struct {
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	HostName string `json:"hostName"`
	Path string `json:"path"`
	Ok bool `json:"ok"`
	BucketName string `json:"bucketName"`
}

// dashboard
func GetDashboardLogs(dashboardInput *DashboardInput) (*Response, error) {
	
}

func GetDashboardLogsAggregated(dashboardInput *DashboardInput) (*Response, error) {

}

func GetDashboardLogsByBucket(dashboardInputByBucket *DashboardInputByBucket) (*Response, error) {

}

func GetDashboardLogsAggregatedByBucket(dashboardInputByBucket *DashboardInputByBucket) (*Response, error) {

}