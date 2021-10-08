// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmdSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd test suite")
}

func ResponseWithPayload(payload interface{}) *http.Response {
	encoded, err := json.Marshal(payload)
	Expect(err).ToNot(HaveOccurred())

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(encoded)),
	}
}
