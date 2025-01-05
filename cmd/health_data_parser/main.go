// Package main provides a command-line tool for analyzing and visualizing Apple Health data.
package main

import (
	"os"

	"github.com/gsweene2/health-data-parser/internal/config"
	"github.com/gsweene2/health-data-parser/internal/logger"
	"github.com/gsweene2/health-data-parser/internal/processor"
	"github.com/gsweene2/health-data-parser/internal/writer"
)

func main() {
	log := logger.New()
	defer log.Sync()

	cfg, err := config.Parse()
	if err != nil {
		log.Errorw("Failed to parse configuration", "error", err)
		os.Exit(1)
	}

	proc := processor.NewDefaultProcessor(log)
	writer := writer.NewDefaultFileWriter(log)

	records, err := proc.Process(cfg)
	if err != nil {
		log.Errorw("Failed to process health data", "error", err)
		os.Exit(1)
	}

	if err := writer.WriteCSV(records, cfg.GetCSVPath()); err != nil {
		log.Errorw("Failed to write CSV", "error", err)
		os.Exit(1)
	}

	log.Infow("CSV file written", "path", cfg.GetCSVPath())

	if err := writer.WritePlot(records, cfg); err != nil {
		log.Errorw("Failed to create plot", "error", err)
		os.Exit(1)
	}
}
