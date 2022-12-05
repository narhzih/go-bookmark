package services

import (
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/models"
	"google.golang.org/api/option"
	"path/filepath"
)

func (s Services) CreateTwitterPipeShareNotification(tweetUrl, pipeName string, userId int64) error {
	message := "The following tweet has been successfully saved to " + pipeName + ": " + tweetUrl
	_, err := s.Repositories.Notification.CreateNotification(userId, message, "")
	if err != nil {
		return err
	}
	return nil

}

func (s Services) CreatePrivatePipeShareNotification(sharedPipeCode string, sharedPipeId, sharerId, sharedToId int64) error {
	sharedPipe, err := s.Repositories.Pipe.GetPipe(sharedPipeId, sharerId)
	if err != nil {
		return err
	}
	sharer, err := s.Repositories.User.GetUserById(sharerId)
	if err != nil {
		return err
	}
	metadata := models.MDPrivatePipeShare{
		Pipe:   sharedPipe,
		Sharer: sharer,
		Code:   sharedPipeCode,
	}
	mdToJson, _ := json.Marshal(metadata)
	message := sharer.Username + " shared a pipe with you"
	_, err = s.Repositories.Notification.CreateNotification(sharedToId, message, string(mdToJson))
	if err != nil {
		return err
	}

	// try to send push notification
	userDeviceTokens, err := s.Repositories.User.GetUserDeviceTokens(sharedToId)
	s.Logger.Info().Msg(fmt.Sprintf("user device tokens -> %+v", userDeviceTokens))
	switch {
	case errors.Is(err, nil):
		pnErr := s.SendPushNotification("Pipe share", message, userDeviceTokens)
		if pnErr != nil {
			s.Logger.Err(err).Msg("An error occurred while sending push notification")
		}
	default:
		s.Logger.Err(err).Msg("An error occurred but not while sending push notifications")
	}
	return nil
}

func (s Services) SendPushNotification(title, message string, deviceTokens []string) error {
	serviceAccountKeyFilePath, err := filepath.Abs("./fbServiceAccount.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opts := option.WithCredentialsFile(serviceAccountKeyFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opts)
	if err != nil {
		return err
	}

	firebaseClient, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	response, err := firebaseClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Tokens: deviceTokens,
	})

	if err != nil {
		s.Logger.Err(err).Msg("An error occurred while sending push notification")
		return err
	}

	s.Logger.Info().Msg(fmt.Sprintf("%#v\n", response))
	s.Logger.Info().Msg(fmt.Sprintf("successfully sent push notification..."))
	return nil
}
