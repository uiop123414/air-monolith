package dbrepo

import (
	"air-monolith/internal/models"
	"database/sql"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeOut = time.Second * 120

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) CreateSale(segments []models.Segment) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	m.setDBTimeout(tx)

	const query = `
	INSERT INTO segments (
		operation_type, operation_time, operation_time_timezone, operation_place, 
		passenger_name, passenger_surname, passenger_patronymic, 
		doc_type, doc_number, birthdate, gender, passenger_type, 
		ticket_number, ticket_type, airline_code, flight_num, 
		depart_place, depart_datetime, depart_datetime_timezone, arrive_place, arrive_datetime, arrive_datetime_timezone,
		pnr_id, serial_number)
	VALUES (
		$1, $2 AT TIME ZONE 'UTC', $3, $4, $5, $6, $7, $8, $9, $10, $11, 
		$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
	);
	`

	for i, segment := range segments {
		_, err := tx.Exec(query,
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
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) getSaleTicketsCountByTicketNumber(tx *sql.Tx, tn string) (int64, error) {
	const query = `
		SELECT serial_number
		FROM segments 
		WHERE ticket_number = $1 and operation_type='sale' FOR NO KEY UPDATE
	`

	var sns []int64

	rows, err := tx.Query(query, tn)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for rows.Next() {
		var sn int64

		rows.Scan(&sn)

		if err := rows.Err(); err != nil {
			_ = tx.Rollback()
			return 0, err
		}

		sns = append(sns, sn)
	}

	return int64(len(sns)), nil
}

func (m *PostgresDBRepo) RefundTicketsByTicketNumber(tn string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	m.setDBTimeout(tx)

	count, err := m.getSaleTicketsCountByTicketNumber(tx, tn)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if count == 0 {
		_ = tx.Rollback()
		return models.ErrTicketRefund
	}

	const query = `
		UPDATE segments
		SET operation_type='refund'
		WHERE ticket_number = $1 and operation_type='sale'`

	result, err := tx.Exec(query, tn)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	num, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if num == 0 {
		return models.ErrTicketRefund
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) setDBTimeout(tx *sql.Tx) error {
	const query = `SET LOCAL lock_timeout = '120s';`

	_, err := tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
