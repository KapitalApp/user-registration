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

// Package data is an abstraction on top of the data models and
// functionalities.
package data

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"user-service.mykapital.io/internal/user"
)

// Possible errors passed from a model.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models represents the internal models for the server.
type Models struct {
	Users user.Model
}

// NewModels creates Models.
//
// For the user model, a DynamoDB client is passed.
func NewModels(client *dynamodb.Client) Models {
	return Models{
		Users: user.Model{DynamoDbClient: client, TableName: "User"},
	}
}
