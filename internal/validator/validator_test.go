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

package validator

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestValidator(t *testing.T) {
	tests := map[string]struct {
		addErrors []string
		checks    []struct {
			ok           bool
			key, message string
		}
		expectedErrors map[string]string
		expectedValid  bool
	}{
		`no errors added`: {
			addErrors:      nil,
			checks:         nil,
			expectedErrors: make(map[string]string),
			expectedValid:  true,
		},
		`errors added`: {
			addErrors: []string{
				"key1:message1",
				"key2:message2",
			},
			checks: nil,
			expectedErrors: map[string]string{
				"key1": "message1",
				"key2": "message2",
			},
			expectedValid: false,
		},
		`errors added and but check passed`: {
			addErrors: []string{
				"key1:message1",
				"key2:message2",
			},
			checks: []struct {
				ok           bool
				key, message string
			}{
				{true, "key3", "message3"},
			},
			expectedErrors: map[string]string{
				"key1": "message1",
				"key2": "message2",
			},
			expectedValid: false,
		},
		`checks fail`: {
			addErrors: nil,
			checks: []struct {
				ok           bool
				key, message string
			}{
				{false, "key1", "message1"},
				{false, "key2", "message2"},
			},
			expectedErrors: map[string]string{
				"key1": "message1",
				"key2": "message2",
			},
			expectedValid: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			v := New()

			for _, err := range tt.addErrors {
				key, message := parseErrorString(err)
				v.AddError(key, message)
			}

			for _, check := range tt.checks {
				v.Check(check.ok, check.key, check.message)
			}

			if len(v.Errors) != len(tt.expectedErrors) {
				t.Errorf("unexpected error count: got %d, want %d", len(v.Errors), len(tt.expectedErrors))
			}

			for key, want := range tt.expectedErrors {
				got, ok := v.Errors[key]
				if !ok {
					t.Errorf("missing expected error for key '%s': got nothing, want '%s'", key, want)
				} else if got != want {
					t.Errorf("unexpected error message for key '%s': got '%s', want '%s'", key, got, want)
				}
			}

			if v.Valid() != tt.expectedValid {
				t.Errorf("unexpected validation error: validation should be '%t', is '%t'", tt.expectedValid, v.Valid())
			}
		})
	}
}

// parseErrorString parses the key and the message from error
func parseErrorString(err string) (key string, message string) {
	parsedString := strings.Split(err, ":")
	key, message = parsedString[0], parsedString[1]
	return
}

func TestIn(t *testing.T) {
	testCases := map[string]struct {
		value     string
		list      []string
		wantFound bool
	}{
		`value in list`: {
			value:     "apple",
			list:      []string{"orange", "apple", "banana"},
			wantFound: true,
		},
		`value not in list`: {
			value:     "grape",
			list:      []string{"orange", "apple", "banana"},
			wantFound: false,
		},
		`empty list`: {
			value:     "apple",
			list:      []string{},
			wantFound: false,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			gotFound := In(tt.value, tt.list...)
			if gotFound != tt.wantFound {
				t.Errorf("In(%q, %q) = %v; want %v", tt.value, tt.list, gotFound, tt.wantFound)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	testCases := map[string]struct {
		value    string
		pattern  string
		expected bool
	}{
		`valid email address`: {
			value:    "test@example.com",
			pattern:  EmailRX.String(),
			expected: true,
		},
		`invalid email address`: {
			value:    "test.example.com",
			pattern:  EmailRX.String(),
			expected: false,
		},
		`valid phone number`: {
			value:    "123-456-7890",
			pattern:  "^[0-9]{3}-[0-9]{3}-[0-9]{4}$",
			expected: true,
		},
		`invalid phone number`: {
			value:    "123-456-789",
			pattern:  "^[0-9]{3}-[0-9]{3}-[0-9]{4}$",
			expected: false,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			rx, err := regexp.Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern %s: %v", tt.pattern, err)
			}

			if Matches(tt.value, rx) != tt.expected {
				t.Errorf("Matches() = %v, expected %v", !tt.expected, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := map[string]struct {
		input    []string
		expected bool
	}{
		`unique`: {
			input:    []string{"apple", "banana", "orange"},
			expected: true,
		},
		`not unique consecutive`: {
			input:    []string{"apple", "banana", "banana"},
			expected: false,
		},
		`not unique empty strings`: {
			input:    []string{"", ""},
			expected: false,
		},
		`not unique non consecutive`: {
			input:    []string{"apple", "banana", "orange", "apple"},
			expected: false,
		},
		`unique empty`: {
			input:    []string{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%v", tt.input), func(t *testing.T) {
			result := Unique(tt.input)
			if result != tt.expected {
				t.Errorf("expected Unique(%v) to be %v, but got %v", tt.input, tt.expected, result)
			}
		})
	}
}
