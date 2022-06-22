// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"fmt"
)

type Pagination struct {
	Enable   bool  `json:"enable"` // TODO: I see this when product list request returns. Maybe a bug?
	Enabled  bool  `json:"enabled"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"pagesize"`
}

func (p *Pagination) QueryString() string {
	return fmt.Sprintf("pagination={%%22page%%22:%d,%%22pageSize%%22:%d}", p.Page, p.PageSize)
}
