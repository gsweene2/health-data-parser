// Package processor provides functionality for processing Apple Health data XML files.
package processor

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/gsweene2/health-data-parser/internal/config"
	"github.com/gsweene2/health-data-parser/internal/models"
	"github.com/gsweene2/health-data-parser/internal/timeutil"

	"go.uber.org/zap"
)

// HealthDataProcessor defines the interface for processing health data records.
type HealthDataProcessor interface {
	Process(config *config.Config) ([]models.Record, error)
}

// DefaultProcessor implements HealthDataProcessor using standard XML processing.
type DefaultProcessor struct {
	logger *zap.SugaredLogger
}

// NewDefaultProcessor creates a new DefaultProcessor with the provided logger.
func NewDefaultProcessor(logger *zap.SugaredLogger) *DefaultProcessor {
	return &DefaultProcessor{logger: logger}
}

// shouldIncludeSource checks if a source should be included based on the config.
func (p *DefaultProcessor) shouldIncludeSource(sourceName string, sources []string) bool {
	if len(sources) == 0 {
		return true
	}

	for _, s := range sources {
		if s == sourceName {
			return true
		}
	}
	return false
}

// Process implements HealthDataProcessor.Process. It reads and filters health records
// based on type, time range, and sources specified in the config.
func (p *DefaultProcessor) Process(config *config.Config) ([]models.Record, error) {
	fileContent, err := os.ReadFile(config.InputFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var healthData models.HealthData
	if err := xml.Unmarshal(fileContent, &healthData); err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	var filteredRecords []models.Record
	for _, record := range healthData.Records {
		if record.Type != config.RecordType {
			continue
		}

		if !p.shouldIncludeSource(record.SourceName, config.Sources) {
			continue
		}

		recordTime, err := timeutil.ParseAppleHealthTime(record.StartDate)
		if err != nil {
			p.logger.Warnw("Failed to parse record time", "error", err)
			continue
		}

		if recordTime.After(config.StartTime) && recordTime.Before(config.EndTime) {
			filteredRecords = append(filteredRecords, record)
		}
	}

	if len(filteredRecords) == 0 {
		return nil, fmt.Errorf("no matching records found")
	}

	p.logger.Infow("Filtered records",
		"total", len(filteredRecords),
		"type", config.RecordType,
		"sources", config.Sources)

	return filteredRecords, nil
}
