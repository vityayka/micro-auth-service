package main

import (
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &request)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("bad request json"), http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(request.Email)
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("invalid creds"), http.StatusUnauthorized)
		return
	}

	isValid, err := user.PasswordMatches(request.Password)
	if err != nil || !isValid {
		app.ErrorJson(w, fmt.Errorf("invalid creds"), http.StatusUnauthorized)
		return
	}

	err = app.logRequest("authenticated", fmt.Sprintf("user: %s", user.Email))
	if err != nil {
		app.ErrorJson(w, fmt.Errorf("error writing log"), http.StatusInternalServerError)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("user authenticated: %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, response)
}
