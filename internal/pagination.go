// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"encoding/json"
	"strings"
)

type Pagination struct {
	Enable   bool  `json:"enable,omitempty"` // TODO: I see this when product list request returns. Maybe a bug?
	Enabled  bool  `json:"enabled,omitempty"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
}

func (p *Pagination) QueryString() string {
	value, _ := json.Marshal(p)
	replacer := strings.NewReplacer(`"`, "%22")
	return "pagination=" + replacer.Replace(string(value))
}
