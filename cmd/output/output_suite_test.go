// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOutputSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Output test suite")
}
