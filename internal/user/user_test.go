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

package user

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
	"user-service.mykapital.io/internal/validator"
)

func TestUserGetKey(t *testing.T) {
	tests := map[string]struct {
		input    User
		expected map[string]types.AttributeValue
	}{
		`get primary key`: {
			input: User{ID: "77d1cbe1-f734-4b94-b69e-e9d55b81ed19"},
			expected: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{
					Value: "77d1cbe1-f734-4b94-b69e-e9d55b81ed19",
				},
			},
		},
		`empty primary key`: {
			input: User{ID: ""},
			expected: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{
					Value: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.input.GetKey()

			for key, expectedValue := range tt.expected {
				actualValue, ok := actual[key]
				if !ok {
					t.Errorf("Expected attribute '%v' not found", key)
					continue
				}

				if !attributeValuesEqual(expectedValue, actualValue) {
					t.Errorf("Attribute '%v': Expected '%v', but got '%v'", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// attributeValuesEqual is a helper function to compare two AttributeValues for equality
func attributeValuesEqual(expected, actual types.AttributeValue) bool {
	switch e := expected.(type) {
	case *types.AttributeValueMemberS:
		a, ok := actual.(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		return e.Value == a.Value
	case *types.AttributeValueMemberN:
		a, ok := actual.(*types.AttributeValueMemberN)
		if !ok {
			return false
		}
		return e.Value == a.Value
	default: // Unsupported type
		return false
	}
}

func TestValidateUser(t *testing.T) {
	tests := map[string]struct {
		user     User
		expected map[string]string
	}{
		`valid user`: {
			user: User{
				Email:             "john.doe@example.com",
				FirstName:         "John",
				CountryCodeAlpha2: "US",
				ProvinceCode:      "CA",
				IsMarried:         true,
				Spouse: &FamilyMember{
					Type:      "spouse",
					FirstName: "Jane",
					LastName:  "Doe",
				},
			},
			expected: make(map[string]string),
		},
		`invalid user`: {
			user: User{
				Email:             "invalid-email",
				FirstName:         "",
				CountryCodeAlpha2: "USA",
				ProvinceCode:      "",
				IsMarried:         true,
				Spouse:            nil,
			},
			expected: map[string]string{
				"email":                "must be valid",
				"first_name":           "must be provided",
				"country_code_alpha_2": "must be two letters",
				"province_code":        "must be provided",
				"spouse":               "must be provided",
			},
		},
		`invalid family member`: {
			user: User{
				Email:             "john.doe@example.com",
				FirstName:         "John",
				LastName:          "Doe",
				CountryCodeAlpha2: "US",
				ProvinceCode:      "CA",
				IsMarried:         true,
				Spouse: &FamilyMember{
					Type:     "spouse",
					LastName: "Doe",
				},
				Dependents: []FamilyMember{
					{Type: "dependent"},
					{FirstName: "John", LastName: "Doe"},
				},
			},
			expected: map[string]string{
				"spouse_first_name":      "must be provided",
				"dependent_1_first_name": "must be provided",
				"dependent_2_type":       "must be provided",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockValidator := validator.New()

			ValidateUser(mockValidator, &tt.user)

			for key, expectedErr := range tt.expected {
				if mockValidator.Errors[key] != expectedErr {
					t.Errorf("Expected error '%v' not found", key)
				}
				delete(mockValidator.Errors, key)
			}

			for key, notExpectedErr := range mockValidator.Errors {
				t.Errorf("Unexpected error '%v' with message '%v' found", key, notExpectedErr)
			}
		})
	}
}
