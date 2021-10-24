package db

import "gitlab.com/gowagr/mypipe-api/db/model"

func (db Database) CreatePipe(pipe model.Pipe) (model.Pipe, error) {
	var newPipe model.Pipe
	query := "INSERT INTO pipes (user_id, name, cover_photo) VALUES($1, $2, $3) RETURNING name, cover_photo"
	err := db.Conn.QueryRow(query, pipe.UserID, pipe.Name, pipe.CoverPhoto).Scan(
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
	return pipe, nil
}

func (db Database) GetPipes(userID int64) ([]model.Pipe, error) {
	var pipes []model.Pipe
	return pipes, nil
}

func (db Database) UpdatePipe(userID int64, pipe model.Pipe) (model.Pipe, error) {
	return pipe, nil
}
func (db Database) DeletePipe(userID, pipeID int64) (bool, error) {
	return true, nil
}
