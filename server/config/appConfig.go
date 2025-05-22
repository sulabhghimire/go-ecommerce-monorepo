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

type StripeConfig struct {
	StripeSecretKey string
	SuccessUrl      string
	CancelUrl       string
	PublishableKey  string
}

type AppConfig struct {
	ServerPort   string
	Dsn          string
	AppSecret    string
	TwilioConfig TwilioConfig
	StripeConfig StripeConfig
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
		return AppConfig{}, errors.New("dsn env not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("app secret env not found")
	}

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if len(twilioAccountSID) < 1 {
		return AppConfig{}, errors.New("twilio Account SID env not found")
	}

	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if len(twilioAuthToken) < 1 {
		return AppConfig{}, errors.New("twilio Auth Token env not found")
	}

	twilioFromPhoneNumber := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	if len(twilioFromPhoneNumber) < 1 {
		return AppConfig{}, errors.New("twilio From Phone Number env not found")
	}

	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if len(stripeSecretKey) < 1 {
		return AppConfig{}, errors.New("stripe secret key env not found")
	}

	stripeSuccessUrl := os.Getenv("STRIPE_SUCCESS_URL")
	if len(stripeSuccessUrl) < 1 {
		return AppConfig{}, errors.New("stripe success url env not found")
	}

	stripeCancelUrl := os.Getenv("STRIPE_CANCEL_URL")
	if len(stripeCancelUrl) < 1 {
		return AppConfig{}, errors.New("stripe cancel url env not found")
	}

	publishableKey := os.Getenv("STRIPE_PUB_KEY")
	if len(publishableKey) < 1 {
		return AppConfig{}, errors.New("stripe publishable key env not found")
	}

	twilioConfig := TwilioConfig{
		AccountSID:        twilioAccountSID,
		AuthToken:         twilioAuthToken,
		FromContactNumber: twilioFromPhoneNumber,
	}

	stripeConfig := StripeConfig{
		StripeSecretKey: stripeSecretKey,
		SuccessUrl:      stripeSuccessUrl,
		CancelUrl:       stripeCancelUrl,
		PublishableKey:  publishableKey,
	}

	return AppConfig{ServerPort: httpPort, Dsn: Dsn, AppSecret: appSecret, TwilioConfig: twilioConfig, StripeConfig: stripeConfig}, nil

}
