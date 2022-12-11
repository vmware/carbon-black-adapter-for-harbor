/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
 
package queue

import (
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/model/layers"
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/scan"
)

type BomStatus string

const (
	BomGeneratedSuccessfully   BomStatus = "bom-generated-successfully"
	BomGeneratedUnsuccessfully BomStatus = "bom-generated-unsuccessfully"
	BomUploadedSuccessfully    BomStatus = "bom-uploaded-successfully"
	BomUploadedUnSuccessfully  BomStatus = "bom-uploaded-unsuccessfully"
)

type ImageInfo struct {
	Status        BomStatus
	DockerPullTag string
	Digest        string
	FullTag       string
	Bom           *scan.Bom
	layers        []layers.Layer
	UserName      string
	Password      string
	OperationID   string
	ImageID       string
}
