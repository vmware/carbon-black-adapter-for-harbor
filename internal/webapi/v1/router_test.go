/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/model/harbor"
)

type fakeAdapter struct {
	// scan []string
}

func NewFakeAdapter() fakeAdapter {
	return fakeAdapter{
		// scan: make([]string, 0),
	}
}

const (
	nameMetadata    string = "Carbon-Black"
	vendorMetadata  string = "VMware"
	versionMetadata string = "1.0"
)

var (
	router      = gin.Default()
	scannerInfo = harbor.Scanner{
		Name:    nameMetadata,
		Vendor:  vendorMetadata,
		Version: versionMetadata,
	}
	consumesMimeTypesMetadata = []string{
		"application/vnd.docker.distribution.manifest.v2+json",
	}
	producesMimeTypesMetadata = []string{
		"application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0",
	}
	propertiesMetadata = map[string]string{
		"id":    "wert",
		"name":  "somu",
		"hello": "hi",
	}
)

func TestRouterGroup(t *testing.T) {
	adapter := NewFakeAdapter()
	group := NewGroup(router, adapter)
	group.Register()

	req := httptest.NewRequest("GET", "/api/v1/metadata", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Unexpected resp code get, expected: %d, actual: %d", http.StatusOK, resp.Code)
	}
	var respBody harbor.ScannerAdapterMetadata
	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
		t.Errorf("Failed to decode resp body: %v", err)
	}
	// Check resp body
	if respBody.Scanner.Name != nameMetadata {
		t.Errorf("Unexpected scanner name, expected: %v, actual: %v", nameMetadata, respBody.Scanner.Name)
	}

	payloadObj := harbor.ScanRequest{
		Registry: harbor.Registry{
			URL:           "123",
			Authorization: "",
		},
		Artifact: harbor.Artifact{
			Repository: "123",
			Digest:     "sha256:123456",
			Tag:        "123",
			MimeType:   "",
		},
	}
	payloadBuffer, _ := json.Marshal(payloadObj)
	req = httptest.NewRequest("POST", "/api/v1/scan", bytes.NewBuffer(payloadBuffer))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusAccepted {
		t.Errorf("Unexpected resp code get, expected: %d, actual: %d", http.StatusAccepted, resp.Code)
	}

	var respBodyScan harbor.ScanResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &respBodyScan); err != nil {
		t.Errorf("Failed to decode resp body: %v", err)
	}

	if respBodyScan.ID != payloadObj.Artifact.Digest {
		t.Error("Unexpected scan ID ", respBodyScan.ID)
	}

	scanRequestId := "sha256:123456"
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scan/%s/report", scanRequestId), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Unexpected resp code get, expected: %d, actual: %d", http.StatusOK, resp.Code)
	}

	scanRequestId = "wrong123"
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scan/%s/report", scanRequestId), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected resp code get, expected: %d, actual: %d", http.StatusInternalServerError, resp.Code)
	}
}

func (f fakeAdapter) GetMetadata() harbor.ScannerAdapterMetadata {
	capability := harbor.ScannerCapability{
		ConsumesMimeTypes: consumesMimeTypesMetadata,
		ProducesMimeTypes: producesMimeTypesMetadata,
	}

	scannerMetadata := harbor.ScannerAdapterMetadata{
		Properties:   propertiesMetadata,
		Capabilities: []harbor.ScannerCapability{capability},
		Scanner:      scannerInfo,
	}

	return scannerMetadata
}

func (f fakeAdapter) Scan(payload harbor.ScanRequest) (harbor.ScanResponse, error) {
	var resp harbor.ScanResponse
	identifier := payload.Artifact.Digest

	if identifier == "" {
		return resp, fmt.Errorf("Payload error")
	}

	resp = harbor.ScanResponse{ID: identifier}

	return resp, nil
}

func (f fakeAdapter) GetImageScanStatus(scanID string) (string, error) {
	return "FINISHED", nil
}

func (f fakeAdapter) GetImageVulnerability(scanID string) (harbor.VulnerabilityReport, error) {
	if scanID == "sha256:123456" {
		vulnerabilities := make([]harbor.VulnerabilityItem, 1)
		vulnerabilities[0] = harbor.VulnerabilityItem{
			ID:          "12345",
			Package:     "dummyPackage",
			Version:     "dummyVersion",
			FixVersion:  "dummyFixAvailable",
			Severity:    harbor.ToHarborSeverity("critical"),
			Description: "dummyDescription",
			Links:       []string{"dummyDescription"},
		}
		return harbor.VulnerabilityReport{
			GeneratedAt: time.Now(),
			Artifact: harbor.Artifact{
				Repository: "dummyRepo",
				Digest:     "dummyDigest",
				Tag:        "dummyTag",
				MimeType:   "world",
			},
			Scanner:         scannerInfo,
			Severity:        harbor.ToHarborSeverity("critical"),
			Vulnerabilities: vulnerabilities,
		}, nil
	}

	return harbor.VulnerabilityReport{}, fmt.Errorf("image ID not found")
}
