// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type VmwareProduct struct {
	Id                  int64    `json:"id,omitempty"` // will be deprecated
	ShortName           string   `json:"shortname"`
	DisplayName         string   `json:"displayname"`
	Version             string   `json:"version,omitempty"`
	HideVmwareReadyLogo bool     `json:"hidevmwarereadylogo"`
	LastUpdated         int64    `json:"lastupdated"`
	EntitlementLevel    string   `json:"entitlementlevel"`
	Tiers               []*Tiers `json:"tiersList"`
	VsxId               string   `json:"vsxid"`
}

type CompatibilityMatrix struct {
	ProductId                    string         `json:"productid"`
	VmwareProductId              int32          `json:"vmwareproductid"` // will be deprecated
	VmwareProductName            string         `json:"vmwareproductname"`
	IsPrimary                    bool           `json:"isprimary"`
	PartnerProd                  string         `json:"partnerprod"`
	PartnerProdVer               string         `json:"partnerprodver"`
	ThirdPartyCompany            string         `json:"thirdpartycompany"`
	ThirdPartyProd               string         `json:"thirdpartyprod"`
	ThirdPartyVer                string         `json:"thirdpartyver"`
	SupportStatement             string         `json:"supportstatement"`
	SupportStatementExternalLink string         `json:"supportstatementexternallink"`
	IsVmwareReady                bool           `json:"isvmwareready"`
	CompId                       string         `json:"compid"`
	VmwareProductDetails         *VmwareProduct `json:"vmwareproductdetails"`
	VsxProductId                 string         `json:"vsxproductid"`
	VersionNumber                string         `json:"versionnumber"`
	IsPartnerReady               bool           `json:"ispartnerready"`
	IsNone                       bool           `json:"isnone"`
	CertificationName            string         `json:"certificationname"`
	CertificationDetail          *Certification `json:"certificationdetail"`
}
