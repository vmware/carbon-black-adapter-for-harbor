/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
 
package harbor

import (
	"testing"
)

var (
	dummyRequest ScanRequest
)

func TestValidate(t *testing.T) {
	dummyRequest.Registry.URL = "core.harbor.domain"
	dummyRequest.Artifact.Digest = "sha256:c8d0cdef8c20b68dae65db25924bdcb620c1f121f5f61f4a6f0e402e8d070af0"
	dummyRequest.Artifact.Repository = "demo/postgres"

	result := dummyRequest.Validate()
	if result != nil {
		t.Errorf("Validation was not successful")
	}
}

func TestValidateError(t *testing.T) {
	dummyRequest.Registry.URL = ""
	dummyRequest.Artifact.Digest = "sha256:c8d0cdef8c20b68dae65db25924bdcb620c1f121f5f61f4a6f0e402e8d070af0"
	dummyRequest.Artifact.Repository = "demo/postgres"

	result := dummyRequest.Validate()
	if result == nil {
		t.Errorf("Validation was not successful")
	}
}
