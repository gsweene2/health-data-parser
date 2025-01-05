// Package timeutil provides utilities for handling time formats specific to Apple Health data.
package timeutil

import (
	"fmt"
	"time"
)

// AppleHealthTimeFormat defines the time format used in Apple Health export files.
// The format follows Go's reference time: "2006-01-02 15:04:05 -0700".
const AppleHealthTimeFormat = "2006-01-02 15:04:05 -0700"

// ParseAppleHealthTime converts a time string from Apple Health format to time.Time.
// Returns an error if the time string doesn't match the expected format.
func ParseAppleHealthTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(AppleHealthTimeFormat, timeStr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unsupported time format: %s", timeStr)
}
