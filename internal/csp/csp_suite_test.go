// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmdSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSP test suite")
}
