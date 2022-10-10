package models

import "time"

type GoogleClaim struct {
	Aud             string `json:"aud"`
	Azp             string `json:"azp"`
	Email           string `json:"email"`
	EmailVerifiedAt string `json:"email_verified_at"`
	Exp             string `json:"exp"`
	FamilyName      string `json:"family_name"`
	GivenName       string `json:"given_name"`
	Iat             string `json:"iat"`
	Iss             string `json:"iss"`
	Locale          string `json:"locale"`
	Name            string `json:"name"`
	Picture         string `json:"picture"`
	Sub             string `json:"sub"`
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
