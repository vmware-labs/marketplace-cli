// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lib_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLibraries(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lib Suite")
}
