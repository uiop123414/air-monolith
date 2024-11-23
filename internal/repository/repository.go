package repository

import (
	"air-monolith/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	CreateSale(sg []models.Segment) error
	GetTicketsByTicketNumber(tn string) ([]models.Segment, error)
	RefundTicketsByTicketNumber(tn string, count int) error
}
