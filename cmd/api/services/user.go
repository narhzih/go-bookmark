package services

import (
	"gitlab.com/trencetech/mypipe-api/cmd/api/helpers"
	"gitlab.com/trencetech/mypipe-api/db/actions/postgres"
	"gitlab.com/trencetech/mypipe-api/db/models"
)

func (s Services) GetUserProfileInformation(userID int64) (models.Profile, error) {
	var profile models.Profile
	var err error

	profile.User, err = s.Repositories.User.GetUserById(int(userID))
	if err != nil {
		s.Logger.Err(err).Msg("Error is from getting user by ID")
		s.Logger.Err(err).Msg(err.Error())
		return models.Profile{}, err
	}
	profile.Pipes, err = s.Repositories.Pipe.GetPipesCount(userID)
	if err != nil {
		s.Logger.Err(err).Msg("Error is from getting user pipe counts")
		s.Logger.Err(err).Msg(err.Error())
		return profile, err
	}
	profile.Bookmarks, err = s.Repositories.Bookmark.GetBookmarksCount(userID)
	if err != nil {
		s.Logger.Err(err).Msg("Error is from getting user ID")
		s.Logger.Err(err).Msg(err.Error())
		return profile, err
	}

	return profile, nil
}

func (s Services) UserWithUsernameExists(username string) (bool, error) {
	exits := false

	_, err := s.Repositories.User.GetUserByUsername(username)
	if err != nil {
		if err == postgres.ErrNoRecord {
			// This means that there's no record with that user
			// and we're good to go
			exits = true
			return exits, nil
		}

		return exits, err
	}

	return exits, nil
}

func (s Services) UserWithUsernameExistsWithUser(username string) (bool, error) {
	exits := false

	_, err := s.Repositories.User.GetUserByUsername(username)
	if err != nil {
		if err == postgres.ErrNoRecord {
			// This means that there's no record with that user
			// and we're good to go
			exits = true
			return exits, nil
		}

		return exits, err
	}

	return exits, nil
}

func (s Services) MarkUserAsVerified(user models.User, token string) (models.User, error) {
	var err error
	user, err = s.Repositories.User.VerifyUser(user)
	if err != nil {
		return models.User{}, err
	}
	_, err = s.Repositories.AccountVerification.DeleteVerification(token)
	if err != nil {
		s.Logger.Err(err).Msg("Could not delete verification token from db")
	}
	return user, nil
}

func (s Services) TokenInUserDeviceTokens(userID int64, deviceToken string) (bool, error) {
	userDeviceTokens, err := s.Repositories.User.GetUserDeviceTokens(userID)
	if err != nil {
		return false, err
	}
	return helpers.SliceContains(userDeviceTokens, deviceToken), nil
}

func (s Services) TwitterAccountConnected(twitterID string) (models.User, error) {
	user, err := s.Repositories.User.GetUserByTwitterID(twitterID)
	if err != nil {
		if err == postgres.ErrNoRecord {
			return models.User{}, postgres.ErrNoRecord
		}
		return models.User{}, err
	}
	return user, nil
}
