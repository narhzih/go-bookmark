package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
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

	if err = u.Db.QueryRow(query, twitterId).Scan(
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
func (u userActions) GetUserById(userId int) (user models.User, err error) {
	query := `SELECT id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at FROM users where id=$1 LIMIT 1`
	if err = u.Db.QueryRow(query, userId).Scan(
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
	query := `SELECT id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at FROM users where email=$1 LIMIT 1`
	if err = u.Db.QueryRow(query, userEmail).Scan(
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
	query := `SELECT id, username, email, profile_name, cover_photo, twitter_id, created_at, modified_at  FROM users where username=$1 LIMIT 1`
	if err = u.Db.QueryRow(query, username).Scan(
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

func (u userActions) GetUserAndAuth(user models.User) (userAndAuth models.UserAuth, err error) {
	query := "SELECT hashed_password, origin FROM user_auth WHERE user_id=$1"
	err = u.Db.QueryRow(query, user.ID).Scan(
		&userAndAuth.HashedPassword,
		&userAndAuth.Origin,
	)
	if err != nil {
		return models.UserAuth{}, err
	}
	userAndAuth.User = user

	return userAndAuth, nil
}

func (u userActions) GetUserDeviceTokens(userID int64) ([]string, error) {
	var deviceTokens []string
	query := `SELECT device_tokens FROM users WHERE id=$1`
	err := u.Db.QueryRow(query, userID).Scan(pq.Array(&deviceTokens))
	if err != nil {
		return deviceTokens, err
	}
	return deviceTokens, nil
}

// ----------------------------------------------------------------
// --------------- UPDATE OPERATIONS ------------------------------
// ----------------------------------------------------------------

func (u userActions) UpdateUserPassword(userId int, password string) error {
	authQuery := "UPDATE user_auth SET hashed_password=$1 WHERE user_id=$2"
	_, err := u.Db.Exec(authQuery, password, userId)
	if err != nil {
		u.Logger.Err(err).Msg("Could not create user for auth")
		return err
	}

	return nil
}

func (u userActions) UpdateUser(updatedBody models.User) (models.User, error) {
	var user models.User
	selectQuery := "SELECT id, username, email, profile_name, cover_photo FROM users WHERE id=$1 LIMIT 1"
	err := u.Db.QueryRow(selectQuery, updatedBody.ID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return models.User{}, ErrNoRecord
		}
		return models.User{}, err
	}

	if len(updatedBody.Username) <= 0 && len(updatedBody.CovertPhoto) <= 0 && len(updatedBody.ProfileName) <= 0 {
		// Just return the original user if there's nothing to udpate

		return user, nil
	} else {
		if len(updatedBody.Username) > 0 && len(updatedBody.ProfileName) > 0 && len(updatedBody.CovertPhoto) > 0 {

			query := "UPDATE users SET username=$1, profile_name=$2, cover_photo=$3 WHERE id=$4 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.Username, updatedBody.ProfileName, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		} else if len(updatedBody.Username) > 0 && len(updatedBody.ProfileName) > 0 {

			query := "UPDATE users SET username=$1, profile_name=$2 WHERE id=$3 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.Username, updatedBody.ProfileName, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		} else if len(updatedBody.Username) > 0 && len(updatedBody.CovertPhoto) > 0 {
			// For usual edit

			query := "UPDATE users SET username=$1, cover_photo=$2 WHERE id=$3 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.Username, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		} else if len(updatedBody.ProfileName) > 0 && len(updatedBody.CovertPhoto) > 0 {
			// For usual edit

			query := "UPDATE users SET profile_name=$1, cover_photo=$2 WHERE id=$3 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.ProfileName, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		} else if len(updatedBody.Username) > 0 {

			query := "UPDATE users SET username=$1 WHERE id=$2 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.Username, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		} else if len(updatedBody.CovertPhoto) > 0 {
			u.Logger.Info().Msg("Saving cover photo to database")
			query := "UPDATE users SET cover_photo=$1 WHERE id=$2 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		} else if len(updatedBody.ProfileName) > 0 {

			query := "UPDATE users SET profile_name=$1 WHERE id=$2 RETURNING id, username, email, profile_name, cover_photo, twitter_id, modified_at"
			err = u.Db.QueryRow(query, updatedBody.ProfileName, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
				&user.TwitterId,
				&user.ModifiedAt,
			)
		}

		if err != nil {
			u.Logger.Err(err).Msg("Error from here from updating profile")
			u.Logger.Err(err).Msg(err.Error())
			return user, err
		}

		return user, nil
	}
}
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

func (u userActions) UpdateUserDeviceTokens(userID int64, deviceTokens []string) ([]string, error) {
	var userDeviceTokens []string
	query := `UPDATE users SET device_tokens=$1 WHERE id=$2 RETURNING device_tokens`
	err := u.Db.QueryRow(query, pq.Array(deviceTokens), userID).Scan(pq.Array(&userDeviceTokens))
	if err != nil {
		return []string{}, err
	}
	return userDeviceTokens, nil
}

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
