/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
 
package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/adapter/imagescanning"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/config"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/queue"
	v1 "github.com/vmware/carbon-black-adapter-for-harbor/internal/webapi/v1"
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/scan"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	// start worker for blocking queue
	queue.InitQueue(1000)
	go queue.NewWorker().HandleEvents()

	router := gin.Default()
	register(router)

	log.Fatal(router.Run())
}

func register(router *gin.Engine) {
	handler := scan.NewScanHandler(config.SaasURL(), config.OrgKey(), config.APIID(), config.APIKey(), nil)
	adapter := imagescanning.NewAdapter(*handler)
	group := v1.NewGroup(router, adapter)
	group.Register()
}
