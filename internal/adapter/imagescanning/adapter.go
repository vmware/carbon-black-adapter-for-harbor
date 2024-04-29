/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package imagescanning

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/model/harbor"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/queue"
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/model/image"
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/scan"
)

const (
	nameMetadata    string = "Carbon-Black"
	vendorMetadata  string = "VMware"
	versionMetadata string = "1.0"
)

var (
	scannerInfo = harbor.Scanner{
		Name:    nameMetadata,
		Vendor:  vendorMetadata,
		Version: versionMetadata,
	}
	consumesMimeTypesMetadata = []string{
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
	}
	producesMimeTypesMetadata = []string{
		"application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0",
		"application/vnd.security.vulnerability.report; version=1.1",
	}
	propertiesMetadata = map[string]string{
		"harbor.scanner-adapter/scanner-type": "os-package-vulnerability",
		"env.LOG_LEVEL":                       log.GetLevel().String(),
	}
)

type Adapter struct {
	getVulnHandler scan.Handler
}

func NewAdapter(handler scan.Handler) *Adapter {
	return &Adapter{getVulnHandler: handler}
}

func (a Adapter) GetMetadata() harbor.ScannerAdapterMetadata {
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

func (a Adapter) Scan(payload harbor.ScanRequest) (harbor.ScanResponse, error) {
	var resp harbor.ScanResponse

	u, err := url.Parse(payload.Registry.URL)
	if err != nil {
		return resp, err
	}

	digestPullString := fmt.Sprintf("%s/%s@%s",
		u.Host, payload.Artifact.Repository, payload.Artifact.Digest)

	username, password, err := getCredential(payload.Registry.Authorization)
	if err != nil {
		log.Errorf("Error getting username and password from request: %v", err)
		return resp, err
	}

	tag := payload.Artifact.Tag
	if payload.Artifact.Tag == "" {
		tag = payload.Artifact.Digest[7:]
	}

	fullTag := fmt.Sprintf("%s/%s:%s", u.Host, payload.Artifact.Repository, tag)
	if ref, err := reference.ParseAnyReference(fullTag); err == nil {
		fullTag = ref.String()
	}

	log.Infof("Image Information: PullTag: %s, FullTag: %s", digestPullString, fullTag)

	scanID := queue.Publish(queue.ImageInfo{
		DockerPullTag: digestPullString,
		FullTag:       fullTag,
		UserName:      username,
		Password:      password,
	})

	resp = harbor.ScanResponse{ID: scanID}
	return resp, nil
}

func getCredential(rawCredential string) (string, string, error) {
	auth := strings.Split(strings.TrimSpace(rawCredential), " ")
	if len(auth) < 2 {
		return "", "", fmt.Errorf("raw credential is incorrect")
	}

	userPass := auth[1]
	decodedUserPass, err := base64.StdEncoding.DecodeString(userPass)
	if err != nil {
		return "", "", err
	}

	splitUserPass := strings.SplitN(string(decodedUserPass), ":", 2)
	return splitUserPass[0], splitUserPass[1], nil
}

func (a Adapter) GetImageScanStatus(scanID string) (string, error) {
	imageInfo, ok := queue.Fetch(scanID)
	if !ok {
		return "", fmt.Errorf("cannot fetch status for task %v in the queue", scanID)
	}

	switch imageInfo.Status {
	case queue.BomUploadedUnSuccessfully, queue.BomGeneratedUnsuccessfully:
		defer queue.Remove(scanID)
		return "", fmt.Errorf("failed to generate/upload sbom to backend: %s", imageInfo.Status)
	case queue.BomUploadedSuccessfully:
		status, err := a.getVulnHandler.GetImageAnalysisStatus(imageInfo.Digest, imageInfo.OperationID)
		return string(status.OperationStatus), err
	case queue.BomGeneratedSuccessfully:
		fallthrough
	default:
		return "UNSTARTED", nil
	}
}

func (a Adapter) GetImageVulnerability(scanID string) (harbor.VulnerabilityReport, error) {
	// we will only fetch data once the status is finished, so there is always a valid info for us
	imageInfo, _ := queue.Fetch(scanID)
	digest := imageInfo.Digest
	defer queue.Remove(scanID)

	scannedResult, err := a.getVulnHandler.GetImageVulnerability(digest, "", "")
	if err != nil {
		return harbor.VulnerabilityReport{}, err
	}
	vul := convertVulnerability(*scannedResult)

	return vul, nil
}

func convertVulnerability(scannedResult image.ScannedImage) harbor.VulnerabilityReport {

	maxSeverity := harbor.UNKNOWN
	vulnerabilities := make([]harbor.VulnerabilityItem, len(scannedResult.Vulnerabilities))
	for i, v := range scannedResult.Vulnerabilities {
		vulnerabilities[i] = harbor.VulnerabilityItem{
			ID:          v.ID,
			Package:     v.Package,
			Version:     v.Version,
			FixVersion:  v.FixAvailable,
			Severity:    harbor.ToHarborSeverity(v.Severity),
			Description: v.Description,
			Links:       []string{v.Link},
		}

		if vulnerabilities[i].FixVersion == "None" {
			vulnerabilities[i].FixVersion = "No fix available"
		}

		maxSeverity = harbor.MaxSeverity(maxSeverity, vulnerabilities[i].Severity)
	}

	return harbor.VulnerabilityReport{
		GeneratedAt: time.Now(),
		Artifact: harbor.Artifact{
			Repository: scannedResult.Repo,
			Digest:     scannedResult.ManifestDigest,
			Tag:        scannedResult.Tag,
			MimeType:   "application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0",
		},
		Scanner:         scannerInfo,
		Severity:        maxSeverity,
		Vulnerabilities: vulnerabilities,
	}
}
