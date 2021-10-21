/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
package imagescanning

import (
	"encoding/base64"
	"fmt"
	"testing"
)

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
