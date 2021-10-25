package db

import (
	"database/sql"

	"github.com/lib/pq"
	"gitlab.com/gowagr/mypipe-api/db/model"
)

func (db Database) CreateUserByEmail(user model.User, password string) (newUser model.User, err error) {
	query := `INSERT INTO users (username, email) VALUES ($1) RETURNING id, username, email`
	err = db.Conn.QueryRow(query, user.Email).Scan(
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

	authQuery := "INSERT INTO users_auth (user_id, password) VALUES ($1, $2)"
	_, err = db.Conn.Exec(authQuery, user.ID, password)
	if err != nil {
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
	query := `SELECT id, username, email FROM users where id=$1 LIMIT 1`
	if err = db.Conn.QueryRow(query, userId).Scan(
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
