/*
 * Copyright 2021 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */
 
package harbor

import (
	"fmt"
)

type ErrorResponse struct {
	Error Error `json:"error,omitempty"`
}

func NewErrorResponse(msg string, err error) ErrorResponse {
	return ErrorResponse{
		Error: Error{Message: fmt.Sprintf("%s: %v", msg, err)},
	}
}

type Error struct {
	Message string `json:"message,omitempty"`
}
