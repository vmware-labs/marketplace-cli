// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

const (
	MetaFileTypeCLI    = "CLI"
	MetaFileTypeConfig = "CONFIG"
	MetaFileTypeOther  = "MISC"
)

func (m *Marketplace) AttachMetaFile(metafile, metafileType, metafileVersion string, product *models.Product, version *models.Version) (*models.Product, error) {
	hashString, err := Hash(metafile, models.HashAlgoSHA1)
	if err != nil {
		return nil, err
	}

	uploader, err := m.GetUploader(product.PublisherDetails.OrgId)
	if err != nil {
		return nil, err
	}
	filename, fileUrl, err := uploader.UploadMetaFile(metafile)
	if err != nil {
		return nil, err
	}

	product.PrepForUpdate()
	product.MetaFiles = append(product.MetaFiles, &models.MetaFile{
		FileType:   metafileType,
		Version:    metafileVersion,
		AppVersion: version.Number,
		Objects: []*models.MetaFileObject{
			{
				FileName:      filename,
				TempURL:       fileUrl,
				HashDigest:    hashString,
				HashAlgorithm: models.HashAlgoSHA1,
			},
		},
	})

	return m.PutProduct(product, version.IsNewVersion)
}
