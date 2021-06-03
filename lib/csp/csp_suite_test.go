package csp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmdSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSP Suite")
}
