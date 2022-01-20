// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package test

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func CreateFakeProduct(id, name, slug, status string) *models.Product {
	if id == "" {
		id = uuid.New().String()
	}
	return &models.Product{
		ProductId:   id,
		Slug:        slug,
		DisplayName: name,
		Status:      status,
		AllVersions: []*models.Version{},
		EncryptionDetails: &models.ProductEncryptionDetails{
			List: []string{"userAuthEncryption"},
		},
		PublisherDetails: &models.Publisher{
			UserId:         "test-user",
			OrgId:          uuid.New().String(),
			OrgName:        "my-org",
			OrgDisplayName: "my-org",
		},
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
				FileName:       name,
				Size:           123,
				DownloadCount:  25,
				IsFileBackedUp: true,
			},
		},
	}
}

func AddVerions(product *models.Product, versions ...string) *models.Product {
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
