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
	"expvar"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
	"runtime"
	"time"
	"user-service.kptl.net/internal/data"
	"user-service.kptl.net/internal/jsonlog"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	sdk  struct {
		config aws.Config
		az     string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.sdk.az, "availability-zone", "us-east-1", "AWS Availability Zone")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	err := configSdk(&cfg, logger)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(dynamodb.NewFromConfig(cfg.sdk.config)),
	}

	err = app.serve(logger)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func configSdk(cfg *config, logger *jsonlog.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sdkCfg, err := sdkConfig.LoadDefaultConfig(
		ctx,
		sdkConfig.WithRegion(cfg.sdk.az),
		sdkConfig.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	if cfg.env == "development" {
		sdkCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})
	}

	cfg.sdk.config = sdkCfg
	return nil
}
