package service

import (
	"github.com/appleboy/go-fcm"
	"log"
	"os"
)

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
