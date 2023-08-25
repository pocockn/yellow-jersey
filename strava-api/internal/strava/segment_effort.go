package strava

import "time"

// SegmentEffortDetailed is a wrapper around SegmentEffortSummary.
type SegmentEffortDetailed struct {
	SegmentEffortSummary
}

// SegmentEffortSummary holds an effort for a segment.
type SegmentEffortSummary struct {
	EffortSummary
	Segment          SegmentSummary `json:"segment"`
	AverageCadence   float64        `json:"average_cadence"`
	AveragePower     float64        `json:"average_watts"`
	AverageHeartrate float64        `json:"average_heartrate"`
	MaximumHeartrate float64        `json:"max_heartrate"`
	KOMRank          int            `json:"kom_rank"` // 1-10 rank on segment at time of upload
	PRRank           int            `json:"pr_rank"`  // 1-3 personal record on segment at time of upload
	Hidden           bool           `json:"hidden"`
}

// EffortSummary is the base object for BestEfforts, SegmentEfforts and LapEfforts
type EffortSummary struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Activity struct {
		Id int64 `json:"id"`
	} `json:"activity"`
	Athlete struct {
		Id int64 `json:"id"`
	} `json:"athlete"`
	Distance       float64   `json:"distance"`
	MovingTime     int       `json:"moving_time"`
	ElapsedTime    int       `json:"elapsed_time"`
	StartIndex     int       `json:"start_index"`
	EndIndex       int       `json:"end_index"`
	StartDate      time.Time `json:"start_date"`
	StartDateLocal time.Time `json:"start_date_local"`
}
