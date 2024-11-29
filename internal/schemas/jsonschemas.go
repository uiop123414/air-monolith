package schemas

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

const path = "/app/schemas/jsons"

var (
    BirthdateLoader  = getJsonSchema("Birthdate")
    CustomTimeLoader = getJsonSchema("CustomTime")
    PassengerLoader  = getJsonSchema("Birthdate")
    RefundLoader     = getJsonSchema("Refund")
    RouteLoader      = getJsonSchema("Route")
    SaleLoader       = getJsonSchema("Sale")
)

func getJsonSchema(name string) gojsonschema.JSONLoader {
    return gojsonschema.NewReferenceLoader(fmt.Sprintf(`file://%s/%s.schema.json`, path, name))
}