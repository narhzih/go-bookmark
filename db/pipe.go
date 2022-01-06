package db

import (
	"database/sql"
	"time"

	"gitlab.com/trencetech/mypipe-api/db/model"
)

func (db Database) PipeAlreadyExists(pipeName string, userId int64) (bool, error) {
	var pipe model.Pipe
	query := "SELECT id, name FROM pipes WHERE name=$1 AND user_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(query, pipeName, userId).Scan(&pipe.ID, &pipe.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			// NO record exits
			return false, nil
		}
		return true, err
	}
	return true, nil
}
func (db Database) CreatePipe(pipe model.Pipe) (model.Pipe, error) {
	var newPipe model.Pipe
	query := "INSERT INTO pipes (user_id, name, cover_photo) VALUES($1, $2, $3) RETURNING id, name, cover_photo"
	err := db.Conn.QueryRow(query, pipe.UserID, pipe.Name, pipe.CoverPhoto).Scan(
		&newPipe.ID,
		&newPipe.Name,
		&newPipe.CoverPhoto,
	)

	if err != nil {
		return model.Pipe{}, err
	}

	return newPipe, nil
}

func (db Database) GetPipe(pipeID, userID int64) (model.Pipe, error) {
	var pipe model.Pipe
	query := "SELECT id, name, cover_photo, created_at, user_id FROM pipes WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(query, pipeID, userID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Pipe{}, ErrNoRecord
		}
		return model.Pipe{}, err
	}
	return pipe, nil
}

func (db Database) GetPipeAndResource(pipeID, userID int64) (model.PipeAndResource, error) {
	var pipeAndR model.PipeAndResource
	query := "SELECT id, name, cover_photo, created_at, user_id FROM pipes WHERE id=$1 AND user_id=$2"
	err := db.Conn.QueryRow(query, pipeID, userID).Scan(
		&pipeAndR.Pipe.ID,
		&pipeAndR.Pipe.Name,
		&pipeAndR.Pipe.CoverPhoto,
		&pipeAndR.Pipe.CreatedAt,
		&pipeAndR.Pipe.UserID,
	)
	if err != nil {
		return model.PipeAndResource{}, nil
	}
	// get bookmarks
	pipeAndR.Bookmarks, err = db.GetBookmarks(userID, pipeID)
	if err != nil {
		return model.PipeAndResource{}, nil
	}
	return pipeAndR, nil
}

func (db Database) GetPipes(userID int64) ([]model.Pipe, error) {
	var pipes []model.Pipe
	query := "SELECT id, name, cover_photo, created_at, user_id FROM pipes WHERE user_id=$1"
	rows, err := db.Conn.Query(query, userID)
	if err != nil {
		return pipes, err
	}
	defer rows.Close()
	for rows.Next() {
		var pipe model.Pipe
		if err := rows.Scan(&pipe.ID, &pipe.Name, &pipe.CoverPhoto, &pipe.CreatedAt, &pipe.UserID); err != nil {
			return pipes, err
		}

		pipes = append(pipes, pipe)
	}

	if err := rows.Err(); err != nil {
		return pipes, err
	}
	return pipes, nil
}

func (db Database) GetPipesCount(userID int64) (int, error) {
	var pipesCount int
	query := "SELECT COUNT(id) FROM pipes WHERE user_id=$1"
	err := db.Conn.QueryRow(query, userID).Scan(&pipesCount)
	if err != nil {
		return pipesCount, err
	}

	return pipesCount, nil
}

func (db Database) UpdatePipe(userID int64, pipeID int64, updatedBody model.Pipe) (model.Pipe, error) {
	var pipe model.Pipe
	selectQuery := "SELECT id, name, cover_photo FROM pipes WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := db.Conn.QueryRow(selectQuery, pipeID, userID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Pipe{}, ErrNoRecord
		}
		return model.Pipe{}, err
	}

	if len(updatedBody.Name) <= 0 && len(updatedBody.CoverPhoto) <= 0 {
		return pipe, nil
	} else {
		if len(updatedBody.Name) > 0 && len(updatedBody.CoverPhoto) > 0 {
			query := "UPDATE pipes SET name=$1,cover_photo=$2, modified_at=$3 WHERE id=$4 AND user_id=$5 RETURNING id, name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.Name, updatedBody.CoverPhoto, time.Now(), pipeID, userID).Scan(
				&pipe.ID,
				&pipe.Name,
				&pipe.CoverPhoto,
			)
		} else if len(updatedBody.Name) > 0 {
			query := "UPDATE pipes SET name=$1, modified_at=$2 WHERE id=$3 AND user_id=$4 RETURNING id, name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.Name, time.Now(), pipeID, userID).Scan(
				&pipe.ID,
				&pipe.Name,
				&pipe.CoverPhoto,
			)
		} else if len(updatedBody.CoverPhoto) > 0 {
			query := "UPDATE pipes SET cover_photo=$1, modified_at=$2 WHERE id=$3 AND user_id=$4 RETURNING id, name, cover_photo"
			err = db.Conn.QueryRow(query, updatedBody.CoverPhoto, time.Now(), pipeID, userID).Scan(
				&pipe.ID,
				&pipe.Name,
				&pipe.CoverPhoto,
			)
		}

		if err != nil {

			return pipe, err
		}

		return pipe, nil

	}

}

func (db Database) UpdatePipeName(userID int64, pipeID int64, pipeName string) (model.Pipe, error) {
	var pipe model.Pipe
	selectQuery := "SELECT id, name FROM pipes WHERE user_id=$1 AND id=$2 LIMIT 1"
	err := db.Conn.QueryRow(selectQuery, userID, pipeID).Scan(
		&pipe.ID,
		&pipe.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Pipe{}, ErrNoRecord
		}
		return model.Pipe{}, err
	}

	// Update the pipe
	updateQuery := "UPDATE pipes SET name=$1, modified_at=CURRENT_TIMESTAMP WHERE id=$2 RETURNING id, name, cover_photo, created_at, modified_at, user_id"
	err = db.Conn.QueryRow(updateQuery, pipe.ID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.ModifiedAt,
		&pipe.UserID,
	)

	if err != nil {
		return model.Pipe{}, err
	}
	return pipe, nil
}
func (db Database) UpdatePipeCoverPhoto(userID int64, pipeID int64, pipeCp string) (model.Pipe, error) {
	var pipe model.Pipe
	selectQuery := "SELECT id, cover_photo FROM pipes WHERE user_id=$1 AND id=$2 LIMIT 1"
	err := db.Conn.QueryRow(selectQuery, userID, pipeID).Scan(
		&pipe.ID,
		&pipe.CoverPhoto,
	)
	if err != nil {
		return model.Pipe{}, nil
	}

	// Update the pipe
	updateQuery := "UPDATE pipes SET cover_photo=$1, modified_at=CURRENT_TIMESTAMP WHERE id=$2 RETURNING id, name, cover_photo, created_at, modified_at, user_id"
	err = db.Conn.QueryRow(updateQuery, pipe.CoverPhoto).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.ModifiedAt,
		&pipe.UserID,
	)

	return pipe, nil
}

// Whenever a pipe is deleted, all bookmarks under the pipe must be deleted
// along with the pipe. This has a already been taken care of by enabling
// ON DELETE CASCADE in the bookmarks table

func (db Database) DeletePipe(userID, pipeID int64) (bool, error) {
	deleteQuery := "DELETE FROM pipes WHERE id=$1 AND user_id=$2"
	_, err := db.Conn.Exec(deleteQuery, pipeID, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}
