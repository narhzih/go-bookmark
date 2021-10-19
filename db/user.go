package db

import (
	"database/sql"

	"github.com/lib/pq"
	"gitlab.com/gowagr/mypipe-api/db/model"
)

func (db Database) CreateUser(email, password string) (userID int64, err error) {
	var id int64
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id`
	err = db.Conn.QueryRow(query, email).Scan(&id)
	if err != nil {
		return id, err
	}

	authQuery := `INSERT INTO user_auth (user_id, password) VALUES ($1, $2)`
	_, err = db.Conn.Exec(authQuery, id, password)
	if err != nil {
		// Check if the error is due to duplicate recoreds
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "25505" {
				db.Logger.Err(dbErr).Msg("Duplicate record")
				return id, ErrRecordExists
			}
		}
		return id, err
	}

	return id, err
}

func (db Database) GetUserById(userId int) (user model.User, err error) {
	query := `SELECT * FROM users where id=$1 LIMIT 1`
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
