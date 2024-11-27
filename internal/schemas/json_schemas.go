package schemas

import "github.com/xeipuuv/gojsonschema"

var (
	BirthdateLoader  = gojsonschema.NewReferenceLoader("file://./jsons/Birthdate.json")
	CustomTimeLoader = gojsonschema.NewReferenceLoader("file://./jsons/CustomTime.json")
	PassengerLoader  = gojsonschema.NewReferenceLoader("file://./jsons/Birthdate.json")
	RefundLoader     = gojsonschema.NewReferenceLoader("file://./jsons/Redund.json")
	RouteLoader      = gojsonschema.NewReferenceLoader("file://./jsons/Route.json")
	SaleLoader       = gojsonschema.NewReferenceLoader("file://./jsons/Sale.json")
)
