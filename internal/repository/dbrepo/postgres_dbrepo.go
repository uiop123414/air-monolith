package dbrepo

import (
	"air-monolith/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)


type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeOut = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) CreateSale(segments []models.Segment) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO segments (
		operation_type, operation_time,operation_time_zone, operation_place, 
		passenger_name, passenger_surname, passenger_patronymic, 
		doc_type, doc_number, birthdate, gender, passenger_type, 
		ticket_number, ticket_type, airline_code, flight_num, 
		depart_place, depart_datetime, arrive_place, arrive_datetime, 
		pnr_id, serial_number) 
	VALUES ($1, 
		($2 AT TIME ZONE 'UTC') AT TIME ZONE 'UTC',
		EXTRACT(TIMEZONE FROM $2) / 60,
		$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	for i, segment := range segments {
		_, err := tx.ExecContext(ctx, query,			
			segment.OperationType, segment.OperationTime.Time, segment.OperationPlace,
			segment.PassengerName, segment.PassengerSurname, segment.PassengerPatronymic,
			segment.DocType, segment.DocNumber, segment.Birthdate, segment.Gender,
			segment.PassengerType, segment.TicketNumber, segment.TicketType,
			segment.AirlineCode, segment.FlightNum, segment.DepartPlace,
			segment.DepartDatetime, segment.ArrivePlace, segment.ArriveDatetime,
			segment.PNRID, i+1 ,
		)
		if err != nil {
			_ = tx.Rollback()
			fmt.Println(err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetTicketByTicketNumber(tn string) (*models.Segment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
		SELECT operation_type, operation_time, operation_place, 
			passenger_name, passenger_surname, passenger_patronymic, 
			doc_type, doc_number, birthdate, gender, passenger_type, 
			ticket_number, ticket_type, airline_code, flight_num, 
			depart_place, depart_datetime, arrive_place, arrive_datetime, 
			pnr_id 
		FROM segments 
		WHERE ticket_number = $1
		LIMIT 1
	`

	var segment models.Segment

	err := m.DB.QueryRowContext(ctx, query, tn).Scan(
		&segment.OperationType,
		&segment.OperationTime,
		&segment.OperationPlace,
		&segment.PassengerName,
		&segment.PassengerSurname,
		&segment.PassengerPatronymic,
		&segment.DocType,
		&segment.DocNumber,
		&segment.Birthdate,
		&segment.Gender,
		&segment.PassengerType,
		&segment.TicketNumber,
		&segment.TicketType,
		&segment.AirlineCode,
		&segment.FlightNum,
		&segment.DepartPlace,
		&segment.DepartDatetime,
		&segment.ArrivePlace,
		&segment.ArriveDatetime,
		&segment.PNRID,
	)

	if err != nil {
		switch {
		case errors.Is(err,sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}

	return &segment, nil
}