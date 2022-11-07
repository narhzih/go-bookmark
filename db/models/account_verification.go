package models

import "time"

type AccountVerification struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Token      string    `json:"token"`
	Used       bool      `json:"used"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
