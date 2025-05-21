package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pId string) (*stripe.CheckoutSession, error)
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	faliureUrl      string
}

// CreatePayment implements PaymentClient.
func (p *payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := int64(amount * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(p.successUrl),
		CancelURL:          stripe.String(p.faliureUrl),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("usd"),
					UnitAmount: stripe.Int64(amountInCents),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Electronics"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))
	params.AddMetadata("order_id", orderId)

	session, err := session.New(params)
	if err != nil {
		log.Printf("Error creating Stripe session: %v", err)
		return nil, errors.New("payment creation failed")
	}

	// Log or return session.URL if you want to redirect user
	log.Printf("Stripe session URL: %s", session.URL)

	return session, nil
}

// GetPaymentStatus implements PaymentClient.
func (p *payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {

	stripe.Key = p.stripeSecretKey

	session, err := session.Get(pId, nil)
	if err != nil {
		log.Printf("Error retrieving Stripe session: %v", err)
		return nil, errors.New("payment status retrieval failed")
	}

	return session, nil

}

func NewPaymentClient(stripeSecretKey, successUrl, faliureUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		faliureUrl:      faliureUrl,
	}
}
