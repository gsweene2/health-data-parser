// Package config provides configuration management for health data parsing.
package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// Config holds the application's configuration parameters.
type Config struct {
	InputFile    string
	StartTime    time.Time
	EndTime      time.Time
	RecordType   string
	Sources      []string
	OutputPrefix string
}

// Validate checks if the configuration is valid and returns an error if not.
func (c *Config) Validate() error {
	if c.InputFile == "" {
		return fmt.Errorf("input file is required")
	}
	if c.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if c.EndTime.IsZero() {
		return fmt.Errorf("end time is required")
	}
	if c.StartTime.After(c.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}
	return nil
}

// GetCSVPath returns the path where the CSV output file should be written.
func (c *Config) GetCSVPath() string {
	return fmt.Sprintf("%s_%s_%s.csv",
		c.OutputPrefix,
		c.RecordType,
		c.StartTime.Format("2006-01-02"))
}

// GetPlotPath returns the path where the plot PNG file should be written.
func (c *Config) GetPlotPath() string {
	return fmt.Sprintf("%s_%s_%s.png",
		c.OutputPrefix,
		c.RecordType,
		c.StartTime.Format("2006-01-02"))
}

// Parse processes command-line flags and returns a configured Config struct.
func Parse() (*Config, error) {
	var rawConfig struct {
		InputFile    string
		StartTime    string
		EndTime      string
		RecordType   string
		Sources      string
		OutputPrefix string
	}

	flag.StringVar(&rawConfig.InputFile, "input", "", "Path to Apple Health data export XML file")
	flag.StringVar(&rawConfig.StartTime, "start", "", "Start time in RFC3339 format")
	flag.StringVar(&rawConfig.EndTime, "end", "", "End time in RFC3339 format")
	flag.StringVar(&rawConfig.RecordType, "type", "HKQuantityTypeIdentifierHeartRate", "Type of health record to filter")
	flag.StringVar(&rawConfig.Sources, "sources", "", "Comma-separated list of source names to filter (e.g., 'Apple Watch,WHOOP')")
	flag.StringVar(&rawConfig.OutputPrefix, "output-prefix", "filtered_health_data", "Prefix for output files")

	flag.Parse()

	config := &Config{
		InputFile:    rawConfig.InputFile,
		RecordType:   rawConfig.RecordType,
		OutputPrefix: rawConfig.OutputPrefix,
	}

	if rawConfig.Sources != "" {
		config.Sources = strings.Split(rawConfig.Sources, ",")
		for i := range config.Sources {
			config.Sources[i] = strings.TrimSpace(config.Sources[i])
		}
	}

	var err error
	config.StartTime, err = time.Parse(time.RFC3339, rawConfig.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %v", err)
	}

	config.EndTime, err = time.Parse(time.RFC3339, rawConfig.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %v", err)
	}

	return config, config.Validate()
}
