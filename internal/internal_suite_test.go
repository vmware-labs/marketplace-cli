// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal test suite")
}
