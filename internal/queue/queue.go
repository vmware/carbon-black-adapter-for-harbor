/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package queue

import (
	"github.com/google/uuid"
)

var (
	// queue will store an unique identifier for the scan task
	queue chan string
	// imageInfoMap will store all the information for the images about to scan
	imageInfoMap map[string]*ImageInfo
)

func InitQueue(bufferSize int) {
	queue = make(chan string, bufferSize)
	imageInfoMap = make(map[string]*ImageInfo)
}

func Queue() chan string {
	return queue
}

func Publish(imageInfo ImageInfo) string {
	scanID := uuid.New().String()
	imageInfoMap[scanID] = &imageInfo
	queue <- scanID
	return scanID
}

func Fetch(scanID string) (*ImageInfo, bool) {
	v, ok := imageInfoMap[scanID]
	return v, ok
}

func Remove(scanID string) {
	delete(imageInfoMap, scanID)
}
