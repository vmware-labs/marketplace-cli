// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func (m *Marketplace) AttachLocalContainerImage(imageFile, image, tag, tagType, instructions string, product *models.Product, version *models.Version) (*models.Product, error) {
	if product.HasContainerImage(version.Number, image, tag) {
		return nil, fmt.Errorf("%s %s already has the image %s:%s", product.Slug, version.Number, image, tag)
	}

	uploader, err := m.GetUploader(product.PublisherDetails.OrgId)
	if err != nil {
		return nil, err
	}
	_, fileUrl, err := uploader.UploadProductFile(imageFile)
	if err != nil {
		return nil, err
	}

	product.PrepForUpdate()
	product.DockerLinkVersions = append(product.DockerLinkVersions, &models.DockerVersionList{
		AppVersion: version.Number,
		DockerURLs: []*models.DockerURLDetails{
			{
				Url: image,
				ImageTags: []*models.DockerImageTag{
					{
						Tag:               tag,
						Type:              tagType,
						MarketplaceS3Link: fileUrl,
					},
				},
				DeploymentInstruction: instructions,
				DockerType:            models.DockerTypeUpload,
			},
		},
	})

	return m.PutProduct(product, version.IsNewVersion)
}

func (m *Marketplace) AttachPublicContainerImage(image, tag, tagType, instructions string, product *models.Product, version *models.Version) (*models.Product, error) {
	if product.HasContainerImage(version.Number, image, tag) {
		return nil, fmt.Errorf("%s %s already has the image %s:%s", product.Slug, version.Number, image, tag)
	}

	product.PrepForUpdate()
	product.DockerLinkVersions = append(product.DockerLinkVersions, &models.DockerVersionList{
		AppVersion: version.Number,
		DockerURLs: []*models.DockerURLDetails{
			{
				Url: image,
				ImageTags: []*models.DockerImageTag{
					{
						Tag:  tag,
						Type: tagType,
					},
				},
				DeploymentInstruction: instructions,
				DockerType:            models.DockerTypeRegistry,
			},
		},
	})

	return m.PutProduct(product, version.IsNewVersion)
}
