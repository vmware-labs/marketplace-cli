// Copyright 2021 VMware, Inc.
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
		AppVersion: version,
		FileID:     uuid.New().String(),
		Status:     "STORED",
		ItemJson:   string(detailString),
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
			Tag:  tag,
			Type: tagType,
		})
	}

	return &models.DockerURLDetails{
		ID:        uuid.New().String(),
		Url:       url,
		ImageTags: tagList,
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
		ID:                    "",
		AppVersion:            version,
		DeploymentInstruction: instructions,
		DockerURLs:            []*models.DockerURLDetails{},
		Status:                "",
	}

	imageList.DockerURLs = append(imageList.DockerURLs, images...)

	product.DockerLinkVersions = append(product.DockerLinkVersions, imageList)
	return product
}
