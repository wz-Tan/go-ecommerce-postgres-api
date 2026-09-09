package orders

// Expected Data from POST Request
// Needs capital letters to be exported
type OrderPayload struct {
	CustomerID int64 `json:"customer_id"`
	ProductID  int64 `json:"product_id"`
	Quantity   int32 `json:"quantity"`
	PriceCents int32 `json:"price_cents"`
}
