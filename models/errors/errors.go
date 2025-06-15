package errors

const (
	ErrWebhookUnauthorized       = "Unauthorized webhook request"
	ErrInvalidEventData          = "Invalid event data"
	ErrReadingBodyFailed         = "Reading body failed"
	ErrTransactionIDNotFound     = "Transaction ID not found"
	ErrPaymentURLNotFound        = "Payment URL not found"
	ErrFailedToCreateTransaction = "Failed to create transaction"
	ErrInvalidTransactionData    = "Invalid transaction data"
	ErrInvalidCredential         = "Invalid credential"
	ErrTokenNotFound             = "Token not found"
	ErrTokenNotValid             = "Token not valid"
	// Add more as needed
)
