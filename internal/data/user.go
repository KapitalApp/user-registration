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
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"user-service.kptl.net/internal/validator"
)

type User struct {
	ID                     string         `json:"id" dynamodbav:"ID"`
	Email                  string         `json:"email" dynamodbav:"email"`
	FirstName              string         `json:"first_name" dynamodbav:"firstName"`
	LastName               string         `json:"last_name" dynamodbav:"lastName"`
	ProvinceCode           string         `json:"province_code" dynamodbav:"provinceCode"`
	CountryCodeAlpha2      string         `json:"country_code_alpha_2" dynamodbav:"countryCodeAlpha2"`
	Currency               string         `json:"currency" dynamodbav:"currency"`
	AdministrativeDivision string         `json:"administrative_division" dynamodbav:"administrativeDivision"`
	DateOfBirth            string         `json:"age,omitempty" dynamodbav:"dateOfBirth,omitempty"`
	Income                 string         `json:"income,omitempty" dynamodbav:"income,omitempty"`
	Expenses               string         `json:"expenses,omitempty" dynamodbav:"expenses,omitempty"`
	FamilyMemberNumber     int64          `json:"family_member_number,omitempty" dynamodbav:"familyMemberNumber,omitempty"`
	IsMarried              bool           `json:"is_married,omitempty" dynamodbav:"isMarried,omitempty"`
	Spouse                 *FamilyMember  `json:"spouse,omitempty" dynamodbav:"spouse,omitempty"`
	Dependents             []FamilyMember `json:"dependent,omitempty" dynamodbav:"dependents,omitempty"`
	Milestones             []Milestone    `json:"milestones,omitempty" dynamodbav:"milestones,omitempty"`
	Goals                  []Goal         `json:"goals,omitempty" dynamodbav:"goals,omitempty"`
	Protections            []Protection   `json:"protections,omitempty" dynamodbav:"protections,omitempty"`
	Debts                  []Debt         `json:"debts,omitempty" dynamodbav:"debts,omitempty"`
	CreatedAt              string         `json:"created_at,omitempty" dynamodbav:"createdAt"`
	Version                int64          `json:"-" dynamodbav:"version"`
	Meta                   []MetaField    `json:"meta,omitempty" dynamodbav:"meta,omitempty"`
}

func (user User) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(user.ID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"ID": id}
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be valid")
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(user.LastName != "", "last_name", "must be provided")
	v.Check(len(user.CountryCodeAlpha2) == 2, "country_code_alpha_2", "must be two letters")
	v.Check(user.ProvinceCode != "", "province_code", "must be provided")

	if user.IsMarried {
		v.Check(user.Spouse != nil, "spouse", "must be provided if married")
		v.Check(ValidateFamilyMember(v, user.Spouse), "spouse", "must be valid")
	} else if user.Dependents != nil {
		for i, dep := range user.Dependents {
			v.Check(ValidateFamilyMember(v, &dep), fmt.Sprintf("dependents_%d", i), "must be valid")
		}
	}
}

func ValidateFamilyMember(v *validator.Validator, familyMember *FamilyMember) bool {
	current := len(v.Errors)
	v.Check(familyMember.FirstName != "", "first_name", "must be provided")
	v.Check(familyMember.LastName != "", "last_name", "must be provided")
	if current != len(v.Errors) {
		return false
	}
	return true
}
