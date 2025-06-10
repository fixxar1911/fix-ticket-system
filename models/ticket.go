package models

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the current state of a ticket
type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

// Priority represents the priority level of a ticket
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Ticket represents a support ticket in the system
type Ticket struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Status      Status    `json:"status" gorm:"type:varchar(20);not null;default:'open'"`
	Priority    Priority  `json:"priority" gorm:"type:varchar(20);not null;default:'medium'"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
	AssignedTo  string    `json:"assigned_to"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

// NewTicket creates a new ticket with default values
func NewTicket(title, description, createdBy string) *Ticket {
	now := time.Now()
	return &Ticket{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		Priority:    PriorityMedium,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
