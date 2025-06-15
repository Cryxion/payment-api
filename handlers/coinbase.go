package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"paymentbe/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateCoinbaseCharge(c *gin.Context) {
	paymentModel, err := models.ToPaymentModel(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	coinbaseModel, err := paymentModel.ToCoinbase()
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	url := os.Getenv("COINBASE_API_URL")
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(coinbaseModel)))

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

	// save transaction id to db
	transactionID, ok := data["id"].(string)
	paymentModel.SetTransactionID(transactionID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction_id not found"})
		return
	}

	hostedURL, ok := data["hosted_url"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hosted_url not found"})
		return
	}
	paymentModel.SetPaymentURL(hostedURL)
	c.JSON(http.StatusOK, gin.H{"url": hostedURL})
}

// TODO: Implement Coinbase webhook handler
func CoinbaseWebhook(c *gin.Context) {

	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Reading body failed"})
		return
	}

	sigHeader := c.GetHeader("X-CC-WEBHOOK-SIGNATURE")

	// verify sigHeader with environment variable
	endpointSecret := os.Getenv("COINBASE_WEBHOOK_SECRET")
	if sigHeader != endpointSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//convert payload to json
	var event map[string]interface{}
	json.Unmarshal(payload, &event)

	// the id is kept at event["event"]["data"]["id"]
	if event == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
		return
	}

	eventData, ok := event["event"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
		return
	}

	data, ok := eventData["data"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
		return
	}

	chargeID, ok := data["id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
		return
	}

	println("Coinbase Charge ID:", chargeID)
	// TODO: find paymentmodel by transaction ID of chargeID, then mark as paid

}
