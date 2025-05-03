package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type TwilioConfig struct {
	AccountSID        string
	AuthToken         string
	FromContactNumber string
}

type AppConfig struct {
	ServerPort   string
	Dsn          string
	AppSecret    string
	TwilioConfig TwilioConfig
}

func SetUpEnv() (cfg AppConfig, err error) {

	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("HTTP_PORT env not found")
	}

	Dsn := os.Getenv("DSN")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("DSN env not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("App secret env not found")
	}

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if len(twilioAccountSID) < 1 {
		return AppConfig{}, errors.New("Twilio Account SID env not found")
	}

	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if len(twilioAuthToken) < 1 {
		return AppConfig{}, errors.New("Twilio Auth Token env not found")
	}

	twilioFromPhoneNumber := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	if len(twilioFromPhoneNumber) < 1 {
		return AppConfig{}, errors.New("Twilio From Phone Number env not found")
	}

	twilioConfig := TwilioConfig{
		AccountSID:        twilioAccountSID,
		AuthToken:         twilioAuthToken,
		FromContactNumber: twilioFromPhoneNumber,
	}

	return AppConfig{ServerPort: httpPort, Dsn: Dsn, AppSecret: appSecret, TwilioConfig: twilioConfig}, nil

}
