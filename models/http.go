package models

type VideosResponse struct {
	Videos []Video `json:"videos"`
}

type VideoURIRequest struct {
	VideoID string `uri:"video_id" binding:"required"`
}

type SetVideoProgressRequest struct {
	Progress string `form:"progress" binding:"required"`
}
