/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
 
package adapter

import (
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/model/harbor"
)

type ScannerAdapter interface {
	GetMetadata() harbor.ScannerAdapterMetadata
	Scan(payload harbor.ScanRequest) (harbor.ScanResponse, error)
	GetImageScanStatus(scanID string) (string, error)
	GetImageVulnerability(scanID string) (harbor.VulnerabilityReport, error)
}
