package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"fix-ticket-system/models"
	"fix-ticket-system/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTicketService is a mock for the ticket service
type MockTicketService struct {
	mock.Mock
}

var _ service.TicketServiceInterface = (*MockTicketService)(nil)

func (m *MockTicketService) CreateTicket(title, description, createdBy string) (*models.Ticket, error) {
	args := m.Called(title, description, createdBy)
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *MockTicketService) GetTicket(id uuid.UUID) (*models.Ticket, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *MockTicketService) GetAllTickets() ([]models.Ticket, error) {
	args := m.Called()
	return args.Get(0).([]models.Ticket), args.Error(1)
}

func (m *MockTicketService) UpdateTicket(id uuid.UUID, title, description string, status models.Status, priority models.Priority, assignedTo string) (*models.Ticket, error) {
	args := m.Called(id, title, description, status, priority, assignedTo)
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *MockTicketService) DeleteTicket(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	InitializeRoutes(r)
	return r
}

func TestHealthEndpoint(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestCreateTicket(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	expectedTicket := &models.Ticket{
		ID:          uuid.New(),
		Title:       "Test Title",
		Description: "Test Description",
		CreatedBy:   "test@example.com",
	}

	mockService.On("CreateTicket", "Test Title", "Test Description", "test@example.com").Return(expectedTicket, nil)

	r := setupRouter()
	reqBody, _ := json.Marshal(map[string]string{
		"title":       "Test Title",
		"description": "Test Description",
		"created_by":  "test@example.com",
	})
	req := httptest.NewRequest("POST", "/api/v1/tickets/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Ticket
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedTicket.ID, response.ID)
	assert.Equal(t, expectedTicket.Title, response.Title)
	assert.Equal(t, expectedTicket.Description, response.Description)
	assert.Equal(t, expectedTicket.CreatedBy, response.CreatedBy)
}

func TestGetTickets(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	expectedTickets := []models.Ticket{
		{ID: uuid.New(), Title: "T1", Description: "D1", CreatedBy: "a@example.com"},
		{ID: uuid.New(), Title: "T2", Description: "D2", CreatedBy: "b@example.com"},
	}
	mockService.On("GetAllTickets").Return(expectedTickets, nil)

	r := setupRouter()
	req := httptest.NewRequest("GET", "/api/v1/tickets/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.Ticket
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 2)
}

func TestGetTickets_Error(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService
	mockService.On("GetAllTickets").Return([]models.Ticket{}, fmt.Errorf("db error"))

	r := setupRouter()
	req := httptest.NewRequest("GET", "/api/v1/tickets/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "db error")
}

func TestGetTicket(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	id := uuid.New()
	expectedTicket := &models.Ticket{ID: id, Title: "T", Description: "D", CreatedBy: "a@example.com"}
	mockService.On("GetTicket", id).Return(expectedTicket, nil)

	r := setupRouter()
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tickets/%s", id.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Ticket
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedTicket.ID, response.ID)
}

func TestGetTicket_NotFound(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	id := uuid.New()
	mockService.On("GetTicket", id).Return((*models.Ticket)(nil), fmt.Errorf("not found"))

	r := setupRouter()
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tickets/%s", id.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Ticket not found")
}

func TestGetTicket_InvalidID(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	r := setupRouter()
	req := httptest.NewRequest("GET", "/api/v1/tickets/invalid-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid ticket ID")
}

func TestUpdateTicket(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	id := uuid.New()
	expectedTicket := &models.Ticket{
		ID: id, Title: "Updated", Description: "Updated", CreatedBy: "a@example.com",
		Status: models.StatusInProgress, Priority: models.PriorityHigh, AssignedTo: "assignee@example.com",
	}
	mockService.On("UpdateTicket", id, "Updated", "Updated", models.StatusInProgress, models.PriorityHigh, "assignee@example.com").Return(expectedTicket, nil)

	r := setupRouter()
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":       "Updated",
		"description": "Updated",
		"status":      models.StatusInProgress,
		"priority":    models.PriorityHigh,
		"assigned_to": "assignee@example.com",
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/tickets/%s", id.String()), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Ticket
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedTicket.ID, response.ID)
	assert.Equal(t, expectedTicket.Title, response.Title)
}

func TestUpdateTicket_NotFound(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	id := uuid.New()
	mockService.On("UpdateTicket", id, "Updated", "Updated", models.StatusInProgress, models.PriorityHigh, "assignee@example.com").Return((*models.Ticket)(nil), fmt.Errorf("not found"))

	r := setupRouter()
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":       "Updated",
		"description": "Updated",
		"status":      models.StatusInProgress,
		"priority":    models.PriorityHigh,
		"assigned_to": "assignee@example.com",
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/tickets/%s", id.String()), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestUpdateTicket_InvalidID(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	r := setupRouter()
	reqBody, _ := json.Marshal(map[string]interface{}{
		"title":       "Updated",
		"description": "Updated",
		"status":      models.StatusInProgress,
		"priority":    models.PriorityHigh,
		"assigned_to": "assignee@example.com",
	})
	req := httptest.NewRequest("PUT", "/api/v1/tickets/invalid-uuid", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid ticket ID")
}

func TestDeleteTicket(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	id := uuid.New()
	mockService.On("DeleteTicket", id).Return(nil)

	r := setupRouter()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/tickets/%s", id.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Ticket deleted successfully")
}

func TestDeleteTicket_NotFound(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	id := uuid.New()
	mockService.On("DeleteTicket", id).Return(fmt.Errorf("not found"))

	r := setupRouter()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/tickets/%s", id.String()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestDeleteTicket_InvalidID(t *testing.T) {
	mockService := new(MockTicketService)
	ticketService = mockService

	r := setupRouter()
	req := httptest.NewRequest("DELETE", "/api/v1/tickets/invalid-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid ticket ID")
}
