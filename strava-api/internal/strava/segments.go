package strava

import "time"

// Segment is the response returned when querying for Segments from Strava.
type Segment struct {
	Id                 int       `json:"id"`
	ResourceState      int       `json:"resource_state"`
	Name               string    `json:"name"`
	ActivityType       string    `json:"activity_type"`
	Distance           float64   `json:"distance"`
	AverageGrade       float64   `json:"average_grade"`
	MaximumGrade       float64   `json:"maximum_grade"`
	ElevationHigh      float64   `json:"elevation_high"`
	ElevationLow       float64   `json:"elevation_low"`
	StartLatlng        []float64 `json:"start_latlng"`
	EndLatlng          []float64 `json:"end_latlng"`
	ClimbCategory      int       `json:"climb_category"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	Country            string    `json:"country"`
	Private            bool      `json:"private"`
	Hazardous          bool      `json:"hazardous"`
	Starred            bool      `json:"starred"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	Map                struct {
		Id            string `json:"id"`
		Polyline      string `json:"polyline"`
		ResourceState int    `json:"resource_state"`
	} `json:"map"`
	EffortCount         int                 `json:"effort_count"`
	AthleteCount        int                 `json:"athlete_count"`
	StarCount           int                 `json:"star_count"`
	AthleteSegmentStats AthleteSegmentStats `json:"athlete_segment_stats"`
}

// AthleteSegmentStats holds state for an athlete for a specific segment.
type AthleteSegmentStats struct {
	PrElapsedTime int    `json:"pr_elapsed_time"`
	PrDate        string `json:"pr_date"`
	EffortCount   int    `json:"effort_count"`
}

type SegmentSummary struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	ActivityType  string    `json:"activity_type"`
	Distance      float64   `json:"distance"`
	AverageGrade  float64   `json:"average_grade"`
	MaximumGrade  float64   `json:"maximum_grade"`
	ElevationHigh float64   `json:"elevation_high"`
	ElevationLow  float64   `json:"elevation_low"`
	ClimbCategory int       `json:"climb_category"`
	StartLocation []float64 `json:"start_latlng"`
	EndLocation   []float64 `json:"end_latlng"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	Private       bool      `json:"private"`
	Starred       bool      `json:"starred"`
}
