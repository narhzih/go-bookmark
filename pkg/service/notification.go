package service

import (
	"github.com/appleboy/go-fcm"
	"log"
)

func (s Service) SendPushNotification(message string, deviceToken string) error {
	msg := &fcm.Message{
		To: "sample_device_token",
		Data: map[string]interface{}{
			"message": message,
		},
		Notification: &fcm.Notification{
			Title: "My pipe notification",
			Body:  "Notification body",
		},
	}

	client, err := fcm.NewClient("sample_api_key")
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
