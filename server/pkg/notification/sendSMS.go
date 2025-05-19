package notification

import (
	"ecommerce/config"
	"encoding/json"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type NotificationClient interface {
	SendSMS(phone, message string) error
}

type notificationClient struct {
	config config.AppConfig
}

func NewNotificationClient(config config.AppConfig) NotificationClient {
	return &notificationClient{
		config: config,
	}
}

func (c notificationClient) SendSMS(phone, message string) error {

	fmt.Println(c.config.TwilioConfig)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: c.config.TwilioConfig.AccountSID,
		Password: c.config.TwilioConfig.AuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(c.config.TwilioConfig.FromContactNumber)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
		return err
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

	return nil
}
