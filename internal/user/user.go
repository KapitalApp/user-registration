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

// Package user defines the structure and the functionality of User.
// It validates User instances and modify the "User" table in dynamodb.
package user

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"user-service.mykapital.io/internal/validator"
)

// User struct is the main struct declaring user fields.
type User struct {
	// ID is the UUID of the user.
	ID           string `dynamodbav:"userID"` // dynamodbav is the representation of the field as a dynamodb attribute.
	Email        string `dynamodbav:"email"`
	FirstName    string `dynamodbav:"firstName"`
	LastName     string `dynamodbav:"lastName"`
	ProvinceCode string `dynamodbav:"provinceCode"`
	// CountryCodeAlpha2 represents the two-letter word representing a country.
	//
	// For example Country Code Alpha 2 for Canada is "CA".
	CountryCodeAlpha2 string `dynamodbav:"countryCodeAlpha2"`
	Currency          string `dynamodbav:"currency"`
	// AdministrativeDivision is the type of the division within a country.
	//
	// For example Administrative Division of Canada is "province".
	AdministrativeDivision string `dynamodbav:"administrativeDivision"`
	DateOfBirth            string `dynamodbav:"dateOfBirth,omitempty"`
	Occupation             string `dynamodbav:"occupation,omitempty"`
	// Income represent the amount in the user's currency.
	Income string `dynamodbav:"income,omitempty"`
	// Expenses represent the amount in the user's currency.
	Expenses           string `dynamodbav:"expenses,omitempty"`
	FamilyMemberNumber int64  `dynamodbav:"familyMemberNumber,omitempty"`
	IsMarried          bool   `dynamodbav:"isMarried,omitempty"`
	// Spouse should be a pointer, else dynamodb would reject the field.
	Spouse      *FamilyMember  `dynamodbav:"spouse,omitempty"`
	Dependents  []FamilyMember `dynamodbav:"dependents,omitempty"`
	Milestones  []Milestone    `dynamodbav:"milestones,omitempty"`
	Goals       []Goal         `dynamodbav:"goals,omitempty"`
	Protections []Protection   `dynamodbav:"protections,omitempty"`
	Debts       []Debt         `dynamodbav:"debts,omitempty"`
	// RiskTolerance can be represented in a different metric.
	//
	// TODO: Find the correct metric for RiskTolerance.
	RiskTolerance string `dynamodbav:"riskTolerance,omitempty"`
	CreatedAt     string `dynamodbav:"createdAt"`
	// Version is used to handle data races
	Version int64       `dynamodbav:"version"`
	Meta    []MetaField `dynamodbav:"meta,omitempty"`
}

// FamilyMember struct declares family member fields
type FamilyMember struct {
	// Type is either Spouse or Child
	Type        string
	FirstName   string
	LastName    string
	DateOfBirth string
	// Income represent the amount in the user's currency.
	Income string
	// Expenses represent the amount in the user's currency.
	Expenses string
}

// Goal struct declares the financial goal of the user
type Goal struct {
	Date              string
	Title             string
	ProgressLevel     string
	EstimatedDuration time.Duration
	Description       string
}

// Milestone struct declares the financial achievement of the user
type Milestone struct {
	Date        string
	Title       string
	Type        string
	Description string
}

// Protection struct declares the financial protection the user
// currently posses.
type Protection struct {
	Type           string
	Premium        int64
	ClaimedDate    string
	ExpirationDate string
	Description    string
}

// Debt struct declares the financial debt the user
// currently posses.
type Debt struct {
	Type         string
	Cost         string
	InterestRate int64
	Term         int64
	Collateral   string
	Description  string
}

// MetaField struct declares user's personalized configuration
type MetaField struct {
	Key       string
	Namespace string
	Value     string
	Type      string
}

// GetKey is used to create a primary key for dynamodb.
// The ID of the user is used as the primary key.
func (user User) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(user.ID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"userID": id}
}

// ValidateUser validates User data.
//
// The email address of the user should follow the regex validator.EmailRX.
// First name, last name, province code, spouse (if applicable) and
// dependent (if applicable) must be provided.
// Spouse (if applicable) and dependents (if applicable) must be validated.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be valid")
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(len(user.CountryCodeAlpha2) == 2, "country_code_alpha_2", "must be two letters")
	v.Check(user.ProvinceCode != "", "province_code", "must be provided")

	if user.IsMarried {
		v.Check(user.Spouse != nil, "spouse", "must be provided")
		if user.Spouse != nil {
			ValidateFamilyMember(v, user.Spouse, "spouse")
		}
	}

	if user.Dependents != nil {
		for i, dep := range user.Dependents {
			depName := fmt.Sprintf("dependent_%d", i+1)
			ValidateFamilyMember(v, &dep, depName)
		}
	}
}

// ValidateFamilyMember validates FamilyMember data.
//
// First name and last name must be provided.
func ValidateFamilyMember(v *validator.Validator, familyMember *FamilyMember, uniqueName string) {
	v.Check(familyMember.Type != "", uniqueName+"_type", "must be provided")
	v.Check(familyMember.FirstName != "", uniqueName+"_first_name", "must be provided")
}
