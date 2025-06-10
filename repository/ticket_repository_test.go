package repository

import (
	"fix-ticket-system/config"
	"fix-ticket-system/models"
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

func TestTicketRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	config.DB = db
	repo := NewTicketRepository()

	ticket := models.NewTicket("Test Ticket", "Test Description", "test@example.com")

	err := repo.Create(ticket)
	assert.NoError(t, err)

	// Verify the ticket was created
	var found models.Ticket
	err = db.First(&found, "id = ?", ticket.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, ticket.Title, found.Title)
	assert.Equal(t, ticket.Description, found.Description)
	assert.Equal(t, ticket.CreatedBy, found.CreatedBy)
}

func TestTicketRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	config.DB = db
	repo := NewTicketRepository()

	// Create a test ticket
	ticket := models.NewTicket("Test Ticket", "Test Description", "test@example.com")
	err := db.Create(ticket).Error
	assert.NoError(t, err)

	// Test getting existing ticket
	found, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, ticket.ID, found.ID)
	assert.Equal(t, ticket.Title, found.Title)

	// Test getting non-existent ticket
	nonExistentID := uuid.New()
	_, err = repo.GetByID(nonExistentID)
	assert.Error(t, err)
}

func TestTicketRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	config.DB = db
	repo := NewTicketRepository()

	// Create test tickets
	tickets := []*models.Ticket{
		models.NewTicket("Ticket 1", "Description 1", "test1@example.com"),
		models.NewTicket("Ticket 2", "Description 2", "test2@example.com"),
	}

	for _, ticket := range tickets {
		err := db.Create(ticket).Error
		assert.NoError(t, err)
	}

	// Test getting all tickets
	found, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, found, 2)
}

func TestTicketRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	config.DB = db
	repo := NewTicketRepository()

	// Create a test ticket
	ticket := models.NewTicket("Test Ticket", "Test Description", "test@example.com")
	err := db.Create(ticket).Error
	assert.NoError(t, err)

	// Update the ticket
	ticket.Title = "Updated Title"
	ticket.Description = "Updated Description"
	err = repo.Update(ticket)
	assert.NoError(t, err)

	// Verify the update
	var found models.Ticket
	err = db.First(&found, "id = ?", ticket.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", found.Title)
	assert.Equal(t, "Updated Description", found.Description)
}

func TestTicketRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	config.DB = db
	repo := NewTicketRepository()

	// Create a test ticket
	ticket := models.NewTicket("Test Ticket", "Test Description", "test@example.com")
	err := db.Create(ticket).Error
	assert.NoError(t, err)

	// Delete the ticket
	err = repo.Delete(ticket.ID)
	assert.NoError(t, err)

	// Verify the ticket was deleted
	var found models.Ticket
	err = db.First(&found, "id = ?", ticket.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
