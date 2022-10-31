package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/models"
	"google.golang.org/api/option"
	"os"
)

func (s Services) CreateTwitterPipeShareNotification(tweetUrl, pipeName string, userId int64) error {
	message := "The following tweet has been successfully saved to " + pipeName + ": " + tweetUrl
	_, err := s.Repositories.Notification.CreateNotification(userId, message, "")
	if err != nil {
		return err
	}
	return nil

}
func (s Services) CreatePrivatePipeShareNotification(sharedPipeId, sharerId, sharedToId int64) error {
	sharedPipe, err := s.Repositories.Pipe.GetPipe(sharedPipeId, sharerId)
	if err != nil {
		return err
	}
	sharer, err := s.Repositories.User.GetUserById(sharerId)
	if err != nil {
		return err
	}
	pipeShareRecord, err := s.Repositories.PipeShare.GetSharedPipe(sharedPipeId)
	if err != nil {
		return err
	}
	metadata := models.MDPrivatePipeShare{
		Pipe:   sharedPipe,
		Sharer: sharer,
		Code:   pipeShareRecord.Code,
	}
	mdToJson, _ := json.Marshal(metadata)
	message := sharer.ProfileName + " privately shared you pipe with name: " + sharedPipe.Name
	_, err = s.Repositories.Notification.CreateNotification(sharedToId, message, string(mdToJson))
	if err != nil {
		return err
	}

	// try to send push notification
	userDeviceTokens, err := s.Repositories.User.GetUserDeviceTokens(sharedToId)
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
	decodedKey, err := getDecodedFireBaseKey()
	if err != nil {
		return err
	}
	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}
	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		return err
	}

	firebaseClient, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}
	//msg := &fcm.Message{
	//	To: deviceTokens[0],
	//	Data: map[string]interface{}{
	//		"message": message,
	//	},
	//	Notification: &fcm.Notification{
	//		Title: "My pipe notification",
	//		Body:  "Notification body",
	//	},
	//}
	//client, err := fcm.NewClient(os.Getenv("GFCM_SERVER_KEY"))
	//if err != nil {
	//	return err
	//}
	//response, err := client.Send(msg)
	response, err := firebaseClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Tokens: deviceTokens,
	})

	if err != nil {
		return err
	}

	s.Logger.Info().Msg(fmt.Sprintf("%#v\n", response))
	return nil
}

func getDecodedFireBaseKey() ([]byte, error) {

	fireBaseAuthKey := os.Getenv("GFCM_SERVER_KEY")

	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}
