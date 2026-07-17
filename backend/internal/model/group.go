package model

import "time"

type Group struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required,min=3,max=100"`
	Description string    `json:"description" db:"description" binding:"max=255"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type GroupMember struct {
	GroupID  int64     `json:"group_id" db:"group_id"`
	UserID   int64     `json:"user_id" db:"user_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=255"`
}

type AddMemberRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}