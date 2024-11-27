package repository

import (
	"air-monolith/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	CreateSale(sg []models.Segment) error
	RefundTicketsByTicketNumber(tn string) error
}
