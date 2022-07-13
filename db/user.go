package db

import (
	"database/sql"

	"github.com/lib/pq"
	"gitlab.com/trencetech/mypipe-api/db/model"
)

func (db Database) CreateUserByEmail(user model.User, password string, authOrigin string) (newUser model.User, err error) {
	query := `INSERT INTO users (email, username, profile_name) VALUES ($1, $2, $3) RETURNING id, email, user, profile_name`
	err = db.Conn.QueryRow(query, user.Email, user.Username, user.ProfileName).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Username,
		&newUser.ProfileName,
	)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "23505" {
				db.Logger.Err(dbErr).Msg("duplicate record")
				return newUser, ErrRecordExists
			}
		}
		return model.User{}, err
	}
	var authQuery string
	var secondQueryValue string
	if authOrigin == "" || authOrigin == "DEFAULT" {
		secondQueryValue = password
		authQuery = "INSERT INTO user_auth (user_id, hashed_password) VALUES ($1, $2)"

	} else {
		secondQueryValue = authOrigin
		authQuery = "INSERT INTO user_auth (user_id, origin) VALUES ($1, $2)"

	}
	_, err = db.Conn.Exec(authQuery, newUser.ID, secondQueryValue)
	if err != nil {
		db.Logger.Err(err).Msg("Could not create user for auth")
		return model.User{}, err
	}

	return newUser, err
}

func (db Database) CreateUser(user model.User) (newUser model.User, err error) {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email`
	err = db.Conn.QueryRow(query, user.Username, user.Email).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
	)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "23505" {
				db.Logger.Err(dbErr).Msg("duplicate record")
				return newUser, ErrRecordExists
			}
		}
		return model.User{}, err
	}

	return newUser, err
}

func (db Database) GetUserById(userId int) (user model.User, err error) {
	query := `SELECT id, username, email, profile_name, cover_photo FROM users where id=$1 LIMIT 1`
	if err = db.Conn.QueryRow(query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrNoRecord
		}
		return model.User{}, err
	}
	return user, err
}

func (db Database) GetUserByEmail(userEmail string) (user model.User, err error) {
	query := `SELECT id, username, email FROM users where email=$1 LIMIT 1`
	if err = db.Conn.QueryRow(query, userEmail).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrNoRecord
		}
		return model.User{}, err
	}
	return user, err
}

func (db Database) GetUserByUsername(username string) (user model.User, err error) {
	query := `SELECT id, username, email, profile_name FROM users where username=$1 LIMIT 1`
	if err = db.Conn.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrNoRecord
		}
		return model.User{}, err
	}
	return user, err
}

func (db Database) GetUserAndAuth(user model.User) (userAndAuth model.UserAuth, err error) {
	query := "SELECT hashed_password, origin FROM user_auth WHERE user_id=$1"
	err = db.Conn.QueryRow(query, user.ID).Scan(
		&userAndAuth.HashedPassword,
		&userAndAuth.Origin,
	)
	if err != nil {
		return model.UserAuth{}, err
	}
	userAndAuth.User = user

	return userAndAuth, nil
}

func (db Database) UpdateUserPassword(userId int, password string) error {
	authQuery := "UPDATE user_auth SET hashed_password=$1 WHERE user_id=$2"
	_, err := db.Conn.Exec(authQuery, password, userId)
	if err != nil {
		db.Logger.Err(err).Msg("Could not create user for auth")
		return err
	}

	return nil
}

// Find a way to improve sql update query

func (db Database) UpdateUser(updatedBody model.User) (model.User, error) {
	var user model.User
	selectQuery := "SELECT id, username, email, profile_name, cover_photo FROM users WHERE id=$1 LIMIT 1"
	err := db.Conn.QueryRow(selectQuery, updatedBody.ID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileName,
		&user.CovertPhoto,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return model.User{}, ErrNoRecord
		}
		return model.User{}, err
	}

	if len(updatedBody.Username) <= 0 && len(updatedBody.CovertPhoto) <= 0 && len(updatedBody.ProfileName) <= 0 {
		// Just return the original user if there's nothing to udpate

		return user, nil
	} else {
		if len(updatedBody.Username) > 0 && len(updatedBody.ProfileName) > 0 && len(updatedBody.CovertPhoto) > 0 {

			query := "UPDATE users SET username=$1, profile_name=$2, cover_photo=$3 WHERE id=$4 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.Username, updatedBody.ProfileName, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		} else if len(updatedBody.Username) > 0 && len(updatedBody.ProfileName) > 0 {

			query := "UPDATE users SET username=$1, profile_name=$2 WHERE id=$3 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.Username, updatedBody.ProfileName, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		} else if len(updatedBody.Username) > 0 && len(updatedBody.CovertPhoto) > 0 {
			// For usual edit

			query := "UPDATE users SET username=$1, cover_photo=$2 WHERE id=$3 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.Username, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		} else if len(updatedBody.ProfileName) > 0 && len(updatedBody.CovertPhoto) > 0 {
			// For usual edit

			query := "UPDATE users SET profile_name=$1, cover_photo=$2 WHERE id=$3 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.ProfileName, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		} else if len(updatedBody.Username) > 0 {

			query := "UPDATE users SET username=$1 WHERE id=$2 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.Username, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		} else if len(updatedBody.CovertPhoto) > 0 {
			db.Logger.Info().Msg("Saving cover photo to database")
			query := "UPDATE users SET cover_photo=$1 WHERE id=$2 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.CovertPhoto, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		} else if len(updatedBody.ProfileName) > 0 {

			query := "UPDATE users SET profile_name=$1 WHERE id=$2 RETURNING id, username, email, profile_name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.ProfileName, updatedBody.ID).Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileName,
				&user.CovertPhoto,
			)
		}

		if err != nil {
			db.Logger.Err(err).Msg("Error from here from updating profile")
			db.Logger.Err(err).Msg(err.Error())
			return user, err
		}

		return user, nil
	}
}

func (db Database) VerifyUser(user model.User) (model.User, error) {
	query := `UPDATE users SET email_verified=true WHERE id=$1 RETURNING id, email, username, profile_name, cover_photo`
	err := db.Conn.QueryRow(query, user.ID).Scan(
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

func (db Database) GetUserDeviceTokens(userID int64) ([]string, error) {
	var deviceTokens []string
	query := `SELECT device_tokens FROM users WHERE user_id=$1`
	err := db.Conn.QueryRow(query, userID).Scan(pq.Array(&deviceTokens))
	if err != nil {
		return deviceTokens, err
	}
	return deviceTokens, nil
}

func (db Database) UpdateUserDeviceTokens(userID int64, deviceTokens []string) ([]string, error) {
	var userDeviceTokens []string
	query := `UPDATE users SET device_tokens=$1 WHERE user_id=$2 RETURNING device_tokens`
	err := db.Conn.QueryRow(query, pq.Array(deviceTokens), userID).Scan(pq.Array(&userDeviceTokens))
	if err != nil {
		return []string{}, err
	}
	return userDeviceTokens, nil
}
