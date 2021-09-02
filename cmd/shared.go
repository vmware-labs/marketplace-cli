// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

// Variables set from CLI flags
var (
	Marketplace *pkg.Marketplace

	OutputFormat   string
	Output         output.Format
	ProductSlug    string
	ProductVersion string

	UploadCredentials = aws.Credentials{}

	ImageRepository string
	ImageTag        string
	ImageTagType    string

	ChartName           string
	ChartVersion        string
	ChartRepositoryName string
	ChartRepositoryURL  string
	ChartURL            string

	DeploymentInstructions string
)
