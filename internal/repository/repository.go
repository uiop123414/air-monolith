package repository

import (
	"air-monolith/internal/models"
	"context"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	CreateSale(ctx context.Context, sg []models.Segment) error
	RefundTicketsByTicketNumber(ctx context.Context, tn string) error
}
