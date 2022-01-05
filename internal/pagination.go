// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"fmt"
	"net/url"
)

type Pagination struct {
	Enabled  bool  `json:"enabled"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"pagesize"`
}

func (p Pagination) QueryString() string {
	return fmt.Sprintf("pagination={%%22page%%22:%d,%%22pageSize%%22:%d}", p.Page, p.PageSize)
}

func (p Pagination) Apply(input *url.URL) *url.URL {
	values := input.Query()
	delete(values, "pagination")

	output := *input
	output.RawQuery = values.Encode()
	if len(values) == 0 {
		output.RawQuery = p.QueryString()
	} else {
		output.RawQuery += "&" + p.QueryString()
	}

	return &output
}
