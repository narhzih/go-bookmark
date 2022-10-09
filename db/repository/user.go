package repository

import "gitlab.com/trencetech/mypipe-api/db/models"

type UserRepository interface {
	CreateUserByEmail(user models.User, password, authOrigin string) (newUser models.User, err error)
	CreateUser(user models.User) (newUser models.User, err error)
	GetUserByTwitterID(twitterId string) (user models.User, err error)
	GetUserById(userId int) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
	GetUserByEmail(userEmail string) (user models.User, err error)
	GetUserAndAuth(user models.User) (userAndAuth models.UserAuth, err error)
	UpdateUserPassword(userId int, password string) error
	UpdateUser(updatedBody models.User) (models.User, error)
	VerifyUser(user models.User) (models.User, error)
	GetUserDeviceTokens(userId int64) ([]string, error)
	UpdateUserDeviceTokens(userId int64, deviceTokens []string) ([]string, error)
	ConnectToTwitter(user models.User, twitterId string) (models.User, error)
}
