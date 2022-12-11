/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package queue

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vmware/carbon-black-adapter-for-harbor/internal/config"
	"github.com/vmware/carbon-black-cloud-container-cli/pkg/scan"
)

type Worker struct{}

func NewWorker() Worker {
	return Worker{}
}

func (w Worker) HandleEvents() {
	chanForUploadingBom := make(chan string)

	go func() {
		for scanID := range chanForUploadingBom {
			imageInfo, ok := Fetch(scanID)
			if !ok {
				log.Errorf("Cannot get information for task: %v", scanID)
				continue
			}

			if imageInfo.Status != BomGeneratedSuccessfully {
				log.Errorf("No valid bom detected for task: %s", scanID)
				continue
			}

			scanHandler := scan.NewScanHandler(config.SaasURL(), config.OrgKey(), config.APIID(), config.APIKey(), imageInfo.Bom, imageInfo.layers)

			operationID := uuid.New().String()
			log.WithField("operation_id", operationID).Info("Starting an operation")

			imageInfo.OperationID = operationID

			if _, err := scanHandler.PutBomAndLayersToAnalysisAPI(operationID, scan.Option{ForceScan: true}); err != nil {
				log.Errorf("Error putting BOM for analysis of image %v", err)

				imageInfo.Status = BomUploadedUnSuccessfully
				continue
			}

			imageInfo.Status = BomUploadedSuccessfully
		}
	}()

	for scanID := range Queue() {
		imageInfo, ok := Fetch(scanID)
		if !ok {
			log.Errorf("Cannot get information for task: %v", scanID)
			continue
		}

		scanner := scan.NewScanner()

		opts := scan.Option{
			FullTag:    imageInfo.FullTag,
			Credential: fmt.Sprintf("%v:%v", imageInfo.UserName, imageInfo.Password)}
		bomGenerated, imgLayers, hasErr := scanner.ExtractDataFromImage(imageInfo.DockerPullTag, opts)

		if hasErr {
			log.Errorf("Error generating bom for the image %v:", imageInfo)
			imageInfo.Status = BomGeneratedUnsuccessfully
			continue
		}

		imageInfo.layers = imgLayers
		imageInfo.Bom = bomGenerated
		imageInfo.Digest = bomGenerated.ManifestDigest
		imageInfo.Status = BomGeneratedSuccessfully
		chanForUploadingBom <- scanID
	}
}
