package models

type SalePayload struct {
	OperationType  string     `json:"operation_type"`
	OperationTime  CustomTime `json:"operation_time"`
	OperationPlace string     `json:"operation_place"`
	Passenger      Passenger  `json:"passenger"`
	Routes         []Route    `json:"routes"`
}

type Passenger struct {
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Patronymic    string    `json:"patronymic"`
	DocType       string    `json:"doc_type"`
	DocNumber     string    `json:"doc_number"`
	Birthdate     Birthdate `json:"birthdate"`
	Gender        string    `json:"gender"`
	PassengerType string    `json:"passenger_type"`
	TicketNumber  string    `json:"ticket_number"`
	TicketType    int       `json:"ticket_type"`
}

type Route struct {
	AirlineCode    string     `json:"airline_code"`
	FlightNum      int        `json:"flight_num"`
	DepartPlace    string     `json:"depart_place"`
	DepartDatetime CustomTime `json:"depart_datetime"`
	ArrivePlace    string     `json:"arrive_place"`
	ArriveDatetime CustomTime `json:"arrive_datetime"`
	PNRID          string     `json:"pnr_id"`
}
