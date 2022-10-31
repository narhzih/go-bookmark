package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type userActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewUserActions(db *sql.DB, logger zerolog.Logger) repository.UserRepository {
	return userActions{
		Db:     db,
		Logger: logger,
	}
}

// ----------------------------------------------------------------
// --------------- CREATION OPERATIONS ----------------------------
// ----------------------------------------------------------------

// CreateUserByEmail - creates a basic user record in users table
// and also adds a record for that user in user_auth table
func (u userActions) CreateUserByEmail(user models.User, password string, authOrigin string) (models.User, error) {
	tx, err := u.Db.Begin()
	if err != nil {
		return models.User{}, err
	}
	var newUser models.User
	query := `
	INSERT INTO users 
	    (email, username, profile_name) 
	VALUES ($1, $2, $3) 
	RETURNING id, username, email, profile_name, cover_photo, twitter_id, email_verified, created_at, modified_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, query, user.Email, user.Username, user.ProfileName).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.ProfileName,
		&newUser.CovertPhoto,
		&newUser.TwitterId,
		&newUser.EmailVerified,
		&newUser.CreatedAt,
		&newUser.ModifiedAt,
	)

	if err != nil {
		tx.Rollback()
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "23505" {
				return newUser, ErrRecordExists
			}
		}
		return models.User{}, err
	}

	// create authentication record for the user
	var authQuery string
	var secondQueryValue string
	if authOrigin == "" || authOrigin == "DEFAULT" {
		secondQueryValue = password
		authQuery = "INSERT INTO user_auth (user_id, hashed_password) VALUES ($1, $2)"

	} else {
		secondQueryValue = authOrigin
		authQuery = "INSERT INTO user_auth (user_id, origin) VALUES ($1, $2)"

	}
	_, err = tx.ExecContext(ctx, authQuery, newUser.ID, secondQueryValue)
	if err != nil {
		tx.Rollback()
		u.Logger.Err(err).Msg("Could not create user for auth")
		return models.User{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.User{}, err
	}
	return newUser, nil
}

// ----------------------------------------------------------------
// --------------- RETRIEVAL OPERATIONS ---------------------------
// ----------------------------------------------------------------

// GetUserByTwitterID - Retrieves a user by their twitter id value
func (u userActions) GetUserByTwitterID(twitterId string) (user models.User, err error) {
	query := `
	SELECT id, username, email, profile_name, cover_photo, twitter_id, email_verified, created_at, modified_at 
	FROM users 
	WHERE twitter_id=$1 
	LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = u.Db.QueryRowContext(ctx, query, twitterId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
		&user.TwitterId,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.ModifiedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNoRecord
		}
		return models.User{}, err
	}
	return user, nil
}

// GetUserById - Retrieves a user by their registered ID
func (u userActions) GetUserById(userId int64) (user models.User, err error) {
	query := `
	SELECT id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at 
	FROM users 
	WHERE id=$1 
	LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = u.Db.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
		&user.TwitterId,
		&user.CreatedAt,
		&user.ModifiedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNoRecord
		}
		return models.User{}, err
	}
	return user, err
}

// GetUserByEmail - Retrieves a user by their email
func (u userActions) GetUserByEmail(userEmail string) (user models.User, err error) {
	query := `
	SELECT id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at 
	FROM users 
	WHERE email=$1 
	LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = u.Db.QueryRowContext(ctx, query, userEmail).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
		&user.TwitterId,
		&user.CreatedAt,
		&user.ModifiedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNoRecord
		}
		return models.User{}, err
	}
	return user, err
}

// GetUserByUsername - Retrieves a user by their username
func (u userActions) GetUserByUsername(username string) (user models.User, err error) {
	query := `
	SELECT id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at  
	FROM users 
	WHERE username=$1 
	LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = u.Db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
		&user.TwitterId,
		&user.CreatedAt,
		&user.ModifiedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNoRecord
		}
		return models.User{}, err
	}
	return user, err
}

// GetUserAndAuth gets user and their authentication credentials
func (u userActions) GetUserAndAuth(userId int64) (models.UserAuth, error) {
	var userAndAuth models.UserAuth

	query := `
	SELECT 
	    ua.hashed_password, 
	    ua.origin, 
	    u.id, 
	    u.username, 
	    u.email, 
	    u.email_verified, 
	    u.profile_name, 
	    u.cover_photo, 
	    u.twitter_id
	FROM user_auth ua
		INNER JOIN users u on u.id = ua.user_id
	WHERE ua.user_id=$1 
	LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, query, userId).Scan(
		&userAndAuth.HashedPassword,
		&userAndAuth.Origin,
		&userAndAuth.User.ID,
		&userAndAuth.User.Username,
		&userAndAuth.User.Email,
		&userAndAuth.User.EmailVerified,
		&userAndAuth.User.ProfileName,
		&userAndAuth.User.CovertPhoto,
		&userAndAuth.User.TwitterId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.UserAuth{}, ErrNoRecord
		}
		return models.UserAuth{}, err
	}

	return userAndAuth, nil
}

// GetUserDeviceTokens gets user device tokens
func (u userActions) GetUserDeviceTokens(userID int64) ([]string, error) {
	var deviceTokens []string
	query := `SELECT device_tokens FROM users WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, query, userID).Scan(pq.Array(&deviceTokens))
	if err != nil {
		return deviceTokens, err
	}
	return deviceTokens, nil
}

// ----------------------------------------------------------------
// --------------- UPDATE OPERATIONS ------------------------------
// ----------------------------------------------------------------

// UpdateUserPassword updates a user's password
func (u userActions) UpdateUserPassword(userId int64, password string) error {
	authQuery := "UPDATE user_auth SET hashed_password=$1 WHERE user_id=$2"
	_, err := u.Db.Exec(authQuery, password, userId)
	if err != nil {
		u.Logger.Err(err).Msg("Could not create user for auth")
		return err
	}

	return nil
}

// UpdateUser updates user information
func (u userActions) UpdateUser(updatedBody models.User) (models.User, error) {
	var user models.User
	query := `
	UPDATE users 
	SET
	    username=$2,
	    email=$3,
	    profile_name=$4,
	    cover_photo=$5,
	    twitter_id=$6
	    
	WHERE id=$1 
	RETURNING id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := u.Db.QueryRowContext(
		ctx,
		query,
		updatedBody.ID,
		updatedBody.Username,
		updatedBody.Email,
		updatedBody.ProfileName,
		updatedBody.CovertPhoto,
		updatedBody.TwitterId,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
		&user.TwitterId,
		&user.CreatedAt,
		&user.ModifiedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNoRecord
		}

		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "23505" {
				errKey := parseUserUpdateDuplicateMessage(dbErr.Error())
				switch errKey {
				case "username":
					return models.User{}, ErrDuplicateUsername
				case "email":
					return models.User{}, ErrDuplicateEmail
				case "twitter_id":
					return models.User{}, ErrDuplicateTwitterID
				}
			}
		}
		return models.User{}, err
	}

	return user, nil
}

// VerifyUser sets a user's email_verified field to true
func (u userActions) VerifyUser(user models.User) (models.User, error) {
	query := `UPDATE users SET email_verified=true WHERE id=$1 RETURNING id, email, username, profile_name, cover_photo`
	err := u.Db.QueryRow(query, user.ID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.ProfileName,
		&user.CovertPhoto,
	)

	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUserDeviceTokens updates user device tokens
func (u userActions) UpdateUserDeviceTokens(userID int64, deviceTokens []string) ([]string, error) {
	var userDeviceTokens []string
	query := `UPDATE users SET device_tokens=$1 WHERE id=$2 RETURNING device_tokens`
	err := u.Db.QueryRow(query, pq.Array(deviceTokens), userID).Scan(pq.Array(&userDeviceTokens))
	if err != nil {
		return []string{}, err
	}
	return userDeviceTokens, nil
}

// ConnectToTwitter sets user twitter_id field
func (u userActions) ConnectToTwitter(user models.User, twitterId string) (models.User, error) {
	updatedUser := models.User{}
	query := `UPDATE users SET twitter_id=$1 WHERE id=$2 RETURNING id, username, profile_name, email, twitter_id, cover_photo`
	err := u.Db.QueryRow(query, twitterId, user.ID).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.ProfileName,
		&updatedUser.Email,
		&updatedUser.TwitterId,
		&updatedUser.CovertPhoto,
	)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

// DisconnectTwitter  set's user twitter_id field to empty
func (u userActions) DisconnectTwitter(user models.User) (models.User, error) {
	updatedUser := models.User{}
	query := `UPDATE users SET twitter_id='', modified_at=now() WHERE id=$1 RETURNING id, username, profile_name, email, twitter_id, cover_photo`
	err := u.Db.QueryRow(query, user.ID).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.ProfileName,
		&updatedUser.Email,
		&updatedUser.TwitterId,
		&updatedUser.CovertPhoto,
	)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func parseUserUpdateDuplicateMessage(errMessage string) string {
	firstSplit := strings.Split(errMessage, ":")[1]
	secondSplit := strings.Split(firstSplit, "\"")[1]
	desiredKey := strings.Split(secondSplit, "_")[1]

	return desiredKey
}
