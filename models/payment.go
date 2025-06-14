package models

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentModel struct {
	TransactionID string      `json:"transaction_id"` // Unique ID for the transaction
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	Provider      string      `json:"provider"`
	SuccessURL    string      `json:"success_url"`
	CancelURL     string      `json:"cancel_url"`
	WebhookURL    string      `json:"webhook_url"`
	ProductData   ProductData `json:"product_data"`
	Quantity      int64       `json:"quantity"` // For Stripe, this is the quantity of the product
	// Add more fields as needed
}

func (p *PaymentModel) SetTransactionID(id string) {
	p.TransactionID = id
}

func ToPaymentModel(c *gin.Context) (PaymentModel, error) {
	amount, err := strconv.ParseFloat(c.Query("amount"), 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid amount"})
		return PaymentModel{}, err
	}
	return PaymentModel{
		Amount:     amount,
		Currency:   c.Query("currency"),
		Provider:   c.Query("provider"),
		SuccessURL: c.Query("success_url"),
		CancelURL:  c.Query("cancel_url"),
		WebhookURL: c.Query("webhook_url"),
		ProductData: ProductData{
			Name:        c.Query("product_name"),
			Description: c.Query("product_description"),
		},
		Quantity: 1,
	}, nil
}

func (p PaymentModel) ToStripe() StripeModel {
	return StripeModel{
		ProductData: p.ProductData,
		Currency:    p.Currency,
		Amount:      p.Amount,
		Quantity:    p.Quantity,
	}
}

func (c PaymentModel) ToCoinbase() ([]byte, error) {
	req := CoinbaseModel{
		Name:        c.ProductData.Name,
		Description: c.ProductData.Description,
		PricingType: "fixed_price",
	}
	req.LocalPrice.Amount = formatAmount(c.Amount)
	req.LocalPrice.Currency = c.Currency

	return json.Marshal(req)
}

type StripeModel struct {
	ProductData ProductData `json:"product_data"`
	Currency    string      `json:"currency"`
	Amount      float64     `json:"amount"`
	Quantity    int64       `json:"quantity"`
}

type CoinbaseModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PricingType string `json:"pricing_type"`
	LocalPrice  struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"local_price"`
}

type ProductData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Helper to format amount as string with 2 decimals
func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
