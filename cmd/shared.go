// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	Marketplace pkg.MarketplaceInterface
	Output      output.Format
)
