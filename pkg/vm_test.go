// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"errors"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("VM", func() {
	var (
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
		uploader    *internalfakes.FakeUploader
	)

	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
			Host:   "marketplace.vmware.example",
		}
		uploader = &internalfakes.FakeUploader{}
		marketplace.SetUploader(uploader)
	})

	Describe("UploadVM", func() {
		var vmFilePath string

		BeforeEach(func() {
			vmFile, err := ioutil.TempFile("", "mkpcli-uploadvm-test-vm.iso")
			Expect(err).ToNot(HaveOccurred())
			vmFilePath = vmFile.Name()
			uploader.UploadProductFileReturns("uploaded-file.iso", "https://example.com/uploaded-file.iso", err)

			httpClient.DoStub = PutProductEchoResponse
		})

		AfterEach(func() {
			Expect(os.Remove(vmFilePath)).To(Succeed())
		})

		It("uploads and attaches the vm file", func() {
			product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
			test.AddVerions(product, "1.2.3")

			updatedProduct, err := marketplace.UploadVM(vmFilePath, product, &models.Version{Number: "1.2.3"})
			Expect(err).ToNot(HaveOccurred())

			By("uploading the file", func() {
				Expect(uploader.UploadProductFileCallCount()).To(Equal(1))
				uploadedFilePath := uploader.UploadProductFileArgsForCall(0)
				Expect(uploadedFilePath).To(Equal(vmFilePath))
			})

			By("updating the product in the marketplace", func() {
				Expect(updatedProduct.ProductDeploymentFiles).To(HaveLen(1))
				uploadedFile := updatedProduct.ProductDeploymentFiles[0]
				Expect(uploadedFile.Name).To(Equal("uploaded-file.iso"))
				Expect(uploadedFile.AppVersion).To(Equal("1.2.3"))
				Expect(uploadedFile.Url).To(Equal("https://example.com/uploaded-file.iso"))
				Expect(uploadedFile.HashAlgo).To(Equal("SHA1"))
				Expect(uploadedFile.HashDigest).To(Equal("da39a3ee5e6b4b0d3255bfef95601890afd80709"))
				Expect(uploadedFile.IsRedirectUrl).To(BeFalse())
				Expect(uploadedFile.UniqueFileID).To(MatchRegexp("fileuploader[0-9]+.url"))
			})
		})

		When("hashing fails", func() {
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1.2.3")

				_, err := marketplace.UploadVM("this/file/does/not/exist", product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to open this/file/does/not/exist: open this/file/does/not/exist: no such file or directory"))
			})
		})

		When("getting an uploader fails", func() {
			BeforeEach(func() {
				marketplace.SetUploader(nil)
				httpClient.DoReturns(nil, errors.New("get uploader failed"))
			})

			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1.2.3")

				_, err := marketplace.UploadVM(vmFilePath, product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get upload credentials: marketplace request failed: get uploader failed"))
			})
		})

		When("uploading the VM image fails", func() {
			BeforeEach(func() {
				uploader.UploadProductFileReturns("", "", errors.New("upload product file failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1.2.3")

				_, err := marketplace.UploadVM(vmFilePath, product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("upload product file failed"))
			})
		})

		When("updating the product fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("put product failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1.2.3")

				_, err := marketplace.UploadVM(vmFilePath, product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the update for product \"hyperspace-database\" failed: marketplace request failed: put product failed"))
			})
		})
	})
})
