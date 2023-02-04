// Package data /*
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
package data

import (
	"github.com/google/uuid"
	"time"
	"user-service.kptl.net/internal/validator"
)

type User struct {
	ID                     uuid.UUID      `json:"id"`
	Email                  string         `json:"email"`
	FirstName              string         `json:"first_name"`
	LastName               string         `json:"last_name"`
	ProvinceCode           string         `json:"province_code"`
	CountryCodeAlpha2      string         `json:"country_code_alpha_2"`
	AdministrativeDivision string         `json:"administrative_division"`
	AgeRange               RangeNumber    `json:"age_range,omitempty"`
	IncomeRange            RangeNumber    `json:"income_range,omitempty"`
	ExpensesRange          RangeNumber    `json:"expenses_range,omitempty"`
	FamilyMemberNumber     int64          `json:"family_member_number,omitempty"`
	IsMarried              bool           `json:"is_married,omitempty"`
	Spouse                 FamilyMember   `json:"spouse,omitempty"`
	Dependent              []FamilyMember `json:"dependent,omitempty"`
	Milestones             []Milestone    `json:"milestones,omitempty"`
	Goals                  []Goal         `json:"goals,omitempty"`
	CreatedAt              time.Time      `json:"created_at,omitempty"`
	Meta                   []MetaField    `json:"meta,omitempty"`
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be valid")
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(user.LastName != "", "last_name", "must be provided")
	v.Check(len(user.CountryCodeAlpha2) == 2, "country_code_alpha_2", "must be two letters")
	v.Check(user.ProvinceCode != "", "province_code", "must be provided")
}
