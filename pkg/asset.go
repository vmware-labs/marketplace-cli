// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

type Asset struct {
	DisplayName            string `json:"displayname"`
	Filename               string `json:"filename"`
	Version                string `json:"version"`
	Type                   string `json:"type"`
	Size                   int64  `json:"size"`
	Downloadable           bool   `json:"downloadable"`
	Downloads              int64  `json:"downloads"`
	DownloadRequestPayload *DownloadRequestPayload
}

const (
	AssetTypeVM             = "VM"
	AssetTypeChart          = "Chart"
	AssetTypeContainerImage = "Container Image"
	AssetTypeMetaFile       = "MetaFile"
)

func GetAssets(product *models.Product, version string) []*Asset {
	var assets []*Asset
	if !product.HasVersion(version) {
		return assets
	}

	for _, file := range product.GetFilesForVersion(version) {
		assets = append(assets, &Asset{
			DisplayName:  file.Name,
			Filename:     file.Name,
			Version:      "",
			Type:         AssetTypeVM,
			Size:         file.Size,
			Downloads:    file.DownloadCount,
			Downloadable: file.Status != models.DeploymentStatusInactive,
			DownloadRequestPayload: &DownloadRequestPayload{
				ProductId:        product.ProductId,
				AppVersion:       version,
				DeploymentFileId: file.FileID,
			},
		})
	}

	for _, chart := range product.GetChartsForVersion(version) {
		assets = append(assets, &Asset{
			DisplayName:  chart.HelmTarUrl,
			Filename:     "chart.tgz",
			Version:      chart.Version,
			Type:         AssetTypeChart,
			Size:         chart.Size,
			Downloadable: chart.IsUpdatedInMarketplaceRegistry,
			Downloads:    chart.DownloadCount,
			DownloadRequestPayload: &DownloadRequestPayload{
				ProductId:    product.ProductId,
				AppVersion:   version,
				ChartVersion: chart.Version,
			},
		})
	}

	containerImages := product.GetContainerImagesForVersion(version)
	if len(containerImages) > 0 {
		for _, containerImage := range containerImages {
			for _, imageURL := range containerImage.DockerURLs {
				for _, tag := range imageURL.ImageTags {
					assets = append(assets, &Asset{
						DisplayName:  fmt.Sprintf("%s:%s", imageURL.Url, tag.Tag),
						Filename:     "image.tar",
						Version:      tag.Tag,
						Type:         AssetTypeContainerImage,
						Size:         tag.Size,
						Downloads:    tag.DownloadCount,
						Downloadable: tag.IsUpdatedInMarketplaceRegistry,
						DownloadRequestPayload: &DownloadRequestPayload{
							ProductId:           product.ProductId,
							AppVersion:          version,
							DockerlinkVersionID: containerImage.ID,
							DockerUrlId:         imageURL.ID,
							ImageTagId:          tag.ID,
						},
					})
				}
			}
		}
	}

	for _, metafile := range product.GetMetaFilesForVersion(version) {
		for _, object := range metafile.Objects {
			assets = append(assets, &Asset{
				DisplayName:  object.FileName,
				Filename:     object.FileName,
				Version:      metafile.Version,
				Type:         AssetTypeMetaFile,
				Size:         object.Size,
				Downloads:    object.DownloadCount,
				Downloadable: object.IsFileBackedUp, // Is this valid?
				DownloadRequestPayload: &DownloadRequestPayload{
					ProductId:        product.ProductId,
					AppVersion:       version,
					MetaFileID:       metafile.ID,
					MetaFileObjectID: object.FileID,
				},
			})
		}
	}

	return assets
}
