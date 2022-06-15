// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	DownloadProductSlug    string
	DownloadProductVersion string
	DownloadFilter         string
	DownloadFilename       string
	DownloadAcceptEULA     bool
)

func init() {
	rootCmd.AddCommand(DownloadCmd)

	DownloadCmd.Flags().StringVarP(&DownloadProductSlug, "product", "p", "", "Product slug (required)")
	_ = DownloadCmd.MarkFlagRequired("product")
	DownloadCmd.Flags().StringVarP(&DownloadProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	DownloadCmd.Flags().StringVar(&DownloadFilter, "filter", "", "Filter assets by display name")
	DownloadCmd.Flags().StringVarP(&AssetType, "type", "t", "", "Filter assets by type (one of "+strings.Join(assetTypesList(), ", ")+")")
	DownloadCmd.Flags().StringVarP(&DownloadFilename, "filename", "f", "", "Output file name")
	DownloadCmd.Flags().BoolVar(&DownloadAcceptEULA, "accept-eula", false, "Accept the product EULA")
}

func filterAssets(filter string, assets []*pkg.Asset) []*pkg.Asset {
	var filteredAssets []*pkg.Asset
	for _, asset := range assets {
		if strings.Contains(asset.DisplayName, filter) {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return filteredAssets
}

var DownloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download an asset from a product",
	Long:    "Download an asset attached to a product in the VMware Marketplace",
	Example: fmt.Sprintf("%s download -p hyperspace-database-chart1 -v 1.2.3", AppName),
	Args:    cobra.NoArgs,
	PreRunE: RunSerially(ValidateAssetTypeFilter, GetRefreshToken),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(DownloadProductSlug, DownloadProductVersion)
		if err != nil {
			return err
		}

		assetType := ""
		if AssetType != "" {
			assetType = assetTypeMapping[AssetType] + " "
		}
		var asset *pkg.Asset
		assets := pkg.GetAssetsByType(assetTypeMapping[AssetType], product, version.Number)
		if len(assets) == 0 {
			return fmt.Errorf("product %s %s does not have any downloadable %sassets", product.Slug, version.Number, assetType)
		}

		if DownloadFilter == "" {
			asset = assets[0]
			if len(assets) > 1 {
				_ = Output.RenderAssets(assets)
				return fmt.Errorf("product %s %s has multiple downloadable %sassets, please use the --filter parameter", product.Slug, version.Number, assetType)
			}
		} else {
			filterAssets := filterAssets(DownloadFilter, assets)
			if len(filterAssets) == 0 {
				return fmt.Errorf("product %s %s does not have any downloadable %sassets that match the filter \"%s\", please adjust the --filter parameter", product.Slug, version.Number, assetType, DownloadFilter)
			}

			asset = filterAssets[0]
			if len(filterAssets) > 1 {
				_ = Output.RenderAssets(filterAssets)
				return fmt.Errorf("product %s %s has multiple downloadable %sassets that match the filter \"%s\", please adjust the --filter parameter", product.Slug, version.Number, assetType, DownloadFilter)
			}
		}

		filename := asset.Filename
		if DownloadFilename != "" {
			filename = DownloadFilename
		}

		if !DownloadAcceptEULA && !product.EulaDetails.Signed {
			cmd.PrintErrln("The EULA must be accepted before downloading")
			if product.EulaDetails.Text != "" {
				cmd.PrintErrf("EULA: %s\n\n", product.EulaDetails.Text)
			} else if product.EulaDetails.Url != "" {
				cmd.PrintErrf("EULA: %s\n\n", product.EulaDetails.Url)
			}
			return fmt.Errorf("please review the EULA and re-run with --accept-eula")
		}

		asset.DownloadRequestPayload.EulaAccepted = DownloadAcceptEULA
		return Marketplace.Download(filename, asset.DownloadRequestPayload)
	},
}
