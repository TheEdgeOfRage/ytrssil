package models

type Channel struct {
	// YouTube ID of the channel
	ID string `json:"channel_id" uri:"channel_id" binding:"required"`
	// Name of the channel
	Name string `json:"name"`
}
