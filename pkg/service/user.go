package service

import (
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/db/model"
	"gitlab.com/trencetech/mypipe-api/pkg/helpers"
)

func (s Service) GetUserProfileInformation(userID int64) (model.Profile, error) {
	var profile model.Profile
	var err error

	profile.User, err = s.DB.GetUserById(int(userID))
	if err != nil {
		s.DB.Logger.Err(err).Msg("Error is from getting user by ID")
		s.DB.Logger.Err(err).Msg(err.Error())
		return model.Profile{}, err
	}
	profile.Pipes, err = s.DB.GetPipesCount(userID)
	if err != nil {
		s.DB.Logger.Err(err).Msg("Error is from getting user pipe counts")
		s.DB.Logger.Err(err).Msg(err.Error())
		return profile, err
	}
	profile.Bookmarks, err = s.DB.GetBookmarksCount(userID)
	if err != nil {
		s.DB.Logger.Err(err).Msg("Error is from getting user ID")
		s.DB.Logger.Err(err).Msg(err.Error())
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

func (s Service) UserWithUsernameExistsWithUser(username string) (bool, error) {
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

func (s Service) MarkUserAsVerified(user model.User, token string) (model.User, error) {
	var err error
	user, err = s.DB.VerifyUser(user)
	if err != nil {
		return model.User{}, err
	}
	_, err = s.DB.DeleteVerification(token)
	if err != nil {
		s.DB.Logger.Err(err).Msg("Could not delete verification token from db")
	}
	return user, nil
}

func (s Service) TokenInUserDeviceTokens(userID int64, deviceToken string) (bool, error) {
	userDeviceTokens, err := s.DB.GetUserDeviceTokens(userID)
	if err != nil {
		return false, err
	}
	return helpers.SliceContains(userDeviceTokens, deviceToken), nil
}

func (s Service) TwitterAccountConnected(twitterID string) (model.User, error) {
	user, err := s.DB.GetUserByTwitterID(twitterID)
	if err != nil {
		if err == db.ErrNoRecord {
			return model.User{}, db.ErrNoRecord
		}
		return model.User{}, err
	}
	return user, nil
}
