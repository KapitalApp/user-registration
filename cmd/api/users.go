/*
Copyright 2023 The Kapital Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"time"
	"user-service.kptl.net/internal/data"
	"user-service.kptl.net/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email             string `json:"email"`
		FirstName         string `json:"first_name"`
		LastName          string `json:"last_name"`
		ProvinceCode      string `json:"province_code"`
		CountryCodeAlpha2 string `json:"country_code"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		ID:                     uuid.New().String(),
		Email:                  input.Email,
		FirstName:              input.FirstName,
		LastName:               input.LastName,
		ProvinceCode:           input.ProvinceCode,
		CountryCodeAlpha2:      input.CountryCodeAlpha2,
		AdministrativeDivision: "province",
		CreatedAt:              time.Now().Format("2006-01-02"),
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(app.readParam(r, "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.Get(id.String())
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(app.readParam(r, "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.Get(id.String())
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	input := data.User{}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	newAttributes := make(map[string]interface{})

	inType := reflect.TypeOf(input)
	inValue := reflect.ValueOf(input)
	for i := 0; i < inValue.NumField(); i++ {
		fieldValue := inValue.Field(i).String()
		if fieldValue != "" && !inValue.Field(i).IsNil() {
			fmt.Printf("%s: %s\n", inType.Field(i).Name, fieldValue)
			newAttributes[inType.Field(i).Name] = fieldValue
		}
	}

	attributes, err := app.models.Users.Update(user, newAttributes)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": attributes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
