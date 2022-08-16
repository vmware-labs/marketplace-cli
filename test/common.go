// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

func CreateFakeProduct(id, name, slug, solutionType string) *models.Product {
	if id == "" {
		id = uuid.New().String()
	}
	return &models.Product{
		ProductId:    id,
		Slug:         slug,
		DisplayName:  name,
		SolutionType: solutionType,
		Status:       "pending",
		AllVersions:  []*models.Version{},
		EncryptionDetails: &models.ProductEncryptionDetails{
			List: []string{"userAuthEncryption"},
		},
		PublisherDetails: &models.Publisher{
			UserId:         "test-user",
			OrgId:          uuid.New().String(),
			OrgName:        "my-org",
			OrgDisplayName: "my-org",
		},
		EulaDetails: &models.EULADetails{
			Text: "This is the EULA text",
		},
	}
}

func CreateFakeOtherFile(name, version string) *models.AddOnFile {
	return &models.AddOnFile{
		ID:            uuid.New().String(),
		Name:          name,
		URL:           "https://marketplace.example.com/product-files/" + name,
		Status:        models.DeploymentStatusActive,
		FileID:        uuid.New().String(),
		AppVersion:    version,
		DownloadCount: 18,
		Size:          140,
	}
}

func CreateFakeOVA(name, version string) *models.ProductDeploymentFile {
	details := &models.ProductItemDetails{
		Name: name,
		Files: []*models.ProductItemFile{
			{
				Name: "some-huge-file.vmdk",
				Size: 1000000,
			},
			{
				Name: "some-small-file.txt",
				Size: 100,
			},
		},
		Type: "fake.ovf",
	}
	detailString, err := json.Marshal(details)
	Expect(err).ToNot(HaveOccurred())

	return &models.ProductDeploymentFile{
		FileID:        uuid.New().String(),
		Name:          name,
		Status:        models.DeploymentStatusActive,
		ItemJson:      string(detailString),
		AppVersion:    version,
		Size:          1000100,
		DownloadCount: 20,
	}
}

func CreateFakeChart(name string) (*chart.Chart, string, string) {
	chartDir, err := os.MkdirTemp("", "mkpcli-test-chart")
	Expect(err).ToNot(HaveOccurred())

	chartFile, err := chartutil.Create(name, chartDir)
	Expect(err).ToNot(HaveOccurred())

	testChart, err := loader.Load(chartFile)
	Expect(err).ToNot(HaveOccurred())

	chartPath, err := chartutil.Save(testChart, chartDir)
	Expect(err).ToNot(HaveOccurred())

	return testChart, chartPath, chartDir
}

func CreateFakeContainerImage(url string, tags ...string) *models.DockerURLDetails {
	var tagList []*models.DockerImageTag

	for _, tag := range tags {
		tagType := "FIXED"
		if tag == "latest" {
			tagType = "FLOATING"
		}

		tagList = append(tagList, &models.DockerImageTag{
			ID:                             uuid.New().String(),
			Tag:                            tag,
			Type:                           tagType,
			Size:                           12345,
			DownloadCount:                  15,
			IsUpdatedInMarketplaceRegistry: true,
		})
	}

	return &models.DockerURLDetails{
		ID:        uuid.New().String(),
		Url:       url,
		ImageTags: tagList,
	}
}

func CreateFakeMetaFile(name, version, productVersion string) *models.MetaFile {
	return &models.MetaFile{
		ID:         uuid.New().String(),
		FileType:   models.MetaFileTypeCLI,
		Version:    version,
		AppVersion: productVersion,
		Objects: []*models.MetaFileObject{
			{
				FileID:         uuid.New().String(),
				FileName:       name,
				Size:           123,
				DownloadCount:  25,
				IsFileBackedUp: true,
			},
		},
	}
}

func AddVersions(product *models.Product, versions ...string) *models.Product {
	for _, version := range versions {
		versionObject := &models.Version{
			Number:       version,
			Details:      fmt.Sprintf("Details for %s", version),
			Status:       "PENDING",
			Instructions: fmt.Sprintf("Instructions for %s", version),
		}
		product.AllVersions = append(product.AllVersions, versionObject)

		if versionObject.Status != "PENDING" {
			product.Versions = append(product.Versions, versionObject)
		}
	}
	return product
}

func AddContainerImages(product *models.Product, version string, instructions string, images ...*models.DockerURLDetails) *models.Product {
	imageList := &models.DockerVersionList{
		ID:                    uuid.New().String(),
		AppVersion:            version,
		DeploymentInstruction: instructions,
		DockerURLs:            []*models.DockerURLDetails{},
		Status:                "",
	}

	imageList.DockerURLs = append(imageList.DockerURLs, images...)

	product.DockerLinkVersions = append(product.DockerLinkVersions, imageList)
	return product
}

type FailingReadWriter struct {
	Message string
}

func (w *FailingReadWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New(w.Message)
}

func (r *FailingReadWriter) Read(p []byte) (n int, err error) {
	return 0, errors.New(r.Message)
}

func MakeJSONResponse(body interface{}) *http.Response {
	bodyBytes, err := json.Marshal(body)
	Expect(err).ToNot(HaveOccurred())
	return MakeBytesResponse(bodyBytes)
}

func MakeBytesResponse(body []byte) *http.Response {
	return &http.Response{
		StatusCode:    http.StatusOK,
		ContentLength: int64(len(body)),
		Body:          io.NopCloser(bytes.NewReader(body)),
	}
}

func MakeStringResponse(body string) *http.Response {
	return &http.Response{
		StatusCode:    http.StatusOK,
		ContentLength: int64(len(body)),
		Body:          io.NopCloser(strings.NewReader(body)),
	}
}

func MakeFailingBodyResponse(errorMessage string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&FailingReadWriter{Message: errorMessage}),
	}
}
