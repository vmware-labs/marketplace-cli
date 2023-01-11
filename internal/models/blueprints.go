// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type Image struct {
	URL         string `json:"url"`
	DownloadURL string `json:"downloadurl"`
	HashType    string `json:"hashtype"`
	HashValue   string `json:"hashvalue"`
}

type File struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	DownloadURL string `json:"downloadurl"`
	HashType    string `json:"hashtype"`
	HashValue   string `json:"hashvalue"`
}

type BlueprintFile struct {
	ID                     string  `json:"id"`
	FileID                 string  `json:"fileid"`
	Title                  string  `json:"string"`
	URL                    string  `json:"url"`
	Status                 string  `json:"status"`
	Metadata               string  `json:"metadata"`
	Images                 []Image `json:"imagesList"`
	Files                  []File  `json:"filesList"`
	VRAVersion             string  `json:"vraversion"`
	DeploymentInstructions string  `json:"deploymentinstructions"`
}

type ProductBlueprintDetails struct {
	Version        string          `json:"version"`
	Instructions   string          `json:"instructions"`
	BlueprintFiles []BlueprintFile `json:"blueprintfilesList"`
	Prerequisites  []string        `json:"prerequisitesList"`
}
