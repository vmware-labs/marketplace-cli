// +build enemy

package features_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnemies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enemy test suite")
}
