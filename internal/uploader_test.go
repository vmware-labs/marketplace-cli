// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	"errors"
	"io"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
)

var _ = Describe("Uploader", func() {
	var (
		client           *internalfakes.FakeS3Client
		file             *os.File
		progressBar      *internalfakes.FakeProgressBar
		progressBarMaker *internalfakes.FakeProgressBarMaker
	)

	BeforeEach(func() {
		var err error
		file, err = os.CreateTemp("", "mkpcli-test-uploader-file-*.txt")
		Expect(err).ToNot(HaveOccurred())
		_, err = file.WriteString("file contents")
		Expect(err).ToNot(HaveOccurred())

		client = &internalfakes.FakeS3Client{}
		client.PutObjectReturns(nil, nil)

		progressBar = &internalfakes.FakeProgressBar{}
		progressBar.WrapReaderStub = func(source io.Reader) io.Reader { return source }
		progressBarMaker = &internalfakes.FakeProgressBarMaker{}
		progressBarMaker.Returns(progressBar)
		internal.MakeProgressBar = progressBarMaker.Spy
	})

	AfterEach(func() {
		Expect(os.Remove(file.Name())).To(Succeed())
	})

	Describe("UploadMediaFile", func() {
		It("properly uploads a media file", func() {
			output := NewBuffer()
			uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client, output)
			filename, fileUrl, err := uploader.UploadMediaFile(file.Name())
			Expect(err).ToNot(HaveOccurred())

			By("sending the object to S3", func() {
				Expect(client.PutObjectCallCount()).To(Equal(1))
				_, putArg, options := client.PutObjectArgsForCall(0)
				Expect(*putArg.Bucket).To(Equal("my-bucket"))
				Expect(*putArg.Key).To(MatchRegexp("^my-org/media-files/[0-9]+/mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(putArg.ContentLength).To(Equal(int64(len("file contents"))))
				Expect(options).To(BeEmpty())
			})

			By("writing to a progress bar", func() {
				Expect(progressBarMaker.CallCount()).To(Equal(1))
				description, size, progressBarOutput := progressBarMaker.ArgsForCall(0)
				Expect(description).To(MatchRegexp("^Uploading mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(size).To(Equal(int64(13)))
				Expect(progressBarOutput).To(Equal(output))
				Expect(progressBar.WrapReaderCallCount()).To(Equal(1))
			})

			By("returning the uploaded filename and url", func() {
				Expect(filename).To(MatchRegexp("^mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(fileUrl).To(MatchRegexp("^https://my-bucket.s3.my-region.amazonaws.com/my-org/media-files/[0-9]+/mkpcli-test-uploader-file-[0-9]+.txt$"))
			})
		})

		When("the upload fails", func() {
			BeforeEach(func() {
				client.PutObjectReturns(nil, errors.New("put object failed"))
			})
			It("returns an error", func() {
				output := NewBuffer()
				uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client, output)
				_, _, err := uploader.UploadMediaFile(file.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload file: put object failed"))
			})
		})
	})

	Describe("UploadMetaFile", func() {
		It("properly uploads a meta file", func() {
			output := NewBuffer()
			uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client, output)
			filename, fileUrl, err := uploader.UploadMetaFile(file.Name())
			Expect(err).ToNot(HaveOccurred())

			By("sending the object to S3", func() {
				Expect(client.PutObjectCallCount()).To(Equal(1))
				_, putArg, options := client.PutObjectArgsForCall(0)
				Expect(*putArg.Bucket).To(Equal("my-bucket"))
				Expect(*putArg.Key).To(MatchRegexp("^my-org/meta-files/[0-9]+/mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(putArg.ContentLength).To(Equal(int64(len("file contents"))))
				Expect(options).To(BeEmpty())
			})

			By("writing to a progress bar", func() {
				Expect(progressBarMaker.CallCount()).To(Equal(1))
				description, size, progressBarOutput := progressBarMaker.ArgsForCall(0)
				Expect(description).To(MatchRegexp("^Uploading mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(size).To(Equal(int64(13)))
				Expect(progressBarOutput).To(Equal(output))
				Expect(progressBar.WrapReaderCallCount()).To(Equal(1))
			})

			By("returning the uploaded filename and url", func() {
				Expect(filename).To(MatchRegexp("^mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(fileUrl).To(MatchRegexp("^https://my-bucket.s3.my-region.amazonaws.com/my-org/meta-files/[0-9]+/mkpcli-test-uploader-file-[0-9]+.txt$"))
			})
		})

		When("the upload fails", func() {
			BeforeEach(func() {
				client.PutObjectReturns(nil, errors.New("put object failed"))
			})
			It("returns an error", func() {
				output := NewBuffer()
				uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client, output)
				_, _, err := uploader.UploadMediaFile(file.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload file: put object failed"))
			})
		})
	})

	Describe("UploadProductFile", func() {
		It("properly uploads a product file", func() {
			output := NewBuffer()
			uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client, output)
			filename, fileUrl, err := uploader.UploadProductFile(file.Name())
			Expect(err).ToNot(HaveOccurred())

			By("sending the object to S3", func() {
				Expect(client.PutObjectCallCount()).To(Equal(1))
				_, putArg, options := client.PutObjectArgsForCall(0)
				Expect(*putArg.Bucket).To(Equal("my-bucket"))
				Expect(*putArg.Key).To(MatchRegexp("^my-org/marketplace-product-files/mkpcli-test-uploader-file-[0-9]+-[0-9]+.txt$"))
				Expect(putArg.ContentLength).To(Equal(int64(len("file contents"))))
				Expect(options).To(BeEmpty())
			})

			By("writing to a progress bar", func() {
				Expect(progressBarMaker.CallCount()).To(Equal(1))
				description, size, progressBarOutput := progressBarMaker.ArgsForCall(0)
				Expect(description).To(MatchRegexp("^Uploading mkpcli-test-uploader-file-[0-9]+.txt$"))
				Expect(size).To(Equal(int64(13)))
				Expect(progressBarOutput).To(Equal(output))
				Expect(progressBar.WrapReaderCallCount()).To(Equal(1))
			})

			By("returning the uploaded filename and url", func() {
				Expect(filename).To(MatchRegexp("^mkpcli-test-uploader-file-[0-9]+-[0-9]+.txt$"))
				Expect(fileUrl).To(MatchRegexp("^https://my-bucket.s3.my-region.amazonaws.com/my-org/marketplace-product-files/mkpcli-test-uploader-file-[0-9]+-[0-9]+.txt$"))
			})
		})

		When("the upload fails", func() {
			BeforeEach(func() {
				client.PutObjectReturns(nil, errors.New("put object failed"))
			})
			It("returns an error", func() {
				output := NewBuffer()
				uploader := internal.NewS3Uploader("my-bucket", "my-region", "my-org", client, output)
				_, _, err := uploader.UploadProductFile(file.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload file: put object failed"))
			})
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
