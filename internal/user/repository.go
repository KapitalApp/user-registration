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
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
	"user-service.mykapital.io/internal/data"
)

// Model is a model that handles CRUD operations for User instances.
// It contains a DynamoDB service client that is used to act on the specified table.
type Model struct {
	// DynamoDbClient is the dynamodb client for User
	DynamoDbClient *dynamodb.Client
	// TableName is the table holding the data for User
	TableName string
}

// TableExists determines whether a DynamoDB table exists.
//
// If the table does not exist, a not found error is returned
// along with false.
func (m Model) TableExists() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DynamoDbClient.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(m.TableName)},
	)

	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			return false, notFoundEx
		}
		return false, fmt.Errorf("couldn't determine existence of table %v. Here's why: %v", m.TableName, err)
	}

	return true, nil
}

// Insert inserts a new user in the table.
//
// If the user already exists, the user get replaced by the new user.
func (m Model) Insert(user *User) error {
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
		return fmt.Errorf("couldn't add item to table. Here's why: %v", err)
	}

	return nil
}

// Get retrieves the user with the specific id.
//
// If no user was found with the given id, nothing will be returned.
func (m Model) Get(id string) (*User, error) {
	user := User{ID: id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := m.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: user.GetKey(), TableName: aws.String(m.TableName),
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get info about %v. Here's why: %v", id, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &user)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal response. Here's why: %v", err)
		}
	}

	return &user, nil
}

// Update updates a user that already exists in the DynamoDB table with the
// new attributes. Current user attributes are not required to be passed.
//
// If the user does not already exist, it adds a new item to the table.
// This function uses the `expression` package to build the update
// expression.
// The Version attribute of the user is automatically updated to handle
// race conditions.
func (m Model) Update(user *User, newAttributes map[string]interface{}) (map[string]interface{}, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]interface{}

	var update expression.UpdateBuilder
	first := true
	for k, v := range newAttributes {
		if first {
			update = expression.Set(expression.Name(k), expression.Value(v))
			first = false
		} else {
			update.Set(expression.Name(k), expression.Value(v))
		}
	}
	update.Set(expression.Name("version"), expression.Value(user.Version+1))

	condition := expression.Name("version").Equal(expression.Value(user.Version))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expression for update. Here's why: %v", err)
	} else {
		response, err = m.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(m.TableName),
			Key:                       user.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ConditionExpression:       expr.Condition(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			var ccf *types.ConditionalCheckFailedException
			switch {
			case errors.As(err, &ccf):
				return nil, data.ErrEditConflict
			default:
				return nil, fmt.Errorf("couldn't update id %v. Here's why: %v", user.ID, err)
			}
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				return nil, fmt.Errorf("couldn't unmarshall update response. Here's why: %v", err)
			}
		}
	}

	return attributeMap, nil
}

// Delete deletes the user from the table in DynamoDB.
//
// The operation is idempotent; running it multiple times on
// the same item or attribute does not result in an error response.
func (m Model) Delete(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DynamoDbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(m.TableName), Key: user.GetKey(),
	})
	if err != nil {
		return fmt.Errorf("couldn't delete %v from the table. Here's why: %v", user.ID, err)
	}

	return nil
}
