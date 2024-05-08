/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
package imagescanning

 import (
	 "encoding/base64"
	 "fmt"
	 "testing"
 
	 "github.com/google/go-cmp/cmp"
	 log "github.com/sirupsen/logrus"
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
 
 func (a fakeAdapter) fakeGetMetadata() harbor.ScannerAdapterMetadata {
	 // this duplicates the Adapter.GetMetadata method
	 // but should use the existing data from Adapter
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
 
 func TestGetCredential(t *testing.T) {
 
	 username := "admin"
	 password := "password"
	 beforeEncoding := fmt.Sprintf("%s:%s", username, password)
	 encodedUserPass := base64.StdEncoding.EncodeToString([]byte(beforeEncoding))
	 encodedCred := fmt.Sprintf("Basic %s", encodedUserPass)
 
	 user, pass, _ := getCredential(encodedCred)
 
	 if user != username {
		 t.Errorf("error getting username")
	 }
 
	 if pass != password {
		 t.Errorf("Pulled getting password")
	 }
 
 }
 
 func TestConsumesMimeTypesMetadata(t *testing.T) {
	 var expectConsumesMimeTypesMetadata = []string{
		 "application/vnd.docker.distribution.manifest.v2+json",
		 "application/vnd.oci.image.manifest.v1+json",
	 }
 
	 ad := NewFakeAdapter()
	 md := ad.fakeGetMetadata()
 
	 if !cmp.Equal(md.Capabilities[0].ConsumesMimeTypes, expectConsumesMimeTypesMetadata) {
		 t.Errorf("ConsumesMimeTypesMetadata did not match expected value %s", expectConsumesMimeTypesMetadata)
	 }
 }
 
 func TestProducesMimeTypesMetadata(t *testing.T) {
	 var expectProducesMimeTypesMetadata = []string{
		 "application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0",
		 "application/vnd.security.vulnerability.report; version=1.1",
	 }
 
	 ad := NewFakeAdapter()
	 md := ad.fakeGetMetadata()
 
	 if !cmp.Equal(md.Capabilities[0].ProducesMimeTypes, expectProducesMimeTypesMetadata) {
		 t.Errorf("ProducesMimeTypesMetadata did not match expected value %s", expectProducesMimeTypesMetadata)
	 }
 }
 
 func TestPropertiesMetadata(t *testing.T) {
	 var expectPropertiesMetadata = map[string]string{
		 "harbor.scanner-adapter/scanner-type": "os-package-vulnerability",
		 "env.LOG_LEVEL":                       log.GetLevel().String(),
	 }
 
	 ad := NewFakeAdapter()
	 md := ad.fakeGetMetadata()
 
	 if !cmp.Equal(md.Properties, expectPropertiesMetadata) {
		 t.Errorf("PropertiesMetadata did not match expected value %s", expectPropertiesMetadata)
	 }
 }
