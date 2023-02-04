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
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"net/http"
	"os"
	"time"
	"user-service.kptl.net/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	sdk  struct {
		config aws.Config
		az     string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.sdk.az, "availability-zone", "us-east-1", "AWS Availability Zone")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err := configSdk(&cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(dynamodb.NewFromConfig(cfg.sdk.config)),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func configSdk(cfg *config, logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sdkCfg, err := sdkConfig.LoadDefaultConfig(
		ctx, sdkConfig.WithRegion(cfg.sdk.az), sdkConfig.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	cfg.sdk.config = sdkCfg
	return nil
}
