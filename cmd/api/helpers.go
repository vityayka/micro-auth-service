package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain the only json value")
	}

	return nil
}

func (app *Config) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse

	payload.Error = true
	payload.Message = err.Error()

	fmt.Println(err)
	return app.WriteJSON(w, statusCode, payload)
}

func (app *Config) logRequest(name, data string) error {
	var requestPayload struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	requestPayload.Name = name
	requestPayload.Data = data

	requestBytes, _ := json.Marshal(requestPayload)

	request, err := http.NewRequest("POST", "http://logger-service", bytes.NewBuffer(requestBytes))
	if err != nil {
		log.Printf("error logging request: %v", err)
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		log.Printf("error logging request: %v", err)
	}
	return err
}
