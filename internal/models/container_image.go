// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type DockerImageTag struct {
	ID                             string `json:"id,omitempty"`
	Tag                            string `json:"tag,omitempty"`
	Type                           string `json:"type,omitempty"`
	IsUpdatedInMarketplaceRegistry bool   `json:"isupdatedinmarketplaceregistry"`
	MarketplaceS3Link              string `json:"marketplaces3link"`
	AppCheckReportLink             string `json:"appcheckreportlink"`
	AppCheckSummaryPdfLink         string `json:"appchecksummarypdflink"`
	S3TarBackupUrl                 string `json:"s3tarbackupurl"`
	ProcessingError                string `json:"processingerror"`
	DownloadCount                  int64  `json:"downloadcount"`
	DownloadURL                    string `json:"downloadurl"`
	HashAlgo                       string `json:"hashalgo"`
	HashDigest                     string `json:"hashdigest"`
	Size                           int64  `json:"size,omitempty"`
}

type DockerURLDetails struct {
	ID                    string            `json:"id,omitempty"`
	Key                   string            `json:"key,omitempty"`
	Url                   string            `json:"url,omitempty"`
	MarketplaceUpdatedUrl string            `json:"marketplaceupdatedurl"`
	ImageTags             []*DockerImageTag `json:"imagetagsList"`
	ImageTagsAsJson       string            `json:"imagetagsasjson"`
	DeploymentInstruction string            `json:"deploymentinstruction"`
}

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
