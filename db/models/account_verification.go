package models

type AccountVerification struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
	Used       bool   `json:"used"`
	ExpiresAt  string `json:"expires_at"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}
