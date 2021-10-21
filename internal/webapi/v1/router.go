/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/scan"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/adapter"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/model/harbor"
)

type APIsGroup struct {
	*gin.RouterGroup
	adapter adapter.ScannerAdapter
}

// NewGroup returns a new router.
func NewGroup(router *gin.Engine, adapter adapter.ScannerAdapter) APIsGroup {
	group := router.Group("api/v1")

	return APIsGroup{
		RouterGroup: group,
		adapter:     adapter,
	}
}

func (g APIsGroup) Register() {
	g.RouterGroup.GET("", g.healthCheck)
	g.RouterGroup.GET("/metadata", g.getMetadata)
	g.RouterGroup.POST("/scan", g.acceptScanRequest)
	g.RouterGroup.GET("/scan/:scan_request_id/report", g.getScanReport)
}

// TODO: add checking to backend here
func (g APIsGroup) healthCheck(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (g APIsGroup) getMetadata(c *gin.Context) {
	log.Debug("Getting metadata was successful")
	c.JSON(http.StatusOK, g.adapter.GetMetadata())
}

func (g APIsGroup) acceptScanRequest(c *gin.Context) {
	var scanPayload harbor.ScanRequest
	if err := c.BindJSON(&scanPayload); err != nil {
		c.JSON(http.StatusBadRequest, harbor.NewErrorResponse("Cannot bind request", err))
		return
	}

	if err := scanPayload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, harbor.NewErrorResponse("Invalid request", err))
		return
	}

	scanResponse, err := g.adapter.Scan(scanPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, harbor.NewErrorResponse("Scan failed", err))
		return
	}

	c.JSON(http.StatusAccepted, scanResponse)
}

func (g APIsGroup) getScanReport(c *gin.Context) {
	scanID := c.Param("scan_request_id")
	log.Infof("Checking report generation status: %s", scanID)

	status, err := g.adapter.GetImageScanStatus(scanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, harbor.NewErrorResponse("Failed to get result", err))
		return
	}

	switch status {
	case string(scan.FinishedStatus):
		report, err := g.adapter.GetImageVulnerability(scanID)
		if err != nil {
			log.Errorf("Failed to get scan report: %v", err)
			c.JSON(http.StatusInternalServerError, harbor.NewErrorResponse("Failed to get result", err))
		}
		c.JSON(http.StatusOK, report)
	case string(scan.FailedStatus):
		c.JSON(http.StatusInternalServerError, harbor.NewErrorResponse("Failed to get result", err))
	case string(scan.QueuedStatus):
		fallthrough
	default:
		c.Writer.Header().Set("Refresh-After", "30")
		c.Redirect(302, c.Request.URL.String())
		log.Info("Scan report status is QUEUED trying after 30 seconds")
		c.Writer.WriteHeader(http.StatusFound)
	}
}
