package service

import (
	"fix-ticket-system/metrics"
	"fix-ticket-system/models"
	"fix-ticket-system/repository"
	"time"

	"github.com/google/uuid"
)

type TicketService struct {
	repo *repository.TicketRepository
}

func NewTicketService() *TicketService {
	return &TicketService{
		repo: repository.NewTicketRepository(),
	}
}

type TicketServiceInterface interface {
	CreateTicket(title, description, createdBy string) (*models.Ticket, error)
	GetTicket(id uuid.UUID) (*models.Ticket, error)
	GetAllTickets() ([]models.Ticket, error)
	UpdateTicket(id uuid.UUID, title, description string, status models.Status, priority models.Priority, assignedTo string) (*models.Ticket, error)
	DeleteTicket(id uuid.UUID) error
}

var _ TicketServiceInterface = (*TicketService)(nil)

func (s *TicketService) CreateTicket(title, description, createdBy string) (*models.Ticket, error) {
	ticket := models.NewTicket(title, description, createdBy)
	if err := s.repo.Create(ticket); err != nil {
		metrics.ErrorTotal.WithLabelValues("create_ticket").Inc()
		return nil, err
	}
	metrics.TicketOperationsTotal.WithLabelValues("create", "success").Inc()
	metrics.TicketStatusGauge.WithLabelValues(string(ticket.Status)).Inc()
	return ticket, nil
}

func (s *TicketService) GetTicket(id uuid.UUID) (*models.Ticket, error) {
	ticket, err := s.repo.GetByID(id)
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("get_ticket").Inc()
		return nil, err
	}
	metrics.TicketOperationsTotal.WithLabelValues("get", "success").Inc()
	return ticket, nil
}

func (s *TicketService) GetAllTickets() ([]models.Ticket, error) {
	tickets, err := s.repo.GetAll()
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("get_all_tickets").Inc()
		return nil, err
	}
	metrics.TicketOperationsTotal.WithLabelValues("get_all", "success").Inc()
	return tickets, nil
}

func (s *TicketService) UpdateTicket(id uuid.UUID, title, description string, status models.Status, priority models.Priority, assignedTo string) (*models.Ticket, error) {
	ticket, err := s.repo.GetByID(id)
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("update_ticket").Inc()
		return nil, err
	}

	// Decrement old status count
	metrics.TicketStatusGauge.WithLabelValues(string(ticket.Status)).Dec()

	ticket.Title = title
	ticket.Description = description
	ticket.Status = status
	ticket.Priority = priority
	ticket.AssignedTo = assignedTo
	ticket.UpdatedAt = time.Now()

	if err := s.repo.Update(ticket); err != nil {
		metrics.ErrorTotal.WithLabelValues("update_ticket").Inc()
		return nil, err
	}

	// Increment new status count
	metrics.TicketStatusGauge.WithLabelValues(string(ticket.Status)).Inc()
	metrics.TicketOperationsTotal.WithLabelValues("update", "success").Inc()
	return ticket, nil
}

func (s *TicketService) DeleteTicket(id uuid.UUID) error {
	ticket, err := s.repo.GetByID(id)
	if err != nil {
		metrics.ErrorTotal.WithLabelValues("delete_ticket").Inc()
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		metrics.ErrorTotal.WithLabelValues("delete_ticket").Inc()
		return err
	}

	// Decrement status count
	metrics.TicketStatusGauge.WithLabelValues(string(ticket.Status)).Dec()
	metrics.TicketOperationsTotal.WithLabelValues("delete", "success").Inc()
	return nil
}
