package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BitPayInvoiceRequest struct {
	Price           float64 `json:"price"`
	Currency        string  `json:"currency"`
	RedirectURL     string  `json:"redirectURL"`
	NotificationURL string  `json:"notificationURL"`
}

func CreateBitPayInvoice(c *gin.Context) {
	body := BitPayInvoiceRequest{
		Price:           10.00,
		Currency:        "USD",
		RedirectURL:     "https://yourdomain.com/payment/success",
		NotificationURL: "https://yourdomain.com/api/payment/webhook/bitpay",
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://bitpay.com/invoices", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept-Version", "2.0.0")
	req.Header.Set("Authorization", "Bearer YOUR_BITPAY_TOKEN")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact BitPay"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if invoice, ok := result["data"].(map[string]interface{}); ok {
		c.JSON(http.StatusOK, gin.H{"url": invoice["url"]})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid BitPay response"})
	}
}

func BitPayWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook"})
		return
	}

	// Handle logic (e.g., check status: complete/confirmed)
	c.Status(http.StatusOK)
}
