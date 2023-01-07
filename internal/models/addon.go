// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type AddOnFile struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	ImageType        string `json:"imagetype"`
	Status           string `json:"status"`
	DeploymentStatus string `json:"deploymentstatus"`
	UploadedOn       int32  `json:"uploadedon"`
	UploadedBy       string `json:"uploadedby"`
	UpdatedOn        int32  `json:"updatedon"`
	UpdatedBy        string `json:"updatedby"`
	FileID           string `json:"fileid"`
	AppVersion       string `json:"appversion"`
	HashDigest       string `json:"hashdigest"`
	HashAlgorithm    string `json:"hashalgo"`
	DownloadCount    int64  `json:"downloadcount"`
	IsRedirectURL    bool   `json:"isredirecturl"`
	IsThirdPartyURL  bool   `json:"isthirdpartyurl"`
	Size             int64  `json:"size"`
}

func (product *Product) GetAddonFilesForVersion(version string) []*AddOnFile {
	var files []*AddOnFile
	versionObj := product.GetVersion(version)

	if versionObj != nil {
		for _, addonFile := range product.AddOnFiles {
			if addonFile.AppVersion == versionObj.Number {
				files = append(files, addonFile)
			}
		}
	}
	return files
}
