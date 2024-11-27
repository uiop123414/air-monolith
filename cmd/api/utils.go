package main

import (
	"air-monolith/internal/models"
	"air-monolith/internal/rww"
	"encoding/json"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const maxBytes = 2 * 1024 * 1024 // 2 kBs

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	rw, ok := w.(*rww.ResponseWriterWrapper)
	if !ok {
		rw = &rww.ResponseWriterWrapper{ResponseWriter: w}
	}

	if rw.HasWritten() {
		return models.ErrAlreadyResponded
	}

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return models.ErrBodyMustConainSingleValue
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

func (app *application) errorJSONWithMSG(w http.ResponseWriter, err error, errors map[string]string, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()
	payload.Data = errors

	return app.writeJSON(w, statusCode, payload)
}
