package app

import (
	"air-monolith/internal/models"
	"air-monolith/internal/schemas"
	utils "air-monolith/pkg/utils"
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
)

func (app *Application) Sale(w http.ResponseWriter, r *http.Request) {
	var ticket models.SalePayload
	err := utils.ReadJSON(w, r, schemas.SaleLoader, &ticket)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrBodyTooLarge):
			utils.ErrorJSON(w, http.StatusRequestEntityTooLarge)
			return
		default:
			utils.ErrorJSON(w)
			return
		}
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

	err = app.DB.CreateSale(r.Context(), segments)
	if err != nil {
		switch pgErr := err.(type) {
		case *pgconn.PgError:
			switch pgErr.Code {
			case models.DublicateCode:
				utils.ErrorJSON(w, http.StatusConflict)
				return
			}
		default:
			utils.ErrorJSON(w)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK)

}

func (app *Application) Refund(w http.ResponseWriter, r *http.Request) {
	var rp models.RefundPayload

	err := utils.ReadJSON(w, r, schemas.RefundLoader, &rp)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrBodyTooLarge):
			utils.ErrorJSON(w, http.StatusRequestEntityTooLarge)
			return
		default:
			utils.ErrorJSON(w)
			return
		}
	}

	err = app.DB.RefundTicketsByTicketNumber(r.Context(), rp.TicketNumber)
	if err != nil {
		switch err {
		case models.ErrTicketRefund:
			utils.ErrorJSON(w, http.StatusConflict)
		default:
			utils.ErrorJSON(w, http.StatusInternalServerError)
		}
	}

	utils.WriteJSON(w, http.StatusOK)
}
