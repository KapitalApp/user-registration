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
	"fmt"
	"github.com/google/uuid"
	"net/http"
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
		Email:             input.Email,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		ProvinceCode:      input.ProvinceCode,
		CountryCodeAlpha2: input.CountryCodeAlpha2,
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(app.readParam(r, "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := data.User{
		ID:                     id,
		Email:                  "email",
		FirstName:              "ff",
		LastName:               "dd",
		ProvinceCode:           "ee",
		CountryCodeAlpha2:      "rr",
		AdministrativeDivision: "rr",
		AgeRange:               data.RangeNumber{UpLimit: 30, DownLimit: 25},
		FamilyMemberNumber:     1,
		CreatedAt:              time.Now(),
		Meta: []data.MetaField{{
			Key:       "Incomplete",
			Namespace: "Registration",
			Value:     "True",
			Type:      "Bool",
		}},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
