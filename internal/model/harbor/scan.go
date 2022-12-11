/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package harbor

import (
	"fmt"
	"time"
)

type ScanRequest struct {
	Registry Registry `json:"registry"`
	Artifact Artifact `json:"artifact"`
}

func (r ScanRequest) Validate() error {
	if r.Registry.URL == "" {
		return fmt.Errorf("empty registry url")
	}

	if r.Artifact.Digest == "" {
		return fmt.Errorf("empty digest")
	}

	if r.Artifact.Repository == "" {
		return fmt.Errorf("empty repository")
	}

	return nil
}

type ScanResponse struct {
	ID string `json:"id"`
}

type Artifact struct {
	// The name of the Docker Registry repository containing the artifact.
	Repository string `json:"repository,omitempty"`
	// The artifact's digest, consisting of an algorithm and hex portion.
	Digest string `json:"digest,omitempty"`
	// The artifact's tag
	Tag string `json:"tag,omitempty"`
	// The MIME type of the artifact.
	MimeType string `json:"mime_type,omitempty"`
}

type Registry struct {
	// A base URL or the Docker Registry v2 API.
	URL string `json:"url,omitempty"`
	// An optional value of the HTTP Authorization header sent with each request to
	// the Docker Registry v2 API. It's used to exchange Base64 encoded robot account credentials
	// to a short lived JWT access token which allows the underlying scanner to pull the artifact
	// from the Docker Registry.
	Authorization string `json:"authorization,omitempty"`
}

type VulnerabilityReport struct {
	GeneratedAt     time.Time           `json:"generated_at,omitempty"`
	Artifact        Artifact            `json:"artifact,omitempty"`
	Scanner         Scanner             `json:"scanner,omitempty"`
	Severity        Severity            `json:"severity,omitempty"`
	Vulnerabilities []VulnerabilityItem `json:"vulnerabilities,omitempty"`
}
