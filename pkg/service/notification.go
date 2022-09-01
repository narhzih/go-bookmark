package service

import (
	"github.com/appleboy/go-fcm"
	"log"
	"os"
)

func (s Service) CreateTwitterPipeShareNotification(tweetUrl, pipeName string, userId int64) error {
	message := "The following tweet has been successfully saved to " + pipeName + ": " + tweetUrl
	_, err := s.DB.CreateNotification(userId, message)
	if err != nil {
		return err
	}
	return nil

}
func (s Service) CreatePrivatePipeShareNotification(sharedPipeId, sharerId, sharedToId int64) error {
	sharedPipe, err := s.DB.GetPipe(sharedPipeId, sharerId)
	if err != nil {
		return err
	}
	sharer, err := s.DB.GetUserById(int(sharerId))
	if err != nil {
		return err
	}
	message := sharer.ProfileName + " privately shared you pipe with name: " + sharedPipe.Name
	_, err = s.DB.CreateNotification(sharedToId, message)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) SendPushNotification(message string, deviceTokens []string) error {
	msg := &fcm.Message{
		To: deviceTokens[0],
		Data: map[string]interface{}{
			"message": message,
		},
		Notification: &fcm.Notification{
			Title: "My pipe notification",
			Body:  "Notification body",
		},
	}

	client, err := fcm.NewClient(os.Getenv("GFCM_SERVER_KEY"))
	if err != nil {
		return err
	}
	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
	return nil
}
