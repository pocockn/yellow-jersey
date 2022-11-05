package strava

import "time"

// AthleteDetailed is the response returned from Strava when querying for an Athlete.
type AthleteDetailed struct {
	AthleteSummary
	Email                 string         `json:"email"`
	FollowerCount         int            `json:"follower_count"`
	FriendCount           int            `json:"friend_count"`
	MutualFriendCount     int            `json:"mutual_friend_count"`
	DatePreference        string         `json:"date_preference"`
	MeasurementPreference string         `json:"measurement_preference"`
	FTP                   int            `json:"ftp"`
	Weight                float64        `json:"weight"` // kilograms
	Clubs                 []*ClubSummary `json:"clubs"`
	Bikes                 []*GearSummary `json:"bikes"`
	Shoes                 []*GearSummary `json:"shoes"`
}

// AthleteSummary is a section of the response returned from Stravas API when querying for an Athlete.
type AthleteSummary struct {
	AthleteMeta
	FirstName        string    `json:"firstname"`
	LastName         string    `json:"lastname"`
	ProfileMedium    string    `json:"profile_medium"` // URL to a 62x62 pixel profile picture
	Profile          string    `json:"profile"`        // URL to a 124x124 pixel profile picture
	City             string    `json:"city"`
	State            string    `json:"state"`
	Country          string    `json:"country"`
	Gender           Gender    `json:"sex"`
	Friend           string    `json:"friend"`   // ‘pending’, ‘accepted’, ‘blocked’ or ‘null’, the authenticated athlete’s following status of this athlete
	Follower         string    `json:"follower"` // this athlete’s following status of the authenticated athlete
	Premium          bool      `json:"premium"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ApproveFollowers bool      `json:"approve_followers"` // if has enhanced privacy enabled
	BadgeTypeId      int       `json:"badge_type_id"`
}

// AthleteMeta holds the id of the athlete
type AthleteMeta struct {
	Id int64 `json:"id"`
}

// ClubSummary is a list of clubs the athlete is associated with.
type ClubSummary struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	ProfileMedium string `json:"profile_medium"` // URL to a 62x62 pixel profile picture
	Profile       string `json:"profile"`        // URL to a 124x124 pixel profile picture
}

// GearSummary is gear the athlete has on their account.
type GearSummary struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Primary  bool    `json:"primary"`
	Distance float64 `json:"distance"`
}

// Gender holds the athletes gender.
type Gender string

// Genders is a list of possible genders from Strava.
var Genders = struct {
	Unspecified Gender
	Male        Gender
	Female      Gender
}{"", "M", "F"}
