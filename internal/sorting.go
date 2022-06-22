// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"fmt"
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
