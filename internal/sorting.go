// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"fmt"
	"net/url"
)

const (
	SortKeyDisplayName      = "displayName"
	SortKeyCreationDate     = "createdOn"
	SortKeyUpdateDate       = "updatedOn"
	SortDirectionAscending  = "ASC"
	SortDirectionDescending = "DESC"
)

type Sorting struct {
	Order     int    `json:"order"`
	Key       string `json:"key"`
	Direction string `json:"direction"`
}

func (s *Sorting) QueryString() string {
	return fmt.Sprintf("sortBy={%%22key%%22:%%22%s%%22,%%22direction%%22:%%22%s%%22}", s.Key, s.Direction)
}

func (s *Sorting) Apply(input *url.URL) *url.URL {
	values := input.Query()
	delete(values, "sorting")

	output := *input
	output.RawQuery = values.Encode()
	if len(values) == 0 {
		output.RawQuery = s.QueryString()
	} else {
		output.RawQuery += "&" + s.QueryString()
	}

	return &output
}
