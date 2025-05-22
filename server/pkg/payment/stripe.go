package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(pId string) (*stripe.PaymentIntent, error)
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	faliureUrl      string
}

// CreatePayment implements PaymentClient.
func (p *payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amountInCents),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))
	params.AddMetadata("order_id", orderId)

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating payment intent: %v", err)
		return nil, errors.New("payment creation failed")
	}

	// Log or return session.URL if you want to redirect user
	return pi, nil
}

// GetPaymentStatus implements PaymentClient.
func (p *payment) GetPaymentStatus(pId string) (*stripe.PaymentIntent, error) {

	stripe.Key = p.stripeSecretKey
	params := &stripe.PaymentIntentParams{}

	result, err := paymentintent.Get(pId, params)
	if err != nil {
		log.Printf("Error retrieving Stripe session: %v", err)
		return nil, errors.New("error fetching payment status")
	}

	return result, nil

}

func NewPaymentClient(stripeSecretKey, successUrl, faliureUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		faliureUrl:      faliureUrl,
	}
}
