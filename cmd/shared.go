// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	Client      pkg.HTTPClient
	Marketplace pkg.MarketplaceInterface
	Output      output.Format

	AssetType        string
	assetTypeMapping = map[string]string{
		"other":    pkg.AssetTypeOther,
		"chart":    pkg.AssetTypeChart,
		"image":    pkg.AssetTypeContainerImage,
		"metafile": pkg.AssetTypeMetaFile,
		"vm":       pkg.AssetTypeVM,
	}
	MetaFileType        string
	metaFileTypeMapping = map[string]string{
		"cli":    pkg.MetaFileTypeCLI,
		"config": pkg.MetaFileTypeConfig,
		"other":  pkg.MetaFileTypeOther,
	}
)

func assetTypesList() []string {
	var assetTypes []string
	for assetType := range assetTypeMapping {
		assetTypes = append(assetTypes, assetType)
	}
	sort.Strings(assetTypes)
	return assetTypes
}

func ValidateAssetTypeFilter(cmd *cobra.Command, args []string) error {
	if AssetType == "" {
		return nil
	}
	if assetTypeMapping[AssetType] != "" {
		return nil
	}
	return fmt.Errorf("Unknown asset type: %s\nPlease use one of %s", AssetType, strings.Join(assetTypesList(), ", "))
}

func metaFileTypesList() []string {
	var metaFileTypes []string
	for metaFileType := range metaFileTypeMapping {
		metaFileTypes = append(metaFileTypes, metaFileType)
	}
	sort.Strings(metaFileTypes)
	return metaFileTypes
}

func ValidateMetaFileType(cmd *cobra.Command, args []string) error {
	if MetaFileType == "" {
		return nil
	}
	if metaFileTypeMapping[MetaFileType] != "" {
		return nil
	}
	return fmt.Errorf("Unknown meta file type: %s\nPlease use one of %s", MetaFileType, strings.Join(metaFileTypesList(), ", "))
}
