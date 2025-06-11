package handlers

import (
	"github.com/gin-gonic/gin"
)

func CreateRazorpayOrder(c *gin.Context) {
	// Use Razorpay API to generate order_id
	// Return order & redirect URL if needed
}
func RazorpayWebhook(c *gin.Context) {
	// Parse & verify Razorpay webhook
}
