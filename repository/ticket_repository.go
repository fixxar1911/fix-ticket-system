package repository

import (
	"fix-ticket-system/config"
	"fix-ticket-system/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		db: config.DB,
	}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) GetByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) GetAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *TicketRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Ticket{}, "id = ?", id).Error
}
