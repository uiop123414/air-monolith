package models

type RefundPayload struct {
	OperationType  string     `json:"operation_type"`
	OperationTime  CustomTime `json:"operation_time"`
	OperationPlace string     `json:"operation_place"`
	TicketNumber   string     `json:"ticket_number"`
}
