package auth

import (
	"fmt"
	"os"
	"paymentbe/models/errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateToken(c *gin.Context) {
	// Implementation for creating a token

	// TODO: fetch credentials from request and validate them with DB

	// mock_credentials
	username := c.PostForm("username")
	password := c.PostForm("password")
	mock_credentials := map[string]string{
		"username": username,
		"password": password,
	}

	// auth mock_credentials
	if mock_credentials["username"] != "testuser" || mock_credentials["password"] != "testpass" {
		c.JSON(401, gin.H{"error": errors.ErrInvalidCredential})
		return
	}

	token := GenerateToken(mock_credentials)

	c.JSON(200, gin.H{
		"token": token,
	})

}

func GenerateToken(credentials map[string]string) string {
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	// Create the JWT claims, which includes the username and expiry time
	claims := jwt.MapClaims{
		"username": credentials["username"],
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}

	return tokenString
}

// ValidateToken checks if the provided token is valid
func ValidateToken(tokenString string) (*jwt.Token, error) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	// Parse the token
	tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix if present

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
