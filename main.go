package main

import (
	"log"
	"net/http"
	"paymentbe/auth"
	"paymentbe/routes"

	"paymentbe/models/errors"

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

		authorize := api.Group("/auth")
		{
			authorize.POST("/token", routes.AuthenticationRoutes)
		}

		//validate JWT token from header before accessing payment routes
		payment := api.Group("/payment", func(c *gin.Context) {
			//bearer token authentication
			token := c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrTokenNotFound})
				c.Abort()
				return
			} else {
				// Validate the token
				_, err := auth.ValidateToken(token)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrTokenNotValid})
					c.Abort()
					return
				}
			}

		})
		{
			payment.POST("/session", routes.CreatePaymentSession)
		}

		api.POST("/payment/webhook/:type", routes.HandleWebhook)

	}

	// Optional redirect handlers, this is for went the page complete payment sucess or cancel,
	// probably here will also check with DB if payment is competed/cancelled, and redirect to frontend UI
	r.GET("/payment/success", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Payment Success"})
	})

	r.GET("/payment/cancel", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Payment Cancelled"})
	})

	r.Run(":8080")
}
