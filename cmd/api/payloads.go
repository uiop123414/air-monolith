package main

import (
	"air-monolith/internal/models"
	"air-monolith/internal/validator"
)

type SalePayload struct {
	OperationType  string            `json:"operation_type"`
	OperationTime  models.CustomTime `json:"operation_time"`
	OperationPlace string            `json:"operation_place"`
	Passenger      Passenger         `json:"passenger"`
	Routes         []Route           `json:"routes"`
}

type Passenger struct {
	Name          string           `json:"name"`
	Surname       string           `json:"surname"`
	Patronymic    string           `json:"patronymic"`
	DocType       string           `json:"doc_type"`
	DocNumber     string           `json:"doc_number"`
	Birthdate     models.Birthdate `json:"birthdate"`
	Gender        string           `json:"gender"`
	PassengerType string           `json:"passenger_type"`
	TicketNumber  string           `json:"ticket_number"`
	TicketType    int              `json:"ticket_type"`
}

type Route struct {
	AirlineCode    string            `json:"airline_code"`
	FlightNum      int               `json:"flight_num"`
	DepartPlace    string            `json:"depart_place"`
	DepartDatetime models.CustomTime `json:"depart_datetime"`
	ArrivePlace    string            `json:"arrive_place"`
	ArriveDatetime models.CustomTime `json:"arrive_datetime"`
	PNRID          string            `json:"pnr_id"`
}

type RefundPayload struct {
	OperationType  string            `json:"operation_type"`
	OperationTime  models.CustomTime `json:"operation_time"`
	OperationPlace string            `json:"operation_place"`
	TicketNumber   string            `json:"ticket_number"`
}

func ValidateTicketNumber(v *validator.Validator, ticketNumber string) {
	v.Check(ticketNumber != "", "ticket_number", "must be provided")
	v.Check(validator.Matches(ticketNumber, validator.NNumbersRegex(13)), "ticket_number", "must be exactly 13 in length and contain only numbers")
}

func ValidateGender(v *validator.Validator, gender string) {
	v.Check(gender != "", "gender", "must be provided")
	v.Check(validator.In(gender, "F", "M"), "gender", "must be F or M")
}

func ValidateDocNumber(v *validator.Validator, DocType, DocNumber string) {
	v.Check(DocType != "", "doc_type", "must be provided")
	v.Check(DocNumber != "", "doc_number", "must be provided")
	if DocType == "00" {
		v.Check(validator.Matches(DocNumber, validator.NNumbersRegex(10)), "doc_number", "must be exactly 10 in length and contain only numbers")
	}
}
