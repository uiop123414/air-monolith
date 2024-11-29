package main

import (
	"air-monolith/internal/models"
	"air-monolith/internal/rww"
	"encoding/json"
	"io"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const RequestSizeLimit = 2 * 1024 * 1024 // 2 kBs

func (app *application) writeJSON(w http.ResponseWriter, status int, headers ...http.Header) error {
	rw, ok := w.(*rww.ResponseWriterWrapper)
	if !ok {
		rw = &rww.ResponseWriterWrapper{ResponseWriter: w}
	}

	if rw.HasWritten() {
		return models.ErrAlreadyResponded
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.WriteHeader(status)
	w.Write([]byte{})
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, loader gojsonschema.JSONLoader, data interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(RequestSizeLimit))
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = app.validateJSON(loader, bodyBytes)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, data)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	return app.writeJSON(w, statusCode)
}

func (app *application) validateJSON(loader gojsonschema.JSONLoader, data []byte) error {
	payloadLoader := gojsonschema.NewBytesLoader(data)

	result, err := gojsonschema.Validate(loader, payloadLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return models.ErrJSONNotValid
	}

	return nil
}
