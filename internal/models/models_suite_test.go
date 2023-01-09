// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models test suite")
}
