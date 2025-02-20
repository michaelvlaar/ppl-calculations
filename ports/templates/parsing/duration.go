package parsing

import (
	"fmt"
	"strconv"
	"time"
)

func parseHHMMToDuration(input string) (time.Duration, error) {
	if len(input) != 4 {
		return 0, fmt.Errorf("invalid format, expected HHmm")
	}
	partHour := input[0:2]
	partMinute := input[2:4]

	hours, err := strconv.Atoi(partHour)
	if err != nil {
		return 0, fmt.Errorf("invalid hour value: %v", err)
	}

	minutes, err := strconv.Atoi(partMinute)
	if err != nil {
		return 0, fmt.Errorf("invalid minute value: %v", err)
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	return duration, nil
}
