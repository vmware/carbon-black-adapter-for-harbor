/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
 
package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var config Config

type Config struct {
	orgKey  string
	saasURL string
	apiID   string
	apiKey  string
}

func InitConfig() error {
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = "debug"
	}
	log.Infof("log level was set to %s", level)
	ll, err := log.ParseLevel(level)
	if err != nil {
		ll = log.DebugLevel
	}
	// set log level
	log.SetLevel(ll)

	url := os.Getenv("CB_URL")
	if url == "" {
		return fmt.Errorf("no URL found in config")
	}
	config.saasURL = url

	orgKey := os.Getenv("CB_ORG_KEY")
	if orgKey == "" {
		return fmt.Errorf("no Org Key found in config")
	}
	config.orgKey = orgKey

	apiID := os.Getenv("CB_API_ID")
	if apiID == "" {
		return fmt.Errorf("no API ID found in config")
	}
	config.apiID = apiID

	apiKey := os.Getenv("CB_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("no API key found in config")
	}
	config.apiKey = apiKey

	log.Info("all the configurations are set")

	return nil
}

func OrgKey() string {
	return config.orgKey
}

func SaasURL() string {
	return config.saasURL
}

func APIID() string {
	return config.apiID
}

func APIKey() string {
	return config.apiKey
}
