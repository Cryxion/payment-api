package routes

import (
	"paymentbe/auth"
	"paymentbe/handlers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(c *gin.Context) {
	auth.CreateToken(c)
}

func CreatePaymentSession(c *gin.Context) {
	provider := c.Query("provider")
	switch provider {
	case "stripe": // TODO: webhook
		handlers.CreateStripeSession(c)
	case "paypal": //TODO: add support
		handlers.CreatePayPalOrder(c)
	case "razorpay": //TODO: add support
		handlers.CreateRazorpayOrder(c)
	case "bitpay": // TODO: add support, due to initial confirmation and deposit require for setting up
		handlers.CreateBitPayInvoice(c)
	case "coinbase": // TODO: webhook
		handlers.CreateCoinbaseCharge(c)
	default:
		c.JSON(400, gin.H{"error": "Unsupported provider"})
	}
}

// TOD: enhance this.
func HandleWebhook(c *gin.Context) {
	provider := c.Param("type")
	switch provider {
	case "stripe":
		handlers.StripeWebhook(c)
	case "paypal":
		handlers.PayPalWebhook(c)
	case "razorpay":
		handlers.RazorpayWebhook(c)
	default:
		c.JSON(400, gin.H{"error": "Unsupported webhook type"})
	}
}
