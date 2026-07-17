package model

import "time"

type Payment struct {
	ID            int64     `json:"id" db:"id"`
	GroupID       int64     `json:"group_id" db:"group_id"`
	PayerID       int64     `json:"payer_id" db:"payer_id"`
	PayeeID       int64     `json:"payee_id" db:"payee_id"`
	Amount        int64     `json:"amount" db:"amount"`
	Notes         string    `json:"notes" db:"notes"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	CreatedBy     int64     `json:"created_by" db:"created_by"`
}

type RecordPaymentRequest struct {
	PayerID       int64  `json:"payer_id" binding:"required"`
	PayeeID       int64  `json:"payee_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Notes         string `json:"notes" binding:"max=255"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=CASH UPI STRIPE"`
}
