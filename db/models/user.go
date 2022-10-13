package models

import "time"

type GoogleClaim struct {
	Aud        string `json:"aud"`
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
}

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	ProfileName   string    `json:"profile_name"`
	TwitterId     string    `json:"twitter_id"`
	CovertPhoto   string    `json:"cover_photo"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type UserAuth struct {
	User           User
	HashedPassword string    `json:"hashed_password"`
	Origin         string    `json:"origin"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedAt     time.Time `json:"modified_at"`
}

type Profile struct {
	User      User
	Pipes     int
	Bookmarks int
}
