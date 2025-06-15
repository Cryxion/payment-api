package handlers

import (
	"fmt"
	"net/http"
	"os"
	"paymentbe/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)

// CreateStripeSession handles creation of a Stripe Checkout Session
func CreateStripeSession(c *gin.Context) {
	paymentModel, err := models.ToPaymentModel(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	stripeModel := paymentModel.ToStripe()

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY") // Use Stripe sandbox/test secret key
	sessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card", "link", "alipay", "grabpay"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(stripeModel.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(stripeModel.ProductData.Name),
						Description: stripe.String(stripeModel.ProductData.Description),
					},
					UnitAmount: stripe.Int64(int64(stripeModel.Amount * 100)), // Convert to cents
				},
				Quantity: stripe.Int64(stripeModel.Quantity),
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

	// save s.ID if not empty
	if s.ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Stripe session ID is empty"})
		return
	}
	paymentModel.SetTransactionID(s.ID)

	paymentModel.SetPaymentURL(s.URL)
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

	eventType := string(event.Type)

	// save payload and event to file
	err = os.WriteFile(eventType+".json", payload, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Example: handle successful payment session
	if event.Type == "checkout.session.completed" {
		stripeSessionID, ok := event.Data.Object["id"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}
		println(stripeSessionID)
		// TODO: find paymentmodel by transaction ID of stripeSessionID, then mark as paid

		// TODO: confirm order, notify user, update DB, etc.
		c.JSON(http.StatusOK, gin.H{"status": "payment succeeded"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "event received"})
}
