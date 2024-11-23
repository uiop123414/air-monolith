package models

import "errors"

var (
	ErrDublicate = errors.New("ERROR: duplicate key value violates unique constraint \"unique_ticket_serial\" (SQLSTATE 23505)")
)