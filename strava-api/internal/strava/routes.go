package strava

import "time"

// Route is a route within Strava.
type Route struct {
	Private  bool    `json:"private"`
	Distance float64 `json:"distance"`
	Athlete  struct {
		Id            int64 `json:"id"`
		ResourceState int   `json:"resource_state"`
	} `json:"athlete"`
	Description         string    `json:"description"`
	CreatedAt           time.Time `json:"created_at"`
	ElevationGain       float64   `json:"elevation_gain"`
	Type                int       `json:"type"`
	EstimatedMovingTime int       `json:"estimated_moving_time"`
	Segments            []struct {
		Country         string  `json:"country"`
		Private         bool    `json:"private"`
		Distance        float64 `json:"distance"`
		AverageGrade    float64 `json:"average_grade"`
		MaximumGrade    float64 `json:"maximum_grade"`
		ClimbCategory   int     `json:"climb_category"`
		City            string  `json:"city"`
		ElevationHigh   float64 `json:"elevation_high"`
		AthletePrEffort struct {
			Distance       float64   `json:"distance"`
			StartDateLocal time.Time `json:"start_date_local"`
			ActivityId     int       `json:"activity_id"`
			ElapsedTime    int       `json:"elapsed_time"`
			IsKom          bool      `json:"is_kom"`
			Id             int       `json:"id"`
			StartDate      time.Time `json:"start_date"`
		} `json:"athlete_pr_effort"`
		AthleteSegmentStats struct {
			PrElapsedTime int       `json:"pr_elapsed_time"`
			PrDate        time.Time `json:"pr_date"`
			EffortCount   int       `json:"effort_count"`
			PrActivityId  int       `json:"pr_activity_id"`
		} `json:"athlete_segment_stats"`
		StartLatlng  string  `json:"start_latlng"`
		ElevationLow float64 `json:"elevation_low"`
		EndLatlng    string  `json:"end_latlng"`
		ActivityType string  `json:"activity_type"`
		Name         string  `json:"name"`
		Id           int     `json:"id"`
		State        string  `json:"state"`
	} `json:"segments"`
	Starred   bool      `json:"starred"`
	UpdatedAt time.Time `json:"updated_at"`
	SubType   int       `json:"sub_type"`
	IdStr     string    `json:"id_str"`
	Name      string    `json:"name"`
	Id        int       `json:"id"`
	Map       struct {
		SummaryPolyline string `json:"summary_polyline"`
		Id              string `json:"id"`
		Polyline        string `json:"polyline"`
	} `json:"map"`
	Timestamp int `json:"timestamp"`
}
