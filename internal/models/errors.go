package models

import (
	"errors"
	"fmt"
)

var (
	ErrDublicate                 = errors.New("ERROR: duplicate key value violates unique constraint \"unique_ticket_serial\" (SQLSTATE 23505)")
	ErrBodyMustConainSingleValue = errors.New("ERROR: body must only contain a single JSON value")
	ErrInvalidCredentils         = errors.New("ERROR: invalid credentials")
	ErrTicketAlreadyExists       = errors.New("ERROR: ticket already exists")
	ErrServerError               = errors.New("ERROR: server error")
	ErrColumnNotSupported        = errors.New("ERROR: column type not supported")
	ErrRequestTimeout            = errors.New("ERROR: request Timeout")
	ErrAlreadyResponded          = errors.New("already responded")

	ErrBodyTooLarge = errors.New("http: request body too large")
)

var (
	DublicateCode = "23505"
)

func ErrNoSale(tn string) error {
	return errors.New(fmt.Sprintf("ERROR: no sale by %s ticket id", tn))
}

func ErrTicketWasRefunded(tn string) error {
	return errors.New(fmt.Sprintf("ERROR: ticket %s was refunded", tn))
}
