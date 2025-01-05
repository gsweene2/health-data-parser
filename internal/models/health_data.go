// Package models defines the data structures for parsing Apple Health export data.
package models

// HealthData represents the root element of the Apple Health export XML.
type HealthData struct {
	Records []Record `xml:"Record"`
}

// Record represents a single health measurement from the Apple Health export.
type Record struct {
	Type       string   `xml:"type,attr"`
	SourceName string   `xml:"sourceName,attr"`
	Unit       string   `xml:"unit,attr"`
	Value      string   `xml:"value,attr"`
	StartDate  string   `xml:"startDate,attr"`
	EndDate    string   `xml:"endDate,attr"`
	Metadata   Metadata `xml:"MetadataEntry"`
}

// Metadata holds additional information attached to a health record.
type Metadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}
