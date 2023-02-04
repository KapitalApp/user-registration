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
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (user User) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(user.ID)
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{"id": id}
}

type UserModel struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func (m UserModel) TableExists() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DynamoDbClient.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(m.TableName)},
	)

	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			return false, fmt.Errorf("Table %v does not exist.\n", m.TableName)
		}
		return false, fmt.Errorf("Couldn't determine existence of table %v. Here's why: %v\n", m.TableName, err)
	}

	return true, nil
}

func (m UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		panic(err)
	}
	_, err = m.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(m.TableName), Item: item,
	})
	if err != nil {
		return fmt.Errorf("Couldn't add item to table. Here's why: %v\n", err)
	}

	return nil
}

func (m UserModel) Get(id uuid.UUID) (*User, error) {
	user := User{ID: id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := m.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: user.GetKey(), TableName: aws.String(m.TableName),
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't get info about %v. Here's why: %v\n", id, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &user)
		if err != nil {
			return nil, fmt.Errorf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User, newAttributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}

	var update expression.UpdateBuilder
	for k, v := range newAttributes {
		update.Set(expression.Name(k), expression.Value(v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, fmt.Errorf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err = m.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(m.TableName),
			Key:                       user.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			return nil, fmt.Errorf("Couldn't update id %v. Here's why: %v\n", user.ID, err)
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				return nil, fmt.Errorf("Couldn't unmarshall update response. Here's why: %v\n", err)
			}
		}
	}

	return attributeMap, nil
}

func (m UserModel) Delete(user *User) error {
	_, err := m.DynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(m.TableName), Key: user.GetKey(),
	})
	if err != nil {
		return fmt.Errorf("Couldn't delete %v from the table. Here's why: %v\n", user.ID, err)
	}
	return nil
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be valid")
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(user.LastName != "", "last_name", "must be provided")
	v.Check(len(user.CountryCodeAlpha2) == 2, "country_code_alpha_2", "must be two letters")
	v.Check(user.ProvinceCode != "", "province_code", "must be provided")
}
