package service

import (
	"gitlab.com/gowagr/mypipe-api/db"
	"gitlab.com/gowagr/mypipe-api/db/model"
)

func (s Service) GetUserProfileInformation(userID int64) (model.Profile, error) {
	var profile model.Profile
	var err error

	profile.User, err = s.DB.GetUserById(int(userID))
	if err != nil {
		return model.Profile{}, err
	}
	profile.Pipes, err = s.DB.GetPipesCount(userID)
	if err != nil {
		return profile, err
	}
	profile.Bookmarks, err = s.DB.GetBookmarksCount(userID)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s Service) UserWithUsernameExists(username string) (bool, error) {
	exits := false

	_, err := s.DB.GetUserByUsername(username)
	if err != nil {
		if err == db.ErrNoRecord {
			// This means that there's no record with that user
			// and we're good to go
			exits = true
			return exits, nil
		}

		return exits, err
	}

	return exits, nil
}
