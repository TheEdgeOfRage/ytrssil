package models

type Channel struct {
	// YouTube ID of the channel
	ID string `json:"channel_id"`
	// Name of the channel
	Name string `json:"name"`
	// Subscribed indicates if the user is subscribed to this channel
	Subscribed bool `json:"subscribed"`
	// UnwatchedCount is the number of unwatched videos from this channel
	UnwatchedCount int `json:"unwatched_count"`
	// ImageURL is the URL of the channel's profile image
	ImageURL string `json:"image_url"`
}
