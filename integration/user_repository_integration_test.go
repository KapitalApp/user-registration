//go:build integration

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

package integration

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbtype "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	xerrors "user-service.mykapital.io/internal/errors"
	"user-service.mykapital.io/internal/user"
)

func TestRepositoryIntegration(t *testing.T) {
	scenarioSteps := []struct {
		name string
		test func(t *testing.T, model user.Model)
	}{
		{`create a new table then confirm the table exists`, testNewTable},
		{`add a new item and get it back to confirm the operation`, testNewItem},
		{`update the new item and get it back to confirm the operation`, testUpdateItem},
		{`try to update the new item, but get an edit conflict error`, testEditConflict},
		{`remove the item and confirm the item is removed`, testRemoveItem},
		{`remove the table and confirm the table is removed`, testRemoveTable},
	}
	runTestsOnDynamoDB(t, scenarioSteps)
}

func testNewTable(t *testing.T, model user.Model) {
	_, err := model.CreateTable()
	if err != nil {
		t.Fatalf("table %s is not created: %v", model.TableName, err)
	}

	exists, err := model.TableExists()
	if err != nil {
		t.Fatalf("table search was intrupted: %v", err)
	}

	if !exists {
		t.Errorf("table %s is created but was not found", model.TableName)
	}
}

func testNewItem(t *testing.T, model user.Model) {
	usr := user.User{
		ID:                     "f8ae3ad1-d5c7-4465-b446-2e931606e938",
		Email:                  "john.doe@example.com",
		FirstName:              "John",
		LastName:               "Doe",
		ProvinceCode:           "CA",
		CountryCodeAlpha2:      "US",
		AdministrativeDivision: "state",
		Currency:               "USD",
		CreatedAt:              time.Now().Format("2006-01-02"),
		Version:                1,
	}

	err := model.Insert(&usr)
	if err != nil {
		t.Fatalf("failed to insert user into %s: %v", model.TableName, err)
	}

	response, err := model.Get(usr.ID)
	if err != nil {
		t.Fatalf("failed to get user from %s: %v", model.TableName, err)
	}

	require.EqualValuesf(t, usr, *response, "user inserted into the table, but was not retrieved")
}

func testUpdateItem(t *testing.T, model user.Model) {
	usr, err := model.Get("f8ae3ad1-d5c7-4465-b446-2e931606e938")
	if err != nil {
		t.Fatalf("failed to get user from %s: %v", model.TableName, err)
	}
	newAttributes := map[string]interface{}{
		"milestones": []map[string]string{
			{"date": "2023-02-05", "title": "Bank Opened", "type": "Debt", "description": ""},
		},
	}

	_, err = model.Update(usr, newAttributes)
	if err != nil {
		t.Fatalf("failed to update the user in %s: %v", model.TableName, err)
	}

	response, err := model.Get(usr.ID)
	if err != nil {
		t.Fatalf("failed to get user from %s: %v", model.TableName, err)
	}

	requiredMileStone := []user.Milestone{
		{Date: "2023-02-05", Title: "Bank Opened", Type: "Debt", Description: ""},
	}
	require.EqualValuesf(t, requiredMileStone, response.Milestones, "failed to confirm that the user is updated")
}

func testEditConflict(t *testing.T, model user.Model) {
	usrWithID := user.User{ID: "f8ae3ad1-d5c7-4465-b446-2e931606e938"}
	newAttributes := map[string]interface{}{
		"milestones": []map[string]string{
			{"date": "2023-02-06", "title": "Bank Opened", "type": "Debt", "description": ""},
		},
	}

	usr, err := model.Get(usrWithID.ID)
	if err != nil {
		t.Fatalf("failed to get user from %s: %v", model.TableName, err)
	}

	usr.Version -= 1

	_, err = model.Update(usr, newAttributes)
	if !errors.Is(err, xerrors.ErrEditConflict) {
		t.Errorf("No edit conflicts detected")
	}
}

func testRemoveItem(t *testing.T, model user.Model) {
	err := model.Delete(&user.User{ID: "f8ae3ad1-d5c7-4465-b446-2e931606e938"})
	if err != nil {
		t.Fatalf("failed to delete user from %s: %v", model.TableName, err)
	}

	usr, err := model.Get("f8ae3ad1-d5c7-4465-b446-2e931606e938")
	if err != nil {
		t.Fatalf("failed to get user from %s: %v", model.TableName, err)
	}

	require.EqualValuesf(t, &user.User{}, usr, "failed to confirm that the user is deleted")
}

func testRemoveTable(t *testing.T, model user.Model) {
	err := model.DeleteTable()
	if err != nil {
		t.Fatalf("failed to delete table %s: %v", model.TableName, err)
	}

	var notFoundEx *dbtype.ResourceNotFoundException

	exists, err := model.TableExists()
	if !errors.As(err, &notFoundEx) {
		t.Fatalf("table search was intrupted: %v", err)
	}

	if exists {
		t.Errorf("table was not removed")
	}
}

// runTestsOnDynamoDB checks the test scenarios.
//
// Clients are created for DynamoDB and Docker. The Docker daemon
// should be running.
func runTestsOnDynamoDB(t *testing.T, scenarioSteps []struct {
	name string
	test func(t *testing.T, model user.Model)
}) {
	// Set up Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("failed to set up docker client: %v", err)
	}

	// Start DynamoDB container
	containerID, err := startDynamoDBContainer(dockerClient)
	if err != nil {
		t.Fatalf("failed to run dynamodb container: %v", err)
	}

	// Set up DynamoDB client
	dynamodbClient, err := getDynamoDBClient()
	if err != nil {
		t.Fatalf("failed to set up dynamodb client: %v", err)
	}

	// Run the scenario steps
	var mutex sync.Mutex
	prevStatus := true
	model := user.Model{
		DynamoDbClient: dynamodbClient,
		TableName:      "User",
		IndexName:      "email",
	}
	for _, step := range scenarioSteps {
		t.Run(step.name, func(t *testing.T) {
			mutex.Lock()
			defer mutex.Unlock()

			if !prevStatus {
				t.Skip("previous test failed")
			}

			step.test(t, model)
			prevStatus = !t.Failed()
		})
	}

	_ = containerID
	// Stop DynamoDB container
	err = stopDynamoDBContainer(dockerClient, containerID)
	if err != nil {
		t.Fatalf("failed to stop dynamodb container: %v", err)
	}
}

// StartDynamoDBContainer starts a DynamoDB container and returns
// its ID.
//
// Docker daemon must be up and running. Image name for DynamoDB
// is "amazon/dynamodb-local" with tag "latest".
func startDynamoDBContainer(cli *client.Client) (string, error) {
	// Check if Docker is running
	if _, err := cli.Info(context.Background()); err != nil {
		return "", fmt.Errorf("docker daemon is not running: %v", err)
	}

	// Pull Docker image
	imageName := "amazon/dynamodb-local:latest"
	out, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("docker image '%s' was not retreived: %v", imageName, err)
	}
	defer func(out io.ReadCloser) {
		_ = out.Close()
	}(out)
	if _, err := io.Copy(os.Stdout, out); err != nil {
		return "", fmt.Errorf("docker image pull is not closed: %v", err)
	}

	// Run Docker container
	containerConfig := &container.Config{
		Image: imageName,
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8000/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "8000",
				},
			},
		},
	}
	resp, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, "dynamodb-local")
	if err != nil {
		return "", fmt.Errorf("docker failed to creat a container from the image: %v", err)
	}
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("docker failed to start the container: %v", err)
	}

	// Check container logs
	containerLogs, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("docker check for the container logs: %v", err)
	}
	defer func(containerLogs io.ReadCloser) {
		_ = containerLogs.Close()
	}(containerLogs)
	if _, err := io.Copy(os.Stdout, containerLogs); err != nil {
		return "", fmt.Errorf("docker failed to close the container logs: %v", err)
	}

	return resp.ID, nil
}

// StopDynamoDBContainer stops and removes a DynamoDB container.
func stopDynamoDBContainer(cli *client.Client, containerID string) error {
	// Stop the container
	timeout := int(time.Second.Milliseconds()) * 10
	if err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("docker failed to stop the container: %v", err)
	}

	// Remove the container
	if err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("docker failed to remove the container: %v", err)
	}

	return nil
}

// getDynamoDBClient returns a new dynamodb client.
//
// The client points to the running local dynamodb.
func getDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				"FAKE_ACCESS_KEY_ID",
				"FAKE_SECRET_ACCESS_KEY",
				"")),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to configure go aws sdk: %v", err)
	}
	return dynamodb.NewFromConfig(cfg), nil
}
