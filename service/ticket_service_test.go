package service

import (
	"fix-ticket-system/config"
	"fix-ticket-system/models"
	"fix-ticket-system/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	// Migrate the schema
	err = db.AutoMigrate(&models.Ticket{})
	assert.NoError(t, err)
	return db
}

func setupService(t *testing.T) *TicketService {
	db := setupTestDB(t)
	config.DB = db
	repo := repository.NewTicketRepository()
	return &TicketService{repo: repo}
}

func TestTicketService_CreateTicket(t *testing.T) {
	svc := setupService(t)
	ticket, err := svc.CreateTicket("Title", "Description", "creator@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "Title", ticket.Title)
	assert.Equal(t, "Description", ticket.Description)
	assert.Equal(t, "creator@example.com", ticket.CreatedBy)
}

func TestTicketService_GetTicket(t *testing.T) {
	svc := setupService(t)
	ticket, _ := svc.CreateTicket("Title", "Description", "creator@example.com")
	found, err := svc.GetTicket(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, ticket.ID, found.ID)
	// Not found
	nonExistentID := uuid.New()
	_, err = svc.GetTicket(nonExistentID)
	assert.Error(t, err)
}

func TestTicketService_GetAllTickets(t *testing.T) {
	svc := setupService(t)
	_, _ = svc.CreateTicket("Title1", "Desc1", "a@example.com")
	_, _ = svc.CreateTicket("Title2", "Desc2", "b@example.com")
	tickets, err := svc.GetAllTickets()
	assert.NoError(t, err)
	assert.Len(t, tickets, 2)
}

func TestTicketService_UpdateTicket(t *testing.T) {
	svc := setupService(t)
	ticket, _ := svc.CreateTicket("Title", "Description", "creator@example.com")
	updated, err := svc.UpdateTicket(ticket.ID, "NewTitle", "NewDesc", models.StatusInProgress, models.PriorityHigh, "assignee@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "NewTitle", updated.Title)
	assert.Equal(t, "NewDesc", updated.Description)
	assert.Equal(t, models.StatusInProgress, updated.Status)
	assert.Equal(t, models.PriorityHigh, updated.Priority)
	assert.Equal(t, "assignee@example.com", updated.AssignedTo)
	// Not found
	nonExistentID := uuid.New()
	_, err = svc.UpdateTicket(nonExistentID, "T", "D", models.StatusOpen, models.PriorityLow, "")
	assert.Error(t, err)
}

func TestTicketService_DeleteTicket(t *testing.T) {
	svc := setupService(t)
	ticket, _ := svc.CreateTicket("Title", "Description", "creator@example.com")
	err := svc.DeleteTicket(ticket.ID)
	assert.NoError(t, err)
	// Not found
	err = svc.DeleteTicket(ticket.ID)
	assert.Error(t, err)
}
