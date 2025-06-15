package models

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentModel struct {
	TransactionID string      `json:"transaction_id"` // Unique ID for the transaction
	ClientID      string      `json:"client_id"`      // Unique ID for the client
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	Provider      string      `json:"provider"`
	PaymentURL    string      `json:"payment_url"`
	SuccessURL    string      `json:"success_url"`
	CancelURL     string      `json:"cancel_url"`
	WebhookURL    string      `json:"webhook_url"`
	ProductData   ProductData `json:"product_data"`
	Quantity      int64       `json:"quantity"` // For Stripe, this is the quantity of the product
	// Add more fields as needed
}

// once transaction id is created, this means the payment is initiated
func (p *PaymentModel) SetTransactionID(id string) {
	p.TransactionID = id
}

func (p *PaymentModel) SetPaymentURL(url string) {
	p.PaymentURL = url
}

func ToPaymentModel(c *gin.Context) (PaymentModel, error) {
	amount, err := strconv.ParseFloat(c.Query("amount"), 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid amount"})
		return PaymentModel{}, err
	}

	// might need to check if client previously initiated a payment without completing.
	// if yes, should we cancel the previous transaction and create a new one?
	// or resuse the same transaction id? and forward the previous created payment link?

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
