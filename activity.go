package warriorrec

import (
	"time"
)

// An ActivityCategory describes the category of an Activity.
type ActivityCategory struct {
	ID   string `json:"id"`
	URL  string `json:"url,omitempty"`
	Name string `json:"name"`
}

// An Activity describes a scheduled Warrior Recreation activity.
type Activity struct {
	ID         string `json:"id"`
	CategoryID string `json:"categoryId"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location"`
	Spots       int    `json:"spots"`

	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Cancelled bool      `json:"cancelled"`
}
