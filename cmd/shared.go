// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

// Variables set from CLI flags
var (
	OutputFormat   string
	ProductSlug    string
	ProductVersion string

	ImageRepository string
	ImageTag        string
	ImageTagType    string

	ChartName           string
	ChartVersion        string
	ChartRepositoryName string
	ChartRepositoryURL  string
	ChartURL            string

	DeploymentInstructions string

	SearchTerm string
)
