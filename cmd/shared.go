// Copyright 2022 VMware, Inc.
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
	Marketplace pkg.MarketplaceInterface
	Output      output.Format
	typeMapping = map[string]string{
		"addon":    pkg.AssetTypeAddon,
		"chart":    pkg.AssetTypeChart,
		"image":    pkg.AssetTypeContainerImage,
		"metafile": pkg.AssetTypeMetaFile,
		"vm":       pkg.AssetTypeVM,
	}
)

func assetTypesList() []string {
	var assetTypes []string
	for assetType := range typeMapping {
		assetTypes = append(assetTypes, assetType)
	}
	sort.Strings(assetTypes)
	return assetTypes
}

func ValidateAssetTypeFilter(cmd *cobra.Command, args []string) error {
	if ListAssetsByType == "" {
		return nil
	}
	if typeMapping[ListAssetsByType] != "" {
		return nil
	}
	return fmt.Errorf("Unknown asset type: %s\nPlease use one of %s", ListAssetsByType, strings.Join(assetTypesList(), ", "))
}
