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

var _ = Describe("Other", func() {
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

	Describe("AttachOtherFile", func() {
		var filePath string

		BeforeEach(func() {
			file, err := ioutil.TempFile("", "mkpcli-attachotherfile-test-file.tgz")
			Expect(err).ToNot(HaveOccurred())
			filePath = file.Name()
			uploader.UploadProductFileReturns("uploaded-file.tgz", "https://example.com/uploaded-file.tgz", err)

			httpClient.PutStub = PutProductEchoResponse
		})

		AfterEach(func() {
			Expect(os.Remove(filePath)).To(Succeed())
		})

		It("uploads and attaches the file", func() {
			product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
			test.AddVersions(product, "1.2.3")

			updatedProduct, err := marketplace.AttachOtherFile(filePath, product, &models.Version{Number: "1.2.3"})
			Expect(err).ToNot(HaveOccurred())

			By("uploading the file", func() {
				Expect(uploader.UploadProductFileCallCount()).To(Equal(1))
				uploadedFilePath := uploader.UploadProductFileArgsForCall(0)
				Expect(uploadedFilePath).To(Equal(filePath))
			})

			By("updating the product in the marketplace", func() {
				Expect(updatedProduct.AddOnFiles).To(HaveLen(1))
				uploadedFile := updatedProduct.AddOnFiles[0]
				Expect(uploadedFile.Name).To(Equal("uploaded-file.tgz"))
				Expect(uploadedFile.AppVersion).To(Equal("1.2.3"))
				Expect(uploadedFile.URL).To(Equal("https://example.com/uploaded-file.tgz"))
				Expect(uploadedFile.HashAlgorithm).To(Equal("SHA1"))
				Expect(uploadedFile.HashDigest).To(Equal("da39a3ee5e6b4b0d3255bfef95601890afd80709"))
			})
		})

		When("hashing fails", func() {
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVersions(product, "1.2.3")

				_, err := marketplace.AttachOtherFile("this/file/does/not/exist", product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to open this/file/does/not/exist: open this/file/does/not/exist: no such file or directory"))
			})
		})

		When("getting an uploader fails", func() {
			BeforeEach(func() {
				marketplace.SetUploader(nil)
				httpClient.GetReturns(nil, errors.New("get uploader failed"))
			})

			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVersions(product, "1.2.3")

				_, err := marketplace.AttachOtherFile(filePath, product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get upload credentials: get uploader failed"))
			})
		})

		When("uploading the file fails", func() {
			BeforeEach(func() {
				uploader.UploadProductFileReturns("", "", errors.New("upload product file failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVersions(product, "1.2.3")

				_, err := marketplace.AttachOtherFile(filePath, product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("upload product file failed"))
			})
		})

		When("updating the product fails", func() {
			BeforeEach(func() {
				httpClient.PutReturns(nil, errors.New("put product failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVersions(product, "1.2.3")

				_, err := marketplace.AttachOtherFile(filePath, product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the update for product \"hyperspace-database\" failed: put product failed"))
			})
		})
	})
})
