package parsing

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseHHMMToDuration(input string) (time.Duration, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format, expected HH:mm")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour value: %v", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute value: %v", err)
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	return duration, nil
}
