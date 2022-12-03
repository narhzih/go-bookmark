package models

import "time"

const (
	PipeShareTypePublic  = "public"
	PipeShareTypePrivate = "private"
)

type SharedPipe struct {
	ID         int64     `json:"id"`
	SharerID   int64     `json:"sharer_id"`
	PipeID     int64     `json:"pipe_id"`
	Type       string    `json:"type"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type SharedPipeReceiver struct {
	ID           int64     `json:"id"`
	SharerId     int64     `json:"sharer_id"`
	SharedPipeId int64     `json:"shared_pipe_id"`
	ReceiverID   int64     `json:"receiver_id"`
	CreatedAt    time.Time `json:"created_at"`
	Code         string    `json:"code"`
	IsAccepted   bool      `json:"is_accepted"`
	ModifiedAt   time.Time `json:"modified_at"`
}

// MDPrivatePipeShare - Metadata definitions for notification
type MDPrivatePipeShare struct {
	Sharer User   `json:"sharer"`
	Pipe   Pipe   `json:"pipe"`
	Code   string `json:"code"`
}
