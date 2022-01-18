// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type AddOnFiles struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	ImageType        string `json:"imagetype"`
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
