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

	const query = `
	INSERT INTO segments (
		operation_type, operation_time, operation_time_timezone, operation_place, 
		passenger_name, passenger_surname, passenger_patronymic, 
		doc_type, doc_number, birthdate, gender, passenger_type, 
		ticket_number, ticket_type, airline_code, flight_num, 
		depart_place, depart_datetime, depart_datetime_timezone, arrive_place, arrive_datetime, arrive_datetime_timezone,
		pnr_id, serial_number)
	VALUES (
		$1, ($2 AT TIME ZONE 'UTC')::TIMESTAMP, $3, $4, $5, $6, $7, $8, $9, $10, $11, 
		$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
	);
	`

	for i, segment := range segments {
		_, err := tx.ExecContext(ctx, query,
			segment.OperationType, segment.OperationTime,
			segment.OperationTime.GetTimezone(), segment.OperationPlace,
			segment.PassengerName, segment.PassengerSurname, segment.PassengerPatronymic,
			segment.DocType, segment.DocNumber, segment.Birthdate, segment.Gender,
			segment.PassengerType, segment.TicketNumber, segment.TicketType,
			segment.AirlineCode, segment.FlightNum, segment.DepartPlace,
			segment.DepartDatetime, segment.DepartDatetime.GetTimezone(), segment.ArrivePlace, segment.ArriveDatetime, segment.ArriveDatetime.GetTimezone(),
			segment.PNRID, i+1,
		)
		if err != nil {
			fmt.Println(err)
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetTicketsByTicketNumber(tn string) ([]models.Segment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	const query = `
		SELECT operation_type, 
			operation_time AT TIME ZONE (-operation_time_timezone)::TEXT,
			operation_place, passenger_name, passenger_surname, passenger_patronymic, 
			doc_type, doc_number, birthdate, gender, passenger_type, 
			ticket_number, ticket_type, airline_code, flight_num, 
			depart_place, depart_datetime, arrive_place, arrive_datetime,
			pnr_id 
		FROM segments 
		WHERE ticket_number = $1
	`

	var segments []models.Segment

	rows, err := m.DB.QueryContext(ctx, query, tn)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}

	for rows.Next() {
		var segment models.Segment
		rows.Scan(
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

		segments = append(segments, segment)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return segments, nil
}

func (m *PostgresDBRepo) GetSaleTicketsCountByTicketNumber(tn string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	const query = `
		SELECT COUNT(serial_number)
		FROM segments 
		WHERE ticket_number = $1
	`

	var sn int

	err := m.DB.QueryRowContext(ctx, query, tn).Scan(&sn)
	if err != nil {
		return 0, err
	}

	return sn, nil
}

func (m *PostgresDBRepo) RefundTicketsByTicketNumber(tn string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	const query = `
		UPDATE segments
		SET operation_type='refund'
		WHERE ticket_number = $1`

	result, err := m.DB.ExecContext(ctx, query, tn)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	num, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	fmt.Println(num)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
