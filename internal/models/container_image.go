// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

const (
	ImageTagTypeFixed    = "FIXED"
	ImageTagTypeFloating = "FLOATING"
)

type DockerImageTag struct {
	ID                             string `json:"id,omitempty"`
	Tag                            string `json:"tag"`
	Type                           string `json:"type"`
	IsUpdatedInMarketplaceRegistry bool   `json:"isupdatedinmarketplaceregistry"`
	MarketplaceS3Link              string `json:"marketplaces3link,omitempty"`
	AppCheckReportLink             string `json:"appcheckreportlink,omitempty"`
	AppCheckSummaryPdfLink         string `json:"appchecksummarypdflink,omitempty"`
	S3TarBackupUrl                 string `json:"s3tarbackupurl,omitempty"`
	ProcessingError                string `json:"processingerror,omitempty"`
	DownloadCount                  int64  `json:"downloadcount,omitempty"`
	DownloadURL                    string `json:"downloadurl,omitempty"`
	HashAlgo                       string `json:"hashalgo,omitempty"`
	HashDigest                     string `json:"hashdigest,omitempty"`
	Size                           int64  `json:"size,omitempty"`
}

type DockerURLDetails struct {
	Key                   string            `json:"key,omitempty"`
	Url                   string            `json:"url,omitempty"`
	MarketplaceUpdatedUrl string            `json:"marketplaceupdatedurl,omitempty"`
	ImageTags             []*DockerImageTag `json:"imagetagsList"`
	ImageTagsAsJson       string            `json:"imagetagsasjson"`
	DockerType            string            `json:"dockertype,omitempty"`
	ID                    string            `json:"id,omitempty"`
	DeploymentInstruction string            `json:"deploymentinstruction"`
	Name                  string            `json:"name"`
	IsMultiArch           bool              `json:"ismultiarch"`
}

const (
	DockerTypeRegistry = "registry"
	DockerTypeUpload   = "upload"
)

func (d *DockerURLDetails) GetTag(tagName string) *DockerImageTag {
	for _, tag := range d.ImageTags {
		if tag.Tag == tagName {
			return tag
		}
	}
	return nil
}

func (d *DockerURLDetails) HasTag(tagName string) bool {
	return d.GetTag(tagName) != nil
}

type DockerVersionList struct {
	ID                    string              `json:"id,omitempty"`
	AppVersion            string              `json:"appversion"`
	DeploymentInstruction string              `json:"deploymentinstruction"`
	DockerURLs            []*DockerURLDetails `json:"dockerurlsList"`
	Status                string              `json:"status,omitempty"`
	ImageTags             []*DockerImageTag   `json:"imagetagsList"`
}

func (product *Product) HasContainerImage(version, imageURL, tag string) bool {
	versionObj := product.GetVersion(version)

	if versionObj != nil {
		for _, dockerVersionLink := range product.DockerLinkVersions {
			if dockerVersionLink.AppVersion == versionObj.Number {
				for _, dockerUrl := range dockerVersionLink.DockerURLs {
					if dockerUrl.Url == imageURL {
						for _, imageTag := range dockerUrl.ImageTags {
							if imageTag.Tag == tag {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (product *Product) GetContainerImagesForVersion(version string) []*DockerVersionList {
	var images []*DockerVersionList
	versionObj := product.GetVersion(version)

	if versionObj != nil {
		for _, dockerVersionLink := range product.DockerLinkVersions {
			if dockerVersionLink.AppVersion == versionObj.Number {
				images = append(images, dockerVersionLink)
			}
		}
	}
	return images
}
