package model

// Settlement represents a simplified transaction to clear group debt
type Settlement struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
	Amount     int64 `json:"amount"` // in cents
}