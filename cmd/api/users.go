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
	"user-registration.kptl.net/internal/data"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(app.readParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	user := data.User{
		ID:                     id,
		Email:                  "email",
		FirstName:              "ff",
		LastName:               "dd",
		ProvinceCode:           "ee",
		CountryCode:            "rr",
		AdministrativeDivision: "rr",
		AgeRange:               data.RangeNumber{UpLimit: 30, DownLimit: 25},
		FamilyNumber:           1,
		CreatedAt:              time.Now(),
		Meta: []data.MetaField{{
			Key:       "Incomplete",
			Namespace: "Registration",
			Value:     "True",
			Type:      "Bool",
		}},
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
