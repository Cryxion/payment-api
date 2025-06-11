package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateCoinbaseCharge(c *gin.Context) {

	url := "https://api.commerce.coinbase.com/charges"
	method := "POST"

	//json body for http
	jsonBody := `{
      "name": "The Human Fund",
      "description": "Money For People",
      "pricing_type": "fixed_price",
      "local_price": {
        "amount": "1.00",
        "currency": "USD"
      }
    }`

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(jsonBody))

	if err != nil {
		fmt.Println(err)
		return
	}

	// get from .env
	req.Header.Add("X-CC-Api-Key", os.Getenv("COINBASE_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// convert body to json and aaccess hosted_url at root
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	// if empty or any status other tahn 200, return error
	if response == nil || res.StatusCode != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create charge"})
	}
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response structure"})
		return
	}
	hostedURL, ok := data["hosted_url"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hosted_url not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": hostedURL})
}

// TODO: Implement Coinbase webhook handler
func CoinbaseWebhook(c *gin.Context) {

}
