// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg test suite")
}

func PutProductEchoResponse(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error) {
	Expect(contentType).To(Equal("application/json"))
	var product *models.Product
	productBytes, err := ioutil.ReadAll(content)
	Expect(err).ToNot(HaveOccurred())
	Expect(json.Unmarshal(productBytes, &product)).To(Succeed())

	body, err := json.Marshal(&pkg.GetProductResponse{
		Response: &pkg.GetProductResponsePayload{
			Message:    "",
			StatusCode: http.StatusOK,
			Data:       product,
		},
	})
	Expect(err).ToNot(HaveOccurred())

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}, nil
}
