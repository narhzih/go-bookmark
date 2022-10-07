package models

type PasswordReset struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	Validated  bool   `json:"validated"`
}
