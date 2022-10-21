package postgres

import (
	"database/sql"
	"github.com/mypipeapp/mypipeapi/db/models"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"github.com/rs/zerolog"
	"time"
)

type pipeActions struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewPipeActions(db *sql.DB, logger zerolog.Logger) repository.PipeRepository {
	return pipeActions{
		Db:     db,
		Logger: logger,
	}
}

func (p pipeActions) PipeAlreadyExists(pipeName string, userId int64) (bool, error) {
	var pipe models.Pipe
	query := "SELECT id, name FROM pipes WHERE name=$1 AND user_id=$2 LIMIT 1"
	err := p.Db.QueryRow(query, pipeName, userId).Scan(&pipe.ID, &pipe.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			// NO record exits
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (p pipeActions) CreatePipe(pipe models.Pipe) (models.Pipe, error) {
	var newPipe models.Pipe
	query := "INSERT INTO pipes (user_id, name, cover_photo) VALUES($1, $2, $3) RETURNING id, name, cover_photo, user_id"
	err := p.Db.QueryRow(query, pipe.UserID, pipe.Name, pipe.CoverPhoto).Scan(
		&newPipe.ID,
		&newPipe.Name,
		&newPipe.CoverPhoto,
		&newPipe.UserID,
	)

	if err != nil {
		return models.Pipe{}, err
	}

	return newPipe, nil
}

func (p pipeActions) GetPipe(pipeID, userID int64) (models.Pipe, error) {
	var pipe models.Pipe
	query := `
	SELECT p.id, p.name, p.cover_photo, p.created_at, p.modified_at, p.user_id, COUNT(b.pipe_id) AS total_bookmarks, u.username
	FROM pipes p
		LEFT JOIN bookmarks b ON p.id=b.pipe_id
		LEFT JOIN users u ON p.user_id=u.id
	WHERE p.user_id=$1 AND p.id = $2
	GROUP BY p.id, u.username
	ORDER BY p.id
	LIMIT 1
	`
	err := p.Db.QueryRow(query, userID, pipeID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.ModifiedAt,
		&pipe.UserID,
		&pipe.Bookmarks,
		&pipe.Creator,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Pipe{}, ErrNoRecord
		}
		return models.Pipe{}, err
	}
	return pipe, nil

}

func (p pipeActions) GetPipeByName(pipeName string, userID int64) (models.Pipe, error) {
	var pipe models.Pipe
	query := "SELECT id, name, cover_photo, created_at, modified_at, user_id FROM pipes WHERE name=$1 AND user_id=$2 LIMIT 1"
	err := p.Db.QueryRow(query, pipeName, userID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.ModifiedAt,
		&pipe.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			p.Logger.Info().Msg("No pipe was found when searching...")
			return models.Pipe{}, ErrNoRecord
		}
		return models.Pipe{}, err
	}
	return pipe, nil
}

func (p pipeActions) GetPipeAndResource(pipeID, userID int64) (models.PipeAndResource, error) {
	var pipeAndR models.PipeAndResource
	query := "SELECT id, name, cover_photo, created_at, modified_at, user_id FROM pipes WHERE id=$1 AND user_id=$2"
	err := p.Db.QueryRow(query, pipeID, userID).Scan(
		&pipeAndR.Pipe.ID,
		&pipeAndR.Pipe.Name,
		&pipeAndR.Pipe.CoverPhoto,
		&pipeAndR.Pipe.CreatedAt,
		&pipeAndR.Pipe.ModifiedAt,
		&pipeAndR.Pipe.UserID,
	)
	if err != nil {
		return models.PipeAndResource{}, nil
	}
	// get bookmarks
	// TODO: uncomment the below line
	bActions := NewBookmarkActions(p.Db, p.Logger)
	pipeAndR.Bookmarks, err = bActions.GetBookmarks(userID, pipeID)
	if err != nil {
		return models.PipeAndResource{}, nil
	}
	return pipeAndR, nil
}

func (p pipeActions) GetPipesOnSteroid(userID int64) ([]models.Pipe, error) {
	var pipes []models.Pipe
	query := `
				SELECT p.id, p.name, p.cover_photo, p.created_at, p.modified_at, p.user_id, COUNT(b.pipe_id) AS total_bookmarks 
				FROM pipes p 
				    LEFT JOIN bookmarks b ON p.id=b.pipe_id 
				WHERE p.user_id=$1 
				GROUP BY p.id
	`
	rows, err := p.Db.Query(query, userID)
	if err != nil {
		return pipes, err
	}
	defer rows.Close()
	for rows.Next() {
		var pipe models.Pipe
		if err := rows.Scan(
			&pipe.ID,
			&pipe.Name,
			&pipe.CoverPhoto,
			&pipe.CreatedAt,
			&pipe.ModifiedAt,
			&pipe.UserID,
			&pipe.Bookmarks,
		); err != nil {
			return pipes, err
		}

		pipes = append(pipes, pipe)
	}
	if err := rows.Err(); err != nil {
		return pipes, err
	}
	return pipes, nil
}

func (p pipeActions) GetPipes(userID int64) ([]models.Pipe, error) {
	var pipes []models.Pipe
	query := `
			SELECT p.id, p.name, p.cover_photo, p.created_at, p.modified_at, p.user_id, COUNT(b.pipe_id) AS total_bookmarks, u.username
			FROM pipes p
				LEFT JOIN bookmarks b ON p.id=b.pipe_id
				LEFT JOIN users u ON p.user_id=u.id
			WHERE p.user_id=$1 OR p.id  IN (
					SELECT spr.shared_pipe_id FROM shared_pipe_receivers spr WHERE receiver_id=$1 AND is_accepted=true
				)
			GROUP BY p.id, u.username
			ORDER BY p.id;
	`
	rows, err := p.Db.Query(query, userID)
	if err != nil {
		return pipes, err
	}
	defer rows.Close()
	for rows.Next() {
		var pipe models.Pipe
		if err := rows.Scan(
			&pipe.ID,
			&pipe.Name,
			&pipe.CoverPhoto,
			&pipe.CreatedAt,
			&pipe.ModifiedAt,
			&pipe.UserID,
			&pipe.Bookmarks,
			&pipe.Creator,
		); err != nil {
			return pipes, err
		}

		pipes = append(pipes, pipe)
	}

	if err := rows.Err(); err != nil {
		return pipes, err
	}
	return pipes, nil
}

func (p pipeActions) GetPipesCount(userID int64) (int, error) {
	var pipesCount int
	query := "SELECT COUNT(id) FROM pipes WHERE user_id=$1"
	err := p.Db.QueryRow(query, userID).Scan(&pipesCount)
	if err != nil {
		return pipesCount, err
	}

	return pipesCount, nil
}

func (p pipeActions) UpdatePipe(userID int64, pipeID int64, updatedBody models.Pipe) (models.Pipe, error) {
	var pipe models.Pipe
	selectQuery := "SELECT id, name, cover_photo FROM pipes WHERE id=$1 AND user_id=$2 LIMIT 1"
	err := p.Db.QueryRow(selectQuery, pipeID, userID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Pipe{}, ErrNoRecord
		}
		return models.Pipe{}, err
	}

	if len(updatedBody.Name) <= 0 && len(updatedBody.CoverPhoto) <= 0 {
		return pipe, nil
	} else {
		if len(updatedBody.Name) > 0 && len(updatedBody.CoverPhoto) > 0 {
			query := "UPDATE pipes SET name=$1,cover_photo=$2, modified_at=$3 WHERE id=$4 AND user_id=$5 RETURNING id, name, cover_photo, modified_at"
			err = p.Db.QueryRow(query, updatedBody.Name, updatedBody.CoverPhoto, time.Now(), pipeID, userID).Scan(
				&pipe.ID,
				&pipe.Name,
				&pipe.CoverPhoto,
				&pipe.ModifiedAt,
			)
		} else if len(updatedBody.Name) > 0 {
			query := "UPDATE pipes SET name=$1, modified_at=$2 WHERE id=$3 AND user_id=$4 RETURNING id, name, cover_photo, modified_at"
			err = p.Db.QueryRow(query, updatedBody.Name, time.Now(), pipeID, userID).Scan(
				&pipe.ID,
				&pipe.Name,
				&pipe.CoverPhoto,
				&pipe.ModifiedAt,
			)
		} else if len(updatedBody.CoverPhoto) > 0 {
			query := "UPDATE pipes SET cover_photo=$1, modified_at=$2 WHERE id=$3 AND user_id=$4 RETURNING id, name, cover_photo, modified_at"
			err = p.Db.QueryRow(query, updatedBody.CoverPhoto, time.Now(), pipeID, userID).Scan(
				&pipe.ID,
				&pipe.Name,
				&pipe.CoverPhoto,
				&pipe.ModifiedAt,
			)
		}

		if err != nil {

			return pipe, err
		}

		return pipe, nil

	}

}

func (p pipeActions) UpdatePipeName(userID int64, pipeID int64, pipeName string) (models.Pipe, error) {
	var pipe models.Pipe
	selectQuery := "SELECT id, name FROM pipes WHERE user_id=$1 AND id=$2 LIMIT 1"
	err := p.Db.QueryRow(selectQuery, userID, pipeID).Scan(
		&pipe.ID,
		&pipe.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Pipe{}, ErrNoRecord
		}
		return models.Pipe{}, err
	}

	// Update the pipe
	updateQuery := "UPDATE pipes SET name=$1, modified_at=now() WHERE id=$2 RETURNING id, name, cover_photo, created_at, modified_at, user_id"
	err = p.Db.QueryRow(updateQuery, pipe.ID).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.ModifiedAt,
		&pipe.UserID,
	)

	if err != nil {
		return models.Pipe{}, err
	}
	return pipe, nil
}

func (p pipeActions) UpdatePipeCoverPhoto(userID int64, pipeID int64, pipeCp string) (models.Pipe, error) {
	var pipe models.Pipe
	selectQuery := "SELECT id, cover_photo FROM pipes WHERE user_id=$1 AND id=$2 LIMIT 1"
	err := p.Db.QueryRow(selectQuery, userID, pipeID).Scan(
		&pipe.ID,
		&pipe.CoverPhoto,
	)
	if err != nil {
		return models.Pipe{}, nil
	}

	// Update the pipe
	updateQuery := "UPDATE pipes SET cover_photo=$1, modified_at=now() WHERE id=$2 RETURNING id, name, cover_photo, created_at, modified_at, user_id"
	err = p.Db.QueryRow(updateQuery, pipe.CoverPhoto).Scan(
		&pipe.ID,
		&pipe.Name,
		&pipe.CoverPhoto,
		&pipe.CreatedAt,
		&pipe.ModifiedAt,
		&pipe.UserID,
	)

	return pipe, nil
}

func (p pipeActions) DeletePipe(userID, pipeID int64) (bool, error) {
	deleteQuery := "DELETE FROM pipes WHERE id=$1 AND user_id=$2"
	_, err := p.Db.Exec(deleteQuery, pipeID, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}
