// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package features

import (
	"os/exec"
	"strings"
	"time"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	CommandSession *gexec.Session
	MkpcliPath     string
)

var _ = BeforeSuite(func() {
	var err error
	MkpcliPath, err = gexec.Build(
		"github.com/vmware-labs/marketplace-cli/v2",
		"-ldflags",
		"-X github.com/vmware-labs/marketplace-cli/v2/cmd.Version=1.2.3",
	)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func DefineCommonSteps(define Definitions) {
	define.When(`^running mkpcli (.*)$`, func(argString string) {
		command := exec.Command(MkpcliPath, strings.Split(argString, " ")...)
		var err error
		CommandSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
	})

	define.Then(`^the command exits without error$`, func() {
		Eventually(CommandSession, time.Minute).Should(gexec.Exit(0))
	})
}
