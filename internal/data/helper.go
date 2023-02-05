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

import "time"

type FamilyMember struct {
	Type        string `json:"type" dynamodbav:"type"`
	FirstName   string `json:"first_name" dynamodbav:"firstName"`
	LastName    string `json:"last_name" dynamodbav:"lastName"`
	DateOfBirth string `json:"age_range" dynamodbav:"dateOfBirth"`
	Income      int64  `json:"income_range" dynamodbav:"income"`
	Expenses    int64  `json:"expenses_range" dynamodbav:"expenses"`
}

type Goal struct {
	Date              string        `json:"date" dynamodbav:"date"`
	Title             string        `json:"title" dynamodbav:"title"`
	ProgressLevel     string        `json:"progress_level" dynamodbav:"progressLevel"`
	EstimatedDuration time.Duration `json:"estimated_duration" dynamodbav:"estimatedDuration"`
	Description       string        `json:"description" dynamodbav:"description"`
}

type Milestone struct {
	Date        string `json:"date" dynamodbav:"date"`
	Title       string `json:"title" dynamodbav:"title"`
	Type        string `json:"type" dynamodbav:"type"`
	Description string `json:"description" dynamodbav:"description"`
}

type Protection struct {
	Type           string `json:"type" dynamodbav:"type"`
	Premium        int64  `json:"premium" dynamodbav:"premium"`
	ClaimedDate    string `json:"claimed_date" dynamodbav:"claimedDate"`
	ExpirationDate string `json:"expiration_date" dynamodbav:"expirationDate"`
	Description    string `json:"description" dynamodbav:"description"`
}

type Debt struct {
	Type         string `json:"type" dynamodbav:"type"`
	Cost         int64  `json:"cost" dynamodbav:"cost"`
	InterestRate int64  `json:"interest_rate" dynamodbav:"interestRate"`
	Term         int64  `json:"term" dynamodbav:"term"`
	Collateral   string `json:"collateral" dynamodbav:"collateral"`
	Description  string `json:"description" dynamodbav:"description"`
}

type MetaField struct {
	Key       string `json:"key" dynamodbav:"key"`
	Namespace string `json:"namespace" dynamodbav:"namespace"`
	Value     string `json:"value" dynamodbav:"value"`
	Type      string `json:"type" dynamodbav:"type"`
}
