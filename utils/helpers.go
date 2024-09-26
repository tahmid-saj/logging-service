package utils

import (
	"fmt"
	"time"
)

// check if a time string is between startTime and endTime
func IsTimeBetween(targetTimeStr, startTimeStr, endTimeStr string) (bool, error) {
	// fefine the layout for time parsing
	const timeLayout = "2006-01-02 15:04:05"
	
	// parse the target time
	targetTime, err := time.Parse(timeLayout, targetTimeStr)
	if err != nil {
		return false, fmt.Errorf("invalid target time format: %v", err)
	}

	if startTimeStr == "" {
		// parse the end time
		endTime, err := time.Parse(timeLayout, endTimeStr)
		if err != nil {
			return false, fmt.Errorf("invalid end time format: %v", err)
		}

		if targetTime.Before(endTime) {
			return true, nil
		}

		return false, nil
	} else if endTimeStr == "" {
		// parse the start time
		startTime, err := time.Parse(timeLayout, startTimeStr)
		if err != nil {
			return false, fmt.Errorf("invalid start time format: %v", err)
		}

		if targetTime.After(startTime) {
			return true, nil
		}

		return false, nil
	}

	// parse the start time
	startTime, err := time.Parse(timeLayout, startTimeStr)
	if err != nil {
		return false, fmt.Errorf("invalid start time format: %v", err)
	}

	// parse the end time
	endTime, err := time.Parse(timeLayout, endTimeStr)
	if err != nil {
		return false, fmt.Errorf("invalid end time format: %v", err)
	}

	// check if targetTime is between startTime and endTime
	if targetTime.After(startTime) && targetTime.Before(endTime) {
		return true, nil
	}

	return false, nil
}