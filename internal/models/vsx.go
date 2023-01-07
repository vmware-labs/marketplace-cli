// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type VSXDeveloper struct {
	ID            int64  `json:"id"`
	VSXID         int64  `json:"vsxid"`
	Partner       bool   `json:"partner"`
	Slug          string `json:"slug"`
	DisplayName   string `json:"displayname"`
	Description   string `json:"description"`
	Icon          string `json:"icon"`
	URL           string `json:"url"`
	SupportEmail  string `json:"supportemail"`
	SupportURL    string `json:"supporturl"`
	SupportPhone  string `json:"supportphone"`
	TwitterUser   string `json:"twitteruser"`
	FacebookURL   string `json:"facebookurl"`
	LinkedInURL   string `json:"linkedinurl"`
	YouTubeURL    string `json:"youtubeurl"`
	State         string `json:"state"`
	ReviewComment string `json:"reviewcomment"`
	CreatedOn     int64  `json:"created"`
	UpdatedOn     int64  `json:"lastupdated"`
	Count         int64  `json:"count"`
}

type Category struct {
	ID       int64        `json:"id"`
	Category *VSXCategory `json:"category"`
}
type VSXCategory struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	First       string `json:"first"`
	Second      string `json:"second"`
	Hidden      bool   `json:"hidden"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Link        string `json:"link"`
	Count       int64  `json:"count"`
}
type VSXContentType struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"displayname"`
	Count       int64  `json:"count"`
}

type RelatedProduct struct {
	ID                    int64  `json:"id"`
	ShortName             string `json:"shortname"`
	DisplayName           string `json:"displayname"`
	Version               string `json:"version"`
	EntitlementLevel      string `json:"entitlementlevel"`
	ContentTypes          string `json:"contenttypes"`
	Count                 int64  `json:"count"`
	ShortNameCount        int64  `json:"shortnamecount"`
	EntitlementLevelCount int64  `json:"entitlementlevelcount"`
	LastUpdated           int64  `json:"lastupdated"`
}

type VSXRelatedProducts struct {
	ID                           int64           `json:"id"`
	Product                      *RelatedProduct `json:"product"`
	VMwareReady                  bool            `json:"vmwareready"`
	SupportStatement             string          `json:"supportstatement"`
	SupportStatementExternalLink bool            `json:"supportstatementexternallink"`
	PartnerProduct               string          `json:"partnerproduct"`
	PartnerProductVersion        string          `json:"partnerproductversion"`
	ThirdPartyCompany            string          `json:"thirdpartycompany"`
	ThirdPartyProduct            string          `json:"thirdpartyproduct"`
	ThirdPartyProductVersion     string          `json:"thirdpartyproductversion"`
	PrimaryProduct               bool            `json:"primaryproduct"`
}

type OperatingSystem struct {
	ID       int64        `json:"id"`
	Category *VSXCategory `json:"category"`
}

type SolutionArea struct {
	ID       int64        `json:"id"`
	Category *VSXCategory `json:"category"`
}

type Technology struct {
	ID       int64        `json:"id"`
	Category *VSXCategory `json:"category"`
}

type VSXDetails struct {
	ID                int64                 `json:"id"`
	UUID              string                `json:"uuid"`
	ShortName         string                `json:"shortname"`
	Version           string                `json:"version"`
	Revision          int32                 `json:"revision"`
	Featured          bool                  `json:"featured"`
	Banner            bool                  `json:"banner"`
	SupportsInProduct bool                  `json:"supportsinproduct"`
	Restricted        bool                  `json:"restricted"`
	Developer         *VSXDeveloper         `json:"developer"`
	ContentType       *VSXContentType       `json:"contenttype"`
	VSGURL            string                `json:"vsgurl"`
	OtherMetadata     string                `json:"othermetadata"`
	NumViews          int64                 `json:"numviews"`
	NumInstalls       int64                 `json:"numinstalls"`
	NumReviews        int64                 `json:"numreviews"`
	AvgRating         float32               `json:"avgrating"`
	Categories        []*Category           `json:"categoriesList"`
	SolutionAreas     []*SolutionArea       `json:"solutionareasList"`
	Technologies      []*Technology         `json:"technologiesList"`
	OperatingSystems  []*OperatingSystem    `json:"operatingsystemsList"`
	Tiers             []*Tiers              `json:"tiersList"`
	Products          []*VSXRelatedProducts `json:"productsList"`
	ParentID          int64                 `json:"parentid"`
}
