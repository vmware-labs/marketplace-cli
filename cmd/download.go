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
)

func init() {
	rootCmd.AddCommand(DownloadCmd)

	DownloadCmd.Flags().StringVarP(&DownloadProductSlug, "product", "p", "", "Product slug (required)")
	_ = DownloadCmd.MarkFlagRequired("product")
	DownloadCmd.Flags().StringVarP(&DownloadProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	DownloadCmd.Flags().StringVar(&DownloadFilter, "filter", "", "filter to select from multiple files")
	DownloadCmd.Flags().StringVarP(&DownloadFilename, "filename", "f", "", "output file name")

}

var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an asset from a product",
	Long:  "Download an asset attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(DownloadProductSlug, DownloadProductVersion)
		if err != nil {
			return err
		}

		var asset *pkg.Asset
		assets := pkg.GetAssets(product, version.Number)
		if len(assets) == 0 {
			return fmt.Errorf("product %s %s does not have any downloadable assets", product.Slug, version.Number)
		} else if len(assets) > 1 {
			if DownloadFilter == "" {
				return fmt.Errorf("product %s %s has multiple downloadable assets, please use the --filter parameter", product.Slug, version.Number)
			} else {
				for _, potentialAsset := range assets {
					if strings.Contains(potentialAsset.Filename, DownloadFilter) {
						if asset == nil {
							asset = potentialAsset
						} else {
							return fmt.Errorf("product %s %s has multiple downloadable assets that match the filter \"%s\", please adjust the --filter parameter", product.Slug, version.Number, DownloadFilter)
						}
					}
				}
				if asset == nil {
					return fmt.Errorf("product %s %s does not have any downloadable assets that match the filter \"%s\", please adjust the --filter parameter", product.Slug, version.Number, DownloadFilter)
				}
			}
		} else {
			asset = assets[0]
		}

		filename := asset.Filename
		if DownloadFilename != "" {
			filename = DownloadFilename
		}
		return Marketplace.Download(product.ProductId, filename, asset.DownloadRequestPayload)
	},
}
