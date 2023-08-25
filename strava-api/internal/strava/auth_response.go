package strava

// AuthorizationResponse is returned as a result of the token exchange
type AuthorizationResponse struct {
	AccessToken  string          `json:"access_token"`
	Athlete      AthleteDetailed `json:"athlete"`
	RefreshToken string          `json:"refresh_token"`
	State        string          `json:"State"`
}
