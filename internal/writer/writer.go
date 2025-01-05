// Package writer provides functionality for outputting health data to various formats.
package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/gsweene2/health-data-parser/internal/config"
	"github.com/gsweene2/health-data-parser/internal/models"
	"github.com/gsweene2/health-data-parser/internal/plotting"
	"github.com/gsweene2/health-data-parser/internal/timeutil"

	"go.uber.org/zap"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

// FileWriter defines the interface for writing health records to different formats.
type FileWriter interface {
	WriteCSV(records []models.Record, path string) error
	WritePlot(records []models.Record, config *config.Config) error
}

// DefaultFileWriter implements FileWriter with standard file operations.
type DefaultFileWriter struct {
	logger *zap.SugaredLogger
}

// NewDefaultFileWriter creates a new DefaultFileWriter with the provided logger.
func NewDefaultFileWriter(logger *zap.SugaredLogger) *DefaultFileWriter {
	return &DefaultFileWriter{logger: logger}
}

// WriteCSV implements FileWriter.WriteCSV. It writes health records to a CSV file
// at the specified path with columns for date, source, type, unit, and value.
func (w *DefaultFileWriter) WriteCSV(records []models.Record, path string) error {
	csvFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating file: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	header := []string{"Date", "Source", "Type", "Unit", "Value"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	for _, record := range records {
		startDateTime, err := timeutil.ParseAppleHealthTime(record.StartDate)
		if err != nil {
			return fmt.Errorf("error parsing record start time: %v", err)
		}

		csvRecord := []string{
			startDateTime.Format(time.RFC3339),
			record.SourceName,
			record.Type,
			record.Unit,
			record.Value,
		}
		if err := writer.Write(csvRecord); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}
	}

	return nil
}

// WritePlot implements FileWriter.WritePlot. It generates a visualization of the
// health records and saves it as a PNG file using the config's plot path.
func (w *DefaultFileWriter) WritePlot(records []models.Record, cfg *config.Config) error {
	data, err := plotting.PreparePlotData(records)
	if err != nil {
		return fmt.Errorf("error preparing plot data: %v", err)
	}

	p := plot.New()
	title := fmt.Sprintf("%s (%s - %s)", cfg.RecordType,
		cfg.StartTime.Format(time.RFC3339),
		cfg.EndTime.Format(time.RFC3339))

	if err := plotting.ConfigurePlot(p, data, title); err != nil {
		return fmt.Errorf("error configuring plot: %v", err)
	}

	if err := p.Save(8*vg.Inch, 6*vg.Inch, cfg.GetPlotPath()); err != nil {
		return fmt.Errorf("error saving plot: %v", err)
	}

	w.logger.Infow("Plot saved", "path", cfg.GetPlotPath())
	return nil
}
