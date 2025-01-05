// Package plotting provides visualization functionality for health data using gonum/plot.
package plotting

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/gsweene2/health-data-parser/internal/models"
	"github.com/gsweene2/health-data-parser/internal/timeutil"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// PlotData contains the processed data ready for plotting.
type PlotData struct {
	SourceMap map[string]plotter.XYs
	Unit      string
	MinTime   time.Time
	MaxTime   time.Time
}

// Collection of colors used for different data sources in the plot.
var sourceColors = []color.RGBA{
	{R: 0, G: 255, B: 0, A: 255},   // Green
	{R: 128, G: 0, B: 128, A: 255}, // Purple
	{R: 255, G: 0, B: 0, A: 255},   // Red
	{R: 0, G: 0, B: 255, A: 255},   // Blue
	{R: 255, G: 165, B: 0, A: 255}, // Orange
	{R: 75, G: 0, B: 130, A: 255},  // Indigo
}

// getColorForSource returns a color for a given source index.
func getColorForSource(index int) color.RGBA {
	return sourceColors[index%len(sourceColors)]
}

// PreparePlotData processes health records into a format suitable for plotting.
// It extracts time series data for each source and determines the time range.
func PreparePlotData(records []models.Record) (PlotData, error) {
	sourceMap := make(map[string]plotter.XYs)
	var minTime, maxTime time.Time

	for i, record := range records {
		t, err := timeutil.ParseAppleHealthTime(record.StartDate)
		if err != nil {
			return PlotData{}, fmt.Errorf("error parsing time: %v", err)
		}

		if i == 0 || t.Before(minTime) {
			minTime = t
		}
		if i == 0 || t.After(maxTime) {
			maxTime = t
		}

		val, err := strconv.ParseFloat(record.Value, 64)
		if err != nil {
			return PlotData{}, fmt.Errorf("error parsing value: %v", err)
		}

		point := plotter.XY{
			X: float64(t.Unix()),
			Y: val,
		}
		sourceMap[record.SourceName] = append(sourceMap[record.SourceName], point)
	}

	return PlotData{
		SourceMap: sourceMap,
		Unit:      records[0].Unit,
		MinTime:   minTime,
		MaxTime:   maxTime,
	}, nil
}

// CustomTimeTicks implements plot.Ticker for custom time axis tick marks.
type CustomTimeTicks struct {
	TimeFormat string
	Interval   time.Duration
}

// Ticks implements the plot.Ticker interface, returning tick marks at specified intervals.
func (t CustomTimeTicks) Ticks(min, max float64) []plot.Tick {
	minTime := time.Unix(int64(min), 0)
	maxTime := time.Unix(int64(max), 0)

	minTime = minTime.Truncate(t.Interval)

	var ticks []plot.Tick
	for curr := minTime; curr.Before(maxTime.Add(t.Interval)); curr = curr.Add(t.Interval) {
		ticks = append(ticks, plot.Tick{
			Value: float64(curr.Unix()),
			Label: curr.Format(t.TimeFormat),
		})
	}
	return ticks
}

// ConfigurePlot sets up the plot with formatting, labels, and data series.
// It handles tick mark spacing based on the time range and adds padding for readability.
func ConfigurePlot(p *plot.Plot, data PlotData, title string) error {
	p.Title.Text = title
	p.X.Label.Text = "Time"
	p.Y.Label.Text = data.Unit
	p.Legend.Top = true
	p.Legend.Left = false

	duration := data.MaxTime.Sub(data.MinTime)
	var tickInterval time.Duration
	var timeFormat string

	if duration <= time.Hour {
		tickInterval = 10 * time.Minute
		timeFormat = "15:04"
	} else if duration <= 6*time.Hour {
		tickInterval = 30 * time.Minute
		timeFormat = "15:04"
	} else if duration <= 24*time.Hour {
		tickInterval = time.Hour
		timeFormat = "15:04"
	} else {
		tickInterval = 24 * time.Hour
		timeFormat = "2006-01-02\n15:04"
	}

	paddingDuration := time.Duration(float64(tickInterval) * 0.5)
	plotMinTime := data.MinTime.Add(-paddingDuration)
	plotMaxTime := data.MaxTime.Add(paddingDuration)

	p.X.Min = float64(plotMinTime.Unix())
	p.X.Max = float64(plotMaxTime.Unix())

	p.X.Tick.Marker = CustomTimeTicks{
		TimeFormat: timeFormat,
		Interval:   tickInterval,
	}

	colorIndex := 0
	for source, points := range data.SourceMap {
		line, err := plotter.NewLine(points)
		if err != nil {
			return fmt.Errorf("error creating line plot for %s: %v", source, err)
		}

		line.Color = getColorForSource(colorIndex)
		p.Legend.Add(source, line)
		p.Add(line)
		colorIndex++
	}

	return nil
}
