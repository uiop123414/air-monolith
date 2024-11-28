package schemas

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

const path = "/app/schemas/jsons"

var (
    BirthdateLoader  = gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/Birthdate.json`, path))
    CustomTimeLoader = gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/CustomTime.json`, path))
    PassengerLoader  = gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/Birthdate.json`, path))
    RefundLoader     = gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/Refund.json`, path))
    RouteLoader      = gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/Route.json`, path))
    SaleLoader       = gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/Sale.json`, path))
)
