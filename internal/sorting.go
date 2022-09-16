// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"encoding/json"
	"strings"
)

const (
	SortKeyDisplayName      = "displayName"
	SortKeyCreationDate     = "createdOn"
	SortKeyUpdateDate       = "updatedOn"
	SortDirectionAscending  = "ASC"
	SortDirectionDescending = "DESC"
)

type Sorting struct {
	Order     int    `json:"order,omitempty"`
	Key       string `json:"key"`
	Direction string `json:"direction"`
}

func (s *Sorting) QueryString() string {
	value, _ := json.Marshal(s)
	replacer := strings.NewReplacer(`"`, "%22")
	return "sortBy=" + replacer.Replace(string(value))
}
