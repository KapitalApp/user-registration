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
	Type        string `json:"type"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"age"`
	Income      int64  `json:"income"`
	Expenses    int64  `json:"expenses"`
}

type Goal struct {
	Date              string        `json:"date"`
	Title             string        `json:"title"`
	ProgressLevel     string        `json:"progress_level"`
	EstimatedDuration time.Duration `json:"estimated_duration"`
	Description       string        `json:"description"`
}

type Milestone struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Protection struct {
	Type           string `json:"type"`
	Premium        int64  `json:"premium"`
	ClaimedDate    string `json:"claimed_date"`
	ExpirationDate string `json:"expiration_date"`
	Description    string `json:"description"`
}

type Debt struct {
	Type         string `json:"type"`
	Cost         string `json:"cost"`
	InterestRate int64  `json:"interest_rate"`
	Term         int64  `json:"term"`
	Collateral   string `json:"collateral"`
	Description  string `json:"description"`
}

type MetaField struct {
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
	Value     string `json:"value"`
	Type      string `json:"type"`
}
