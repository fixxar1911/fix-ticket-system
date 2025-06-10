package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTicket(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		createdBy   string
		wantErr     bool
	}{
		{
			name:        "Valid ticket",
			title:       "Test Ticket",
			description: "Test Description",
			createdBy:   "test@example.com",
			wantErr:     false,
		},
		{
			name:        "Empty title",
			title:       "",
			description: "Test Description",
			createdBy:   "test@example.com",
			wantErr:     false, // NewTicket doesn't validate empty title
		},
		{
			name:        "Empty description",
			title:       "Test Ticket",
			description: "",
			createdBy:   "test@example.com",
			wantErr:     false, // NewTicket doesn't validate empty description
		},
		{
			name:        "Empty createdBy",
			title:       "Test Ticket",
			description: "Test Description",
			createdBy:   "",
			wantErr:     false, // NewTicket doesn't validate empty createdBy
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket := NewTicket(tt.title, tt.description, tt.createdBy)

			if ticket == nil {
				t.Error("NewTicket() returned nil")
				return
			}

			if ticket.Title != tt.title {
				t.Errorf("NewTicket() title = %v, want %v", ticket.Title, tt.title)
			}
			if ticket.Description != tt.description {
				t.Errorf("NewTicket() description = %v, want %v", ticket.Description, tt.description)
			}
			if ticket.CreatedBy != tt.createdBy {
				t.Errorf("NewTicket() createdBy = %v, want %v", ticket.CreatedBy, tt.createdBy)
			}
			if ticket.Status != StatusOpen {
				t.Errorf("NewTicket() status = %v, want %v", ticket.Status, StatusOpen)
			}
			if ticket.Priority != PriorityMedium {
				t.Errorf("NewTicket() priority = %v, want %v", ticket.Priority, PriorityMedium)
			}
			if ticket.CreatedAt.IsZero() {
				t.Error("NewTicket() CreatedAt is zero")
			}
			if ticket.UpdatedAt.IsZero() {
				t.Error("NewTicket() UpdatedAt is zero")
			}
			if ticket.ID == uuid.Nil {
				t.Error("NewTicket() ID is nil")
			}
		})
	}
}
