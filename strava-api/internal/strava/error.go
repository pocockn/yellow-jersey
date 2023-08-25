package strava

import "encoding/json"

// Error is used if we encounter an error when querying the Strava API.
type Error struct {
	Message string           `json:"message"`
	Errors  []*ErrorDetailed `json:"errors"`
}

// ErrorDetailed holds a more detailed error from Strava.
type ErrorDetailed struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

// Error satisfies the error interface.
func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
