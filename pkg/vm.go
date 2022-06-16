// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"
	"time"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func makeUniqueFileID() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("fileuploader%d.url", now)
}

func (m *Marketplace) UploadVM(vmFile string, product *models.Product, version *models.Version) (*models.Product, error) {
	hashString, err := Hash(vmFile, models.HashAlgoSHA1)
	if err != nil {
		return nil, err
	}

	uploader, err := m.GetUploader(product.PublisherDetails.OrgId)
	if err != nil {
		return nil, err
	}
	filename, fileUrl, err := uploader.UploadProductFile(vmFile)
	if err != nil {
		return nil, err
	}

	product.PrepForUpdate()
	product.ProductDeploymentFiles = []*models.ProductDeploymentFile{
		{
			Name:          filename,
			AppVersion:    version.Number,
			Url:           fileUrl,
			HashAlgo:      models.HashAlgoSHA1,
			HashDigest:    hashString,
			IsRedirectUrl: false,
			UniqueFileID:  makeUniqueFileID(),
			VersionList:   []string{},
		},
	}

	return m.PutProduct(product, version.IsNewVersion)
}
