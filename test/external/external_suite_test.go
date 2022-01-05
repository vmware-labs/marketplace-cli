// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "External test suite")
}
