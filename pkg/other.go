// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func (m *Marketplace) AttachOtherFile(file string, product *models.Product, version *models.Version) (*models.Product, error) {
	hashString, err := Hash(file, models.HashAlgoSHA1)
	if err != nil {
		return nil, err
	}

	uploader, err := m.GetUploader(product.PublisherDetails.OrgId)
	if err != nil {
		return nil, err
	}
	filename, fileUrl, err := uploader.UploadProductFile(file)
	if err != nil {
		return nil, err
	}

	product.PrepForUpdate()
	product.AddOnFiles = []*models.AddOnFile{{
		Name:          filename,
		URL:           fileUrl,
		AppVersion:    version.Number,
		HashDigest:    hashString,
		HashAlgorithm: models.HashAlgoSHA1,
	}}

	return m.PutProduct(product, version.IsNewVersion)
}
