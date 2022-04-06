// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package test

import (
	"os"
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
	mkpcliPath     string
)

var _ = BeforeSuite(func() {
	var err error
	mkpcliPath, err = gexec.Build(
		"github.com/vmware-labs/marketplace-cli/v2",
		"-ldflags",
		"-X github.com/vmware-labs/marketplace-cli/v2/cmd.version=1.2.3",
	)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func unsetEnvVars(envVars []string, varsToUnset []string) []string {
	var filtered []string
	for _, envVar := range envVars {
		needsToBeUnset := false
		for _, varToUnset := range varsToUnset {
			if strings.HasPrefix(envVar, varToUnset) {
				needsToBeUnset = true
			}
		}
		if !needsToBeUnset {
			filtered = append(filtered, envVar)
		}
	}
	return filtered
}

func DefineCommonSteps(define Definitions, marketplaceEnvironment string) {
	var (
		envVars   []string
		unsetVars []string
	)

	define.When(`^the environment variable ([_A-Z]*) is set to (.*)$`, func(key, value string) {
		envVars = append(envVars, key+"="+value)
	})

	define.When(`^the environment variable ([_A-Z]*) is not set$`, func(key string) {
		unsetVars = append(unsetVars, key)
	})

	define.When(`^running mkpcli (.*)$`, func(argString string) {
		command := exec.Command(mkpcliPath, strings.Split(argString, " ")...)
		envVars = append(envVars, "MARKETPLACE_ENV="+marketplaceEnvironment)
		command.Env = unsetEnvVars(append(os.Environ(), envVars...), unsetVars)

		var err error
		CommandSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
	})

	define.Then(`^the command exits without error$`, func() {
		Eventually(CommandSession, 5*time.Minute).Should(gexec.Exit(0))
	})

	define.Then(`^the command exits with an error$`, func() {
		Eventually(CommandSession, 5*time.Minute).Should(gexec.Exit(1))
	})
}
