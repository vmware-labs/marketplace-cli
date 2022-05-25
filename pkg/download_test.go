// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Download", func() {
	var (
		filename         string
		httpClient       *pkgfakes.FakeHTTPClient
		marketplace      *pkg.Marketplace
		output           *Buffer
		progressBar      *internalfakes.FakeProgressBar
		progressBarMaker *internalfakes.FakeProgressBarMaker
	)

	BeforeEach(func() {
		filename = ""
		httpClient = &pkgfakes.FakeHTTPClient{}
		output = NewBuffer()
		marketplace = &pkg.Marketplace{
			Host:   "marketplace.example.com",
			Client: httpClient,
			Output: output,
		}

		progressBar = &internalfakes.FakeProgressBar{}
		progressBar.WrapWriterStub = func(source io.Writer) io.Writer { return source }
		progressBarMaker = &internalfakes.FakeProgressBarMaker{}
		progressBarMaker.Returns(progressBar)
		internal.MakeProgressBar = progressBarMaker.Spy

		httpClient.DoReturnsOnCall(0, MakeJSONResponse(&pkg.DownloadResponse{
			Response: &pkg.DownloadResponseBody{
				PreSignedURL: "https://example.com/download/file.txt",
			},
		}), nil)
		httpClient.DoReturnsOnCall(1, MakeStringResponse("file contents!"), nil)
	})

	AfterEach(func() {
		if filename != "" {
			Expect(os.Remove(filename)).To(Succeed())
		}
	})

	It("downloads a file", func() {
		filename = "destination-file.txt"
		requestPayload := &pkg.DownloadRequestPayload{
			ProductId:  "my-product-id",
			AppVersion: "1.2.3",
		}
		err := marketplace.Download(filename, requestPayload)
		Expect(err).ToNot(HaveOccurred())

		Expect(httpClient.DoCallCount()).To(Equal(2))
		By("requesting the download link", func() {
			request := httpClient.DoArgsForCall(0)
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.String()).To(Equal("https://marketplace.example.com/api/v1/products/my-product-id/download"))
			var payload *pkg.DownloadRequestPayload
			bodyBytes, err := ioutil.ReadAll(request.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(json.Unmarshal(bodyBytes, &payload)).To(Succeed())
			Expect(payload.ProductId).To(Equal("my-product-id"))
			Expect(payload.AppVersion).To(Equal("1.2.3"))
		})

		By("sending the download request", func() {
			request := httpClient.DoArgsForCall(1)
			Expect(request.Method).To(Equal("GET"))
			Expect(request.URL.String()).To(Equal("https://example.com/download/file.txt"))
		})

		By("writing to a progress bar", func() {
			Expect(progressBarMaker.CallCount()).To(Equal(1))
			description, size, progressBarOutput := progressBarMaker.ArgsForCall(0)
			Expect(description).To(Equal("Downloading destination-file.txt"))
			Expect(size).To(Equal(int64(14)))
			Expect(progressBarOutput).To(Equal(output))
			Expect(progressBar.WrapWriterCallCount()).To(Equal(1))
		})

		By("copying the asset to the filename", func() {
			content, err := ioutil.ReadFile(filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal("file contents!"))
		})
	})

	When("requesting the download link fails", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, nil, errors.New("download link request failed"))
		})
		It("returns an error", func() {
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download("", requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to get download link: marketplace request failed: download link request failed"))
		})
	})

	When("the response is not 200 OK", func() {
		BeforeEach(func() {
			response := MakeStringResponse("download link request failed")
			response.Status = http.StatusText(http.StatusTeapot)
			response.StatusCode = http.StatusTeapot
			httpClient.DoReturnsOnCall(0, response, nil)
		})
		It("returns an error", func() {
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download("", requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to fetch download link: I'm a teapot\ndownload link request failed"))
		})

		Context("and no response body is sent", func() {
			BeforeEach(func() {
				response := &http.Response{
					Body:       ioutil.NopCloser(&test.FailingReadWriter{Message: "failed to read"}),
					Status:     http.StatusText(http.StatusTeapot),
					StatusCode: http.StatusTeapot,
				}
				httpClient.DoReturnsOnCall(0, response, nil)
			})

			It("returns an error", func() {
				requestPayload := &pkg.DownloadRequestPayload{
					ProductId:  "my-product-id",
					AppVersion: "1.2.3",
				}
				err := marketplace.Download("", requestPayload)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to fetch download link: I'm a teapot"))
			})
		})
	})

	When("the response is not a valid download payload", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, MakeStringResponse("this is not a good payload"), nil)
		})

		It("returns an error", func() {
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download("", requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to parse response: invalid character 'h' in literal true (expecting 'r')"))
		})
	})

	When("creating the target file fails", func() {
		It("returns an error", func() {
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download("/this/path/does/not/exist", requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to create file for download: open /this/path/does/not/exist: no such file or directory"))
		})
	})

	When("making the download request fails", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, MakeJSONResponse(&pkg.DownloadResponse{
				Response: &pkg.DownloadResponseBody{
					PreSignedURL: ": : this is a bad url",
				},
			}), nil)
		})
		It("returns an error", func() {
			filename = "destination-file.txt"
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download(filename, requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to create download file request: parse \": : this is a bad url\": missing protocol scheme"))
		})
	})

	When("downloading the asset fails", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(1, nil, errors.New("download failed"))
		})
		It("returns an error", func() {
			filename = "destination-file.txt"
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download(filename, requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to download file: download failed"))
		})
	})

	When("writing the asset fails", func() {
		BeforeEach(func() {
			progressBar.WrapWriterReturns(&test.FailingReadWriter{Message: "writing failed"})
		})
		It("returns an error", func() {
			filename = "destination-file.txt"
			requestPayload := &pkg.DownloadRequestPayload{
				ProductId:  "my-product-id",
				AppVersion: "1.2.3",
			}
			err := marketplace.Download(filename, requestPayload)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to download file to disk: writing failed"))
		})
	})
})
