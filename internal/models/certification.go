// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type Certification struct {
	ID          string `json:"certificationid"`
	DisplayName string `json:"displayname"`
	Logo        string `json:"logo"`
	URL         string `json:"url"`
	Description string `json:"description"`
	IsEnabled   bool   `json:"isenabled"`
}
