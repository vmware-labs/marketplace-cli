// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type Subscription struct {
	ID               int    `json:"id"`
	ConsumerID       string `json:"consumerid"`
	ProductID        string `json:"productid"`
	ProductName      string `json:"productname"`
	ProductLogo      string `json:"productlogo"`
	PublisherID      string `json:"publisherid"`
	PublisherName    string `json:"publishername"`
	DeploymentStatus string `json:"deploymentstatus"`
	DeployedOn       int    `json:"deployedon"`
	SDDCID           string `json:"sddcid"`
	FolderID         string `json:"folderid"`
	ResourcePoolID   string `json:"resourcepoolid"`
	DatastoreID      string `json:"datastoreid"`
	PowerStatus      string `json:"powerstatus"`
	PoweredOn        int    `json:"poweredon"`
	PowerOn          bool   `json:"poweron"`
	StartedOn        int    `json:"startedon"`
	VMName           string `json:"vmname"`
	SourceOrgID      string `json:"sourceorgid"`
	TargetOrgID      string `json:"targetorgid"`
	SDDCLocation     struct {
		Latitude  int `json:"latitude"`
		Longitude int `json:"longitude"`
	} `json:"sddclocation"`
	ProductVersion        string `json:"productversion"`
	EULAAccepted          bool   `json:"eulaaccepted"`
	DeploymentPlatform    string `json:"deploymentplatform"`
	SubscriptionUUID      string `json:"subscriptionuuid"`
	SubscriptionURL       string `json:"subscriptionurl"`
	PlatformRepoName      string `json:"platformreponame"`
	ContainerSubscription struct {
		AppVersion     string `json:"appversion"`
		ChartVersion   string `json:"chartversion"`
		DeploymentType string `json:"deploymenttype"`
	} `json:"containersubscription"`
	StatusText              string `json:"statustext"`
	SourceOrgName           string `json:"sourceorgname"`
	PublisherOrgDisplayName string `json:"publisherorgdisplayname"`
	UpdatesAvailable        bool   `json:"updatesavailable"`
	AutoUpdate              bool   `json:"autoupdate"`
	ContentCatalogID        string `json:"contentcatalogid"`
	PublisherOrgID          string `json:"publisherorgid"`
	SourceOrgDisplayName    string `json:"sourceorgdisplayname"`
	IsAlreadySubscribed     bool   `json:"isalreadysubscribed"`
}
