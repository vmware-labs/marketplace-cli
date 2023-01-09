// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"encoding/json"
	"strings"
)

type ListProductFilter struct {
	Text    string   `json:"search,omitempty"`
	AllOrgs bool     `json:"-"`
	OrgIds  []string `json:"Publishers,omitempty"`
}

func (f *ListProductFilter) QueryString() string {
	value, _ := json.Marshal(f)
	replacer := strings.NewReplacer(`"`, "%22")
	return "filters=" + replacer.Replace(string(value))
}
