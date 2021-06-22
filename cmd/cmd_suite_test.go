// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/models"
)

func TestCmdSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

func CreateFakeProduct(id, name, slug, status string) *models.Product {
	if id == "" {
		id = uuid.New().String()
	}
	return &models.Product{
		ProductId:   id,
		Slug:        slug,
		DisplayName: name,
		Status:      status,
		Versions:    []*models.Version{},
		EncryptionDetails: &models.ProductEncryptionDetails{
			List: []string{"userAuthEncryption"},
		},
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
		Id:                    "",
		AppVersion:            version,
		DeploymentInstruction: instructions,
		DockerURLs:            []*models.DockerURLDetails{},
		Status:                "",
	}

	for _, image := range images {
		imageList.DockerURLs = append(imageList.DockerURLs, image)
	}

	product.DockerLinkVersions = append(product.DockerLinkVersions, imageList)
	return product
}
