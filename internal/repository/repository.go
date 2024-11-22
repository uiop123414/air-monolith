package repository

import (
	"air-monolith/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	CreateSale(sg []models.Segment) error
	GetTicketByTicketNumber(tn string) (*models.Segment, error)
}
