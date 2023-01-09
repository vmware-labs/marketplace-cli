// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

const (
	MetaFileTypeInvalid = "INVALID_FILE_TYPE"
	MetaFileTypeCLI     = "CLI"
	MetaFileTypeConfig  = "CONFIG"
	MetaFileTypeMisc    = "MISC"
)

type MetaFileObject struct {
	FileID          string `json:"fileid"`
	FileName        string `json:"filename"`
	TempURL         string `json:"tempurl"`
	URL             string `json:"url"`
	IsFileBackedUp  bool   `json:"isfilebackedup"`
	ProcessingError string `json:"processingerror"`
	HashDigest      string `json:"hashdigest"`
	HashAlgorithm   string `json:"hashalgo"`
	Size            int64  `json:"size"`
	UploadedBy      string `json:"uploadedby"`
	UploadedOn      int32  `json:"uploadedon"`
	DownloadCount   int64  `json:"downloadcount"`
}

type MetaFile struct {
	ID         string            `json:"metafileid"`
	GroupId    string            `json:"groupid"`
	GroupName  string            `json:"groupname"`
	FileType   string            `json:"filetype"`
	Version    string            `json:"version"`    // Note: This is the version of this particular file...
	AppVersion string            `json:"appversion"` // Note: and this is the associated Marketplace product version
	Status     string            `json:"status"`
	Objects    []*MetaFileObject `json:"metafileobjectsList"`
	CreatedBy  string            `json:"createdby"`
	CreatedOn  int32             `json:"createdon"`
}

func (product *Product) GetMetaFilesForVersion(version string) []*MetaFile {
	var metafiles []*MetaFile
	versionObj := product.GetVersion(version)

	if versionObj != nil {
		for _, metafile := range product.MetaFiles {
			if metafile.AppVersion == versionObj.Number {
				metafiles = append(metafiles, metafile)
			}
		}
	}

	return metafiles
}
