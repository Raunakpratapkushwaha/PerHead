package model

import "time"

type SplitType string

const (
	SplitEqual      SplitType = "EQUAL"
	SplitExact      SplitType = "EXACT"
	SplitPercentage SplitType = "PERCENTAGE"
	SplitShares     SplitType = "SHARES"
)

type Expense struct {
	ID          int64     `json:"id" db:"id"`
	GroupID     int64     `json:"group_id" db:"group_id"`
	PayerID     int64     `json:"payer_id" db:"payer_id"`
	Amount      int64     `json:"amount" db:"amount"` // in cents
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	SplitType   SplitType `json:"split_type" db:"split_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
}

type ExpenseSplit struct {
	ID        int64   `json:"id" db:"id"`
	ExpenseID int64   `json:"expense_id" db:"expense_id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Amount    int64   `json:"amount" db:"amount"` // in cents
	Percentage float64 `json:"percentage,omitempty" db:"percentage"`
	Share      int     `json:"share,omitempty" db:"share"`
}

// SplitInput is what handlers parse from incoming HTTP POST payloads
type SplitInput struct {
	UserID     int64   `json:"user_id" binding:"required"`
	Amount     int64   `json:"amount,omitempty"` // Used for EXACT
	Percentage float64 `json:"percentage,omitempty"` // Used for PERCENTAGE
	Share      int     `json:"share,omitempty"`      // Used for SHARES
}

type CreateExpenseRequest struct {
	GroupID     int64        `json:"group_id" binding:"required"`
	PayerID     int64        `json:"payer_id" binding:"required"`
	Amount      int64        `json:"amount" binding:"required,gt=0"` // in cents
	Description string       `json:"description" binding:"required,max=255"`
	Category    string       `json:"category"`
	SplitType   SplitType    `json:"split_type" binding:"required,oneof=EQUAL EXACT PERCENTAGE SHARES"`
	Splits      []SplitInput `json:"splits" binding:"required,dive"`
}