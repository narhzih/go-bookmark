package models

import "time"

type Notification struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	Message   string      `json:"message"`
	MetaData  interface{} `json:"meta_data"`
	Read      bool        `json:"read"`
	CreatedAt time.Time   `json:"created_at"`
}
