// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type ProductDeploymentFile struct {
	Id              string   `json:"id,omitempty"` // uuid
	Name            string   `json:"name,omitempty"`
	Url             string   `json:"url,omitempty"`
	ImageType       string   `json:"imagetype,omitempty"`
	Status          string   `json:"status,omitempty"`
	UploadedOn      int32    `json:"uploadedon,omitempty"`
	UploadedBy      string   `json:"uploadedby,omitempty"`
	UpdatedOn       int32    `json:"updatedon,omitempty"`
	UpdatedBy       string   `json:"updatedby,omitempty"`
	ItemJson        string   `json:"itemjson,omitempty"`
	Itemkey         string   `json:"itemkey,omitempty"`
	FileID          string   `json:"fileid,omitempty"`
	IsSubscribed    bool     `json:"issubscribed,omitempty"`
	AppVersion      string   `json:"appversion"` // Mandatory
	HashDigest      string   `json:"hashdigest"`
	IsThirdPartyUrl bool     `json:"isthirdpartyurl,omitempty"`
	ThirdPartyUrl   string   `json:"thirdpartyurl,omitempty"`
	IsRedirectUrl   bool     `json:"isredirecturl,omitempty"`
	Comment         string   `json:"comment,omitempty"`
	HashAlgo        string   `json:"hashalgo"`
	DownloadCount   int64    `json:"downloadcount,omitempty"`
	UniqueFileID    string   `json:"uniqueFileId,omitempty"`
	VersionList     []string `json:"versionList"`
	Size            int64    `json:"size,omitempty"`
}

func (product *Product) GetFilesForVersion(version string) []*ProductDeploymentFile {
	var files []*ProductDeploymentFile
	for _, file := range product.ProductDeploymentFiles {
		if file.AppVersion == product.GetVersion(version).Number {
			files = append(files, file)
		}
	}
	return files
}

func (product *Product) GetFile(fileId string) *ProductDeploymentFile {
	for _, file := range product.ProductDeploymentFiles {
		if file.FileID == fileId {
			return file
		}
	}
	return nil
}

func (product *Product) AddFile(file *ProductDeploymentFile) {
	product.ProductDeploymentFiles = append(product.ProductDeploymentFiles, file)
}
