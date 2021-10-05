// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

import (
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
)

type Repo struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type ChartVersion struct {
	Id         string `json:"id,omitempty"`
	Version    string `json:"version,omitempty"`
	AppVersion string `json:"appversion"`
	Details    string `json:"details,omitempty"`
	Readme     string `json:"readme,omitempty"`
	Repo       *Repo  `json:"repo,omitempty"`
	Values     string `json:"values,omitempty"`
	Digest     string `json:"digest,omitempty"`
	Status     string `json:"status,omitempty"`

	TarUrl                         string `json:"tarurl"` // to use during imgprocessor update & download from UI/API
	IsExternalUrl                  bool   `json:"isexternalurl"`
	HelmTarUrl                     string `json:"helmtarurl"` // to use during UI/API create & update product
	IsUpdatedInMarketplaceRegistry bool   `json:"isupdatedinmarketplaceregistry"`
	ProcessingError                string `json:"processingerror"`
	DownloadCount                  int64  `json:"downloadcount"`
	ValidationStatus               string `json:"validationstatus"`
	InstallOptions                 string `json:"installoptions"`
}

type Tier struct {
	Id            int64  `json:"id,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	Description   string `json:"description,omitempty"`
	Icon          string `json:"icon,omitempty"`
	AssetPackURL  string `json:"assetPackURL,omitempty"`
	LearnMoreURL  string `json:"learnMoreURL,omitempty"`
	EnableTooltip bool   `json:"enableTooltip,omitempty"`
	Count         int64  `json:"count,omitempty"`
	VcgProducts   string `json:"vcgProducts,omitempty"`
}

type Tiers struct {
	Id   int64 `json:"id,omitempty"`
	Tier *Tier `json:"tier,omitempty"`
}

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
}

type TechSpecs struct {
	OsDetails            []string `json:"osdetailsList"`
	ContentTypeDetails   []string `json:"contenttypedetailsList"`
	SolutionAreaDetails  []string `json:"solutionareadetailsList"`
	TechSpecsDescription string   `json:"techspecsdescription"`
}

const (
	DeploymentTypesDocker = "DOCKERLINK"
	DeploymentTypeHelm    = "HELM"
)

type ProductItemFile struct {
	Name string `json:"name,omitempty"`
	Size int    `json:"size"`
}

type ProductItemDetails struct {
	Id    string             `json:"id,omitempty"`
	Name  string             `json:"name,omitempty"`
	Files []*ProductItemFile `json:"files,omitempty"`
	Type  string             `json:"type"`
}

type ProductDeploymentFile struct {
	Id              string   `json:"id,omitempty"` // uuid
	Name            string   `json:"name,omitempty"`
	Url             string   `json:"url,omitempty"`
	ImageType       string   `json:"imagetype,omitempty"`
	Status          string   `json:"status,omitempty"`
	UploadedOn      int32    `json:"uploadedon,omitempty"`
	UploadedBy      string   `json:"uploadedby,omitempty"`
	UpdatedOn       int32    `json:"updatedon,omitempty"`
	UpdatedBy       string   `json:"updatedby,omitempty"`
	ItemJson        string   `json:"itemjson,omitempty"`
	Itemkey         string   `json:"itemkey,omitempty"`
	FileID          string   `json:"fileid,omitempty"`
	IsSubscribed    bool     `json:"issubscribed,omitempty"`
	AppVersion      string   `json:"appversion"` // Mandatory
	HashDigest      string   `json:"hashdigest"`
	IsThirdPartyUrl bool     `json:"isthirdpartyurl,omitempty"`
	ThirdPartyUrl   string   `json:"thirdpartyurl,omitempty"`
	IsRedirectUrl   bool     `json:"isredirecturl,omitempty"`
	Comment         string   `json:"comment,omitempty"`
	HashAlgo        string   `json:"hashalgo"`
	DownloadCount   int64    `json:"downloadcount,omitempty"`
	UniqueFileID    string   `json:"uniqueFileId,omitempty"`
	VersionList     []string `json:"versionList"`
}

const (
	HashAlgoSHA1   = "SHA1"
	HashAlgoSHA256 = "SAH256"
)

const (
	ImageTypeJPG  = "JPG"
	ImageTypePNG  = "PNG"
	ImageTypeJPEG = "JPEG"
)

type DeploymentMediaImage struct {
	Status    string `json:"status,omitempty"`
	ImageUrl  string `json:"imageurl"`
	ImageType string `json:"imagetype"`
	CreatedOn int32  `json:"createdon"`
	UpdatedOn int32  `json:"updatedon"`
}

type ProductDeploymentPlatform struct {
	Id          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	ReadyOn     int32  `json:"readyon"`
	DisplayName string `json:"displayname"`
}

type Description struct {
	Summary     string   `json:"summary,omitempty"`
	ImageUrls   []string `json:"imageurlsList"`
	VideoUrls   []string `json:"videourlsList"`
	YoutubeUrl  string   `json:"youtubeurl"`
	Description string   `json:"description,omitempty"`
}

type License struct {
	LicenseName    string `json:"licensename"`
	LicenseDetails string `json:"licensedetails"`
	LicenseUrl     string `json:"licenseurl"`
}

type SupportDetails struct {
	Url         string   `json:"url,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Email       []string `json:"emailList"`
	PhoneNumber []string `json:"phonenumberList"`
}

type Publisher struct {
	UserId         string `json:"userid"`
	OrgId          string `json:"orgid"`
	OrgName        string `json:"orgname"`
	OrgDisplayName string `json:"orgdisplayname"`
}

type EULADetails struct {
	Url       string `json:"url,omitempty"`
	Text      string `json:"text,omitempty"` // mandatory
	Signed    bool   `json:"signed,omitempty"`
	SignedOn  int32  `json:"signedon"`
	Version   string `json:"version,omitempty"`
	CreatedOn int32  `json:"createdon"`
	UpdatedOn int32  `json:"updatedon"`
}

type ProductEncryptionDetails struct {
	List []string `json:"listList"`
}

type ProductEncryption struct {
	List map[string]bool `json:"list"`
}

type ProductExportCompliance struct {
	Eccn             string `json:"eccn,omitempty"`
	HtsNumber        string `json:"htsnumber"`
	LicenseException string `json:"licenseexception"`
	CcatsNumber      string `json:"ccatsnumber"`
	CcatsDocumentUrl string `json:"ccatsdocumenturl"`
}

type OpenSourceDisclosureURLS struct {
	LicenseDisclosureURL string `json:"licensedisclosureurl"`
	SourceCodePackageURL string `json:"sourcecodepackageurl"`
}

type Logo struct {
	URL          string `json:"url"`
	CreationDate int    `json:"createdon"`
}

type Version struct {
	Number       string `json:"versionnumber"`
	Details      string `json:"versiondetails"`
	Status       string `json:"status,omitempty"`
	Instructions string `json:"versioninstruction"`
}

type DockerImageTag struct {
	ID                             string `json:"id,omitempty"`
	Tag                            string `json:"tag,omitempty"`
	Type                           string `json:"type,omitempty"`
	IsUpdatedInMarketplaceRegistry bool   `json:"isupdatedinmarketplaceregistry"`
	MarketplaceS3Link              string `json:"marketplaces3link"`
	AppCheckReportLink             string `json:"appcheckreportlink"`
	AppCheckSummaryPdfLink         string `json:"appchecksummarypdflink"`
	S3TarBackupUrl                 string `json:"s3tarbackupurl"`
	ProcessingError                string `json:"processingerror"`
	DownloadCount                  int64  `json:"downloadcount"`
	HashAlgo                       int    `json:"hashalgo"`
	HashDigest                     string `json:"hashdigest"`
}

type DockerURLDetails struct {
	ID                    string            `json:"id,omitempty"`
	Key                   string            `json:"key,omitempty"`
	Url                   string            `json:"url,omitempty"`
	MarketplaceUpdatedUrl string            `json:"marketplaceupdatedurl"`
	ImageTags             []*DockerImageTag `json:"imagetagsList"`
	ImageTagsAsJson       string            `json:"imagetagsasjson"`
	DeploymentInstruction string            `json:"deploymentinstruction"`
}

func (d *DockerURLDetails) GetTag(tagName string) *DockerImageTag {
	for _, tag := range d.ImageTags {
		if tag.Tag == tagName {
			return tag
		}
	}
	return nil
}

func (d *DockerURLDetails) HasTag(tagName string) bool {
	return d.GetTag(tagName) != nil
}

type DockerVersionList struct {
	ID                    string              `json:"id,omitempty"`
	AppVersion            string              `json:"appversion"`
	DeploymentInstruction string              `json:"deploymentinstruction"`
	DockerURLs            []*DockerURLDetails `json:"dockerurlsList"`
	Status                string              `json:"status,omitempty"`
	//ImageTags             []*DockerImageTag   `json:"imagetags"` // DEPRECATED
}

func (l *DockerVersionList) GetImage(imageURL string) *DockerURLDetails {
	for _, image := range l.DockerURLs {
		if image.Url == imageURL {
			return image
		}
	}
	return nil
}

type RateCard struct {
	RateCardId        string               `json:"ratecardid"`
	SubscriptionType  string               `json:"subscriptiontype"`
	DimensionPricing  []*RateCardDimension `json:"dimensionpricingList"`
	SubscriptionPrice float32              `json:"subscriptionprice"`
}

type RateCardDimension struct {
	DimensionName  string  `json:"dimensionname"`
	DimensionPrice float32 `json:"dimensionprice"`
	DimensionUnit  string  `json:"dimensionunit"`
}

const (
	SolutionTypeHelm = "HELMCHARTS"
)

type Product struct {
	ProductId        string          `json:"productid"`
	IsParent         bool            `json:"isparent"`
	Slug             string          `json:"slug,omitempty"`
	DisplayName      string          `json:"displayname"`
	Published        bool            `json:"published,omitempty"`
	TechSpecs        *TechSpecs      `json:"techspecs"`
	Description      *Description    `json:"description,omitempty"`
	License          *License        `json:"license,omitempty"`
	Categories       []string        `json:"categoriesList"`
	SupportAvailable bool            `json:"supportavailable"`
	SupportDetails   *SupportDetails `json:"supportdetails"`
	LoginRequired    bool            `json:"loginrequired"`
	PublisherDetails *Publisher      `json:"publisherdetails"` // Mandatory
	Type             string          `json:"type,omitempty"`
	ProductPricing   []*RateCard     `json:"productpricingList"`
	//Resources []*AppProductResources `json:"resourcesList"`
	//MetaDetails *AppProductMetaDetails `json:"metadetails"`
	Status          string       `json:"status,omitempty"`
	ParentProductId string       `json:"parentproductid"`
	Byol            bool         `json:"byol,omitempty"`
	EulaDetails     *EULADetails `json:"euladetails"`
	EulaURL         string       `json:"eulaurl"`
	EulaTempURL     string       `json:"eulatempurl"`
	//Highlights      []string     `json:"highlightsList"`
	ProductDeploymentMediaImages []*DeploymentMediaImage      `json:"productdeploymentmediaimagesList"`
	ProductDeploymentFiles       []*ProductDeploymentFile     `json:"productdeploymentfilesList"`
	SaasURL                      string                       `json:"saasurl"`
	DeploymentPlatforms          []*ProductDeploymentPlatform `json:"deploymentplatformsList"`
	CreationDate                 int                          `json:"createdon"`
	UpdatedDate                  int                          `json:"updatedon"`
	UpdatedBy                    string                       `json:"updatedby"`
	PublishedDate                int                          `json:"publishedon"`
	PublisherOrgName             string                       `json:"publisherorgname"`
	EncryptionDetails            *ProductEncryptionDetails    `json:"encryptiondetails"`
	Encryption                   *ProductEncryption           `json:"encryption"`
	ExportCompliance             *ProductExportCompliance     `json:"exportcompliance"`
	OpenSourceDisclosure         *OpenSourceDisclosureURLS    `json:"opensourcedisclosure"`
	ProductLogo                  *Logo                        `json:"productlogo"`
	InLegalReview                bool                         `json:"inlegalreview"`
	IsVSX                        bool                         `json:"isvsx"`
	Versions                     []*Version                   `json:"versionsList"`
	AllVersions                  []*Version                   `json:"allversiondetailsList"`
	DeploymentTypes              []string                     `json:"deploymenttypesList"`
	SolutionType                 string                       `json:"solutiontype"`
	FormFactor                   string                       `json:"formfactor"`
	ChartVersions                []*ChartVersion              `json:"chartversionsList"`
	//Blueprints []*ProductBlueprintDetails `json:"blueprintsList"`
	//BlueprintDetails             *BlueprintDetails       `json:"blueprintdetails"`
	//DockerURLs                   []*DockerURL            `json:"dockerurlsList"`
	DockerLinkVersions   []*DockerVersionList   `json:"dockerlinkversionsList"`
	RelatedSlugs         []string               `json:"relatedslugsList"`
	IsFeatured           bool                   `json:"isfeatured"`
	IsPopular            bool                   `json:"isPopular,omitempty"`
	IsPrivate            bool                   `json:"isprivate"`
	IsListingProduct     bool                   `json:"islistingproduct"`
	CompatibilityMatrix  []*CompatibilityMatrix `json:"compatibilitymatrixList"` // compatibility-matrix-supported-features needed for vsx.
	CertificationType    []string               `json:"certificationtypeList"`
	SolutionAreaId       []string               `json:"solutionareaidList"`
	SolutionAreaName     []string               `json:"solutionareanameList"`
	SolutionAreaTypeId   []string               `json:"solutionareatypeidList"`
	SolutionAreaTypeName []string               `json:"solutionareatypenameList"`
	Category             []string               `json:"categoryList"`
	SubCategories        []string               `json:"subcategoriesList"`
	SubCategoryId        []string               `json:"subcategoryidList"`

	Version                *Version `json:"version,omitempty"`
	IsDraft                bool     `json:"isdraft"`
	DeploymentType         string   `json:"deploymenttype"`         // TO BE DEPRECATED
	DeploymentInstructions string   `json:"deploymentinstructions"` // TO BE DEPRECATED
	Unsubscribable         bool     `json:"unsubscribable,omitempty"`
	CannotDownload         bool     `json:"cannotdownload"`
	//VsxDetails             *VSXDetails     `json:"vsxdetails"`
	//SaaSProduct            *SaaSProduct    `json:"saasproduct"`
	DraftId     string `pjson:"draftid"`
	ChartId     string `json:"chartid"`
	IsAutoDraft bool   `json:"isautodraft"`
	//ProductAddOnFiles []*AddOnFiles `json:"addonfilesList"`
}

func (product *Product) GetVersion(version string) *Version {
	if version == "" {
		return product.GetLatestVersion()
	}

	for _, v := range product.AllVersions {
		if v.Number == version {
			return v
		}
	}
	return nil
}

func (product *Product) GetLatestVersion() *Version {
	if len(product.AllVersions) == 0 {
		return &Version{Number: "N/A"}
	}

	version, err := product.getLatestVersionSemver()
	if err != nil {
		version = product.getLatestVersionAlphanumeric()
	}

	return version
}

func (product *Product) getLatestVersionSemver() (*Version, error) {
	latestVersion := product.AllVersions[0]
	version, err := semver.NewVersion(latestVersion.Number)
	if err != nil {
		return nil, err
	}
	for _, v := range product.AllVersions {
		otherVersion, err := semver.NewVersion(v.Number)
		if err != nil {
			return nil, err
		}
		if version.LessThan(*otherVersion) {
			latestVersion = v
			version = otherVersion
		}
	}

	return latestVersion, nil
}

func (product *Product) getLatestVersionAlphanumeric() *Version {
	latestVersion := product.AllVersions[0]
	for _, v := range product.AllVersions {
		if strings.Compare(latestVersion.Number, v.Number) < 0 {
			latestVersion = v
		}
	}
	return latestVersion
}

type Versions []*Version

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Versions) Less(i, j int) bool {
	return v[i].LessThan(*v[j])
}

func (a Version) LessThan(b Version) bool {
	semverA, errA := semver.NewVersion(a.Number)
	semverB, errB := semver.NewVersion(b.Number)

	if errA != nil || errB != nil {
		return strings.Compare(a.Number, b.Number) < 0
	}

	return semverA.LessThan(*semverB)
}

func Sort(versions []*Version) {
	sort.Sort(sort.Reverse(Versions(versions)))
}

func (product *Product) HasVersion(version string) bool {
	return product.GetVersion(version) != nil
}

func (product *Product) GetContainerImagesForVersion(version string) *DockerVersionList {
	for _, dockerVersionLink := range product.DockerLinkVersions {
		if dockerVersionLink.AppVersion == product.GetVersion(version).Number {
			return dockerVersionLink
		}
	}
	return nil
}

func (product *Product) GetChartsForVersion(version string) []*ChartVersion {
	var charts []*ChartVersion
	for _, chart := range product.ChartVersions {
		if chart.AppVersion == product.GetVersion(version).Number {
			charts = append(charts, chart)
		}
	}
	return charts
}

func (product *Product) GetChart(chartId string) *ChartVersion {
	for _, chart := range product.ChartVersions {
		if chart.Id == chartId {
			return chart
		}
	}
	return nil
}

func (product *Product) AddChart(chart *ChartVersion) {
	product.ChartVersions = append(product.ChartVersions, chart)
}

func (product *Product) GetFilesForVersion(version string) []*ProductDeploymentFile {
	var files []*ProductDeploymentFile
	for _, file := range product.ProductDeploymentFiles {
		if file.AppVersion == product.GetVersion(version).Number {
			files = append(files, file)
		}
	}
	return files
}

func (product *Product) GetFile(fileId string) *ProductDeploymentFile {
	for _, file := range product.ProductDeploymentFiles {
		if file.FileID == fileId {
			return file
		}
	}
	return nil
}

func (product *Product) AddFile(file *ProductDeploymentFile) {
	product.ProductDeploymentFiles = append(product.ProductDeploymentFiles, file)
}

func (product *Product) PrepForUpdate() {
	// Send an empty compatibility matrix, any entries in here will multiply
	product.CompatibilityMatrix = []*CompatibilityMatrix{}

	// For updates, the encryption hash needs to be populated
	// with the contents of the encryption details list
	product.Encryption = &ProductEncryption{List: map[string]bool{}}
	if product.EncryptionDetails != nil {
		for _, key := range product.EncryptionDetails.List {
			product.Encryption.List[key] = true
		}
	}

	// On updates, there is no Versions, only AllVersions, so
	// make sure AllVersions truly has all versions
	for _, version := range product.Versions {
		if !product.HasVersion(version.Number) {
			product.AllVersions = append(product.AllVersions, version)
		}
	}
	product.Versions = product.AllVersions

	product.ChartVersions = []*ChartVersion{}
	product.DockerLinkVersions = []*DockerVersionList{}
	product.ProductDeploymentFiles = []*ProductDeploymentFile{}
}

func (product *Product) SetDeploymentType(deploymentType string) {
	product.DeploymentTypes = []string{deploymentType}
}
