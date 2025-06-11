package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/webhook"
)

// CreateStripeSession handles creation of a Stripe Checkout Session
func CreateStripeSession(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY") // Use Stripe sandbox/test secret key
	sessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Test Item"),
					},
					UnitAmount: stripe.Int64(1500), // $15.00
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://yourdomain.com/payment/success"),
		CancelURL:  stripe.String("https://yourdomain.com/payment/cancel"),
	}

	s, err := session.New(sessionParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}

// TODO: complete webhook handling for Stripe
// StripeWebhook handles incoming webhook events from Stripe
func StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Reading body failed"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	// endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook signature verification failed"})
		return
	}

	// Example: handle successful payment session
	if event.Type == "checkout.session.completed" {
		// TODO: confirm order, notify user, update DB, etc.
		c.JSON(http.StatusOK, gin.H{"status": "payment succeeded"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "event received"})
}
