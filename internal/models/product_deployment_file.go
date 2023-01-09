// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

import (
	"encoding/json"
)

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
	versionObj := product.GetVersion(version)

	if versionObj != nil {
		for _, file := range product.ProductDeploymentFiles {
			if file.AppVersion == versionObj.Number {
				files = append(files, file)
			}
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

func (f *ProductDeploymentFile) CalculateSize() int64 {
	if f.Size > 0 {
		return f.Size
	}

	details := &ProductItemDetails{}
	err := json.Unmarshal([]byte(f.ItemJson), details)
	if err != nil {
		return 0
	}

	var size int64 = 0
	for _, file := range details.Files {
		size += int64(file.Size)
	}
	return size
}
