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
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"testing"
)

func TestReadParam(t *testing.T) {
	app := &application{}

	tests := []struct {
		name     string
		url      string
		param    string
		expected string
	}{
		{
			name:     "Test case 1: Check if the function correctly returns the value of the parameter",
			url:      "/test?param=value",
			param:    "param",
			expected: "value",
		},
		{
			name:     "Test case 2: Check if the function returns an empty string when the parameter does not exist",
			url:      "/test",
			param:    "param",
			expected: "",
		},
		{
			name:     "Test case 3: Check if the function handles case sensitivity correctly",
			url:      "/test?param=value",
			param:    "PARAM",
			expected: "",
		},
		{
			name:     "Test case 4: Check if the function correctly handles numeric parameters",
			url:      "/test?param=123",
			param:    "param",
			expected: "123",
		},
		{
			name:     "Test case 5: Check if the function correctly handles empty parameters",
			url:      "/test?param=",
			param:    "param",
			expected: "",
		},
		{
			name:     "Test case 6: Check if the function correctly handles context with no parameters",
			url:      "/test",
			param:    "param",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var req *http.Request
			var err error
			if test.url != "" {
				req, err = http.NewRequest("GET", test.url, nil)
				if err != nil {
					t.Fatal(err)
				}
			}

			params := httprouter.Params{}
			if test.expected != "" {
				params = httprouter.Params{httprouter.Param{Key: test.param, Value: test.expected}}
			}

			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			if req != nil {
				req = req.WithContext(ctx)
			}
			result := app.readParam(req, test.param)
			if result != test.expected {
				t.Errorf("Expected: %s, but got: %s", test.expected, result)
			}
		})
	}
}
