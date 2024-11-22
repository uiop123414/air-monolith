package main

import (
	"air-monolith/internal/models"
	"air-monolith/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) Sale(w http.ResponseWriter, r *http.Request) {
	var ticket SalePayload
	fmt.Println("Hello from sale")
	err := app.readJSON(w, r, &ticket)
	if err != nil {
		switch {
		case err.Error() == "http: request body too large":
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
		app.errorJSONWithMSG(w, errors.New("invalid credentials"), v.Errors)
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

	fmt.Println(ticket.OperationTime)

	err = app.DB.CreateSale(segments)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, ticket)

	fmt.Println(app.DB.GetTicketByTicketNumber(ticket.Passenger.TicketNumber))
}
