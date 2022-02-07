// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
)

var _ = Describe("Upload", func() {
	var (
		client *internalfakes.FakeS3Client
		file   *os.File
	)

	BeforeEach(func() {
		var err error
		file, err = os.CreateTemp("", "mkpcli-test-uploader-file-*.txt")
		Expect(err).ToNot(HaveOccurred())
		_, err = file.WriteString("file contents")
		Expect(err).ToNot(HaveOccurred())

		client = &internalfakes.FakeS3Client{}
		client.PutObjectReturns(nil, nil)
	})

	AfterEach(func() {
		err := os.Remove(file.Name())
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("UploadMediaFile", func() {
		It("properly uploads a media file", func() {
			uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client)
			filename, fileUrl, err := uploader.UploadMediaFile(file.Name())
			Expect(err).ToNot(HaveOccurred())

			By("sending the object to S3", func() {
				Expect(client.PutObjectCallCount()).To(Equal(1))
				_, putArg, options := client.PutObjectArgsForCall(0)
				Expect(*putArg.Bucket).To(Equal("my-bucket"))
				Expect(*putArg.Key).To(MatchRegexp("^my-org/media-files/[0-9]+/mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(putArg.Body.(*os.File).Name()).To(Equal(file.Name()))
				Expect(putArg.ContentLength).To(Equal(int64(len("file contents"))))
				Expect(options).To(BeEmpty())
			})
			Expect(filename).To(MatchRegexp("^mkpcli-test-uploader-file-[0-9]+.txt$"))
			Expect(fileUrl).To(MatchRegexp("^https://stg-cdn.market.csp.vmware.com/my-org/media-files/[0-9]+/mkpcli-test-uploader-file-[0-9]+.txt$"))
		})
	})

	Describe("UploadProductFile", func() {
		It("properly uploads a product file", func() {
			uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client)
			filename, fileUrl, err := uploader.UploadProductFile(file.Name())
			Expect(err).ToNot(HaveOccurred())

			By("sending the object to S3", func() {
				Expect(client.PutObjectCallCount()).To(Equal(1))
				_, putArg, options := client.PutObjectArgsForCall(0)
				Expect(*putArg.Bucket).To(Equal("my-bucket"))
				Expect(*putArg.Key).To(MatchRegexp("^my-org/marketplace-product-files/mkpcli-test-uploader-file-[0-9]+-[0-9]+.txt$"))
				Expect(putArg.Body.(*os.File).Name()).To(Equal(file.Name()))
				Expect(putArg.ContentLength).To(Equal(int64(len("file contents"))))
				Expect(options).To(BeEmpty())
			})
			Expect(filename).To(MatchRegexp("^mkpcli-test-uploader-file-[0-9]+-[0-9]+.txt$"))
			Expect(fileUrl).To(MatchRegexp("^https://my-bucket.s3.my-region.amazonaws.com/my-org/marketplace-product-files/mkpcli-test-uploader-file-[0-9]+-[0-9]+.txt$"))
		})
	})

	Describe("MakeUniqueFilename", func() {
		It("Returns a new unique filename", func() {
			Expect(internal.MakeUniqueFilename("test.binary")).To(MatchRegexp("test-[0-9]*.binary"))
			Expect(internal.MakeUniqueFilename("two.dots.test")).To(MatchRegexp("two.dots-[0-9]*.test"))
			Expect(internal.MakeUniqueFilename("no-dots")).To(MatchRegexp("no-dots-[0-9]*"))
		})
	})
})
