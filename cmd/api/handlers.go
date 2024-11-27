package main

import (
	"air-monolith/internal/models"
	"air-monolith/internal/validator"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
)

func (app *application) Sale(w http.ResponseWriter, r *http.Request) {
	var ticket SalePayload
	err := app.readJSON(w, r, &ticket)
	if err != nil {
		switch {
		case err == models.ErrBodyTooLarge:
			app.errorJSON(w, err, http.StatusRequestEntityTooLarge)
			return
		default:
			app.errorJSON(w, err)
			return
		}
	}

	v := validator.New()

	ValidateGender(v, ticket.Passenger.Gender)
	ValidateTicketNumber(v, ticket.Passenger.TicketNumber)
	ValidateDocNumber(v, ticket.Passenger.DocType, ticket.Passenger.DocNumber)

	if !v.Valid() {
		app.errorJSONWithMSG(w, models.ErrInvalidCredentils, v.Errors)
		return
	}

	segments := []models.Segment{}

	for _, v := range ticket.Routes {
		sg := models.Segment{
			OperationType:       ticket.OperationType,
			OperationTime:       ticket.OperationTime,
			OperationPlace:      ticket.OperationPlace,
			PassengerName:       ticket.Passenger.Name,
			PassengerSurname:    ticket.Passenger.Surname,
			PassengerPatronymic: ticket.Passenger.Patronymic,
			DocType:             ticket.Passenger.DocType,
			DocNumber:           ticket.Passenger.DocNumber,
			Birthdate:           ticket.Passenger.Birthdate,
			Gender:              ticket.Passenger.Gender,
			PassengerType:       ticket.Passenger.PassengerType,
			TicketNumber:        ticket.Passenger.TicketNumber,
			TicketType:          ticket.Passenger.TicketType,
			AirlineCode:         v.AirlineCode,
			FlightNum:           v.FlightNum,
			DepartPlace:         v.DepartPlace,
			DepartDatetime:      v.DepartDatetime,
			ArrivePlace:         v.ArrivePlace,
			ArriveDatetime:      v.ArriveDatetime,
			PNRID:               v.PNRID,
		}

		segments = append(segments, sg)
	}

	err = app.DB.CreateSale(segments)
	if err != nil {
		switch pgErr := err.(type) {
		case *pgconn.PgError:
			switch pgErr.Code {
			case models.DublicateCode:
				app.errorJSON(w, models.ErrTicketAlreadyExists, http.StatusConflict)
				return
			}
		default:
			app.errorJSON(w, err)
		}
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "ticket was sold"

	app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) Refund(w http.ResponseWriter, r *http.Request) {
	var rp RefundPayload

	err := app.readJSON(w, r, &rp)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	segments, err := app.DB.GetTicketsByTicketNumber(rp.TicketNumber)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if segments == nil {
		app.errorJSON(w, models.ErrNoSale(rp.TicketNumber), http.StatusConflict)
		return
	}

	fmt.Println(segments[0].OperationTime)

	for _, sg := range segments {
		if sg.OperationType == models.OperationRefund {
			app.errorJSON(w, models.ErrTicketWasRefunded(rp.TicketNumber), http.StatusConflict)
			return
		}
	}

	err = app.DB.RefundTicketsByTicketNumber(rp.TicketNumber, len(segments))
	if err != nil {
		app.errorJSON(w, models.ErrServerError, http.StatusInternalServerError)
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "ticket was refunded"

	app.writeJSON(w, http.StatusOK, payload)
}
