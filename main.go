package main

import (
	"log"
	"paymentbe/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	r := gin.Default()
	api := r.Group("/api")
	{
		payment := api.Group("/payment")
		{
			payment.POST("/session", routes.CreatePaymentSession)
			payment.POST("/webhook/:type", routes.HandleWebhook)
		}
	}

	// Optional redirect handlers
	r.GET("/payment/success", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Payment Success"})
	})

	r.GET("/payment/cancel", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Payment Cancelled"})
	})

	r.Run(":8080")
}
