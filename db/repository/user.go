package repository

import "github.com/mypipeapp/mypipeapi/db/models"

type UserRepository interface {
	CreateUserByEmail(user models.User, password, authOrigin string) (newUser models.User, err error)
	GetUserByTwitterID(twitterId string) (user models.User, err error)
	GetUserById(userId int64) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
	GetUserByEmail(userEmail string) (user models.User, err error)
	GetUserAndAuth(userId int64) (models.UserAuth, error)
	UpdateUserPassword(userId int64, password string) error
	UpdateUser(updatedBody models.User) (models.User, error)
	VerifyUser(user models.User) (models.User, error)
	GetUserDeviceTokens(userId int64) ([]string, error)
	UpdateUserDeviceTokens(userId int64, deviceTokens []string) ([]string, error)
	ConnectToTwitter(user models.User, twitterId string) (models.User, error)
	DisconnectTwitter(user models.User) (models.User, error)
}
