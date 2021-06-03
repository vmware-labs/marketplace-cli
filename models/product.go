// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

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
	DeploymentTypeHelm    = "HELM"
	DeploymentTypesDocker = "DOCKERLINK"
)

type ProductDeploymentFile struct {
	Id              string `json:"id,omitempty"` // uuid
	Name            string `json:"name,omitempty"`
	Url             string `json:"url,omitempty"`
	ImageType       string `json:"imagetype"`
	Status          string `json:"status,omitempty"`
	UploadedOn      int32  `json:"uploadedon"`
	UploadedBy      string `json:"uploadedby"`
	UpdatedOn       int32  `json:"updatedon"`
	UpdatedBy       string `json:"updatedby"`
	ItemJson        string `json:"itemjson"`
	Itemkey         string `json:"itemkey,omitempty"`
	FileID          string `json:"fileid"`
	IsSubscribed    bool   `json:"issubscribed"`
	AppVersion      string `json:"appversion"` // Mandatory
	HashDigest      string `json:"hashdigest"`
	IsThirdPartyUrl bool   `json:"isthirdpartyurl"`
	ThirdPartyUrl   string `json:"thirdpartyurl"`
	Comment         string `json:"comment,omitempty"`
	HashAlgo        string `json:"hashalgo"`
	DownloadCount   int64  `json:"downloadcount"`
}

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
	Id                    string              `json:"id,omitempty"`
	AppVersion            string              `json:"appversion"`
	DeploymentInstruction string              `json:"deploymentinstruction"`
	DockerURLs            []*DockerURLDetails `json:"dockerurlsList"`
	Status                string              `json:"status,omitempty"`
	//ImageTags             []*DockerImageTag   `json:"imagetags"` // DEPRECATED
}

func (l *DockerVersionList) GetImage(imageName string) *DockerURLDetails {
	for _, image := range l.DockerURLs {
		if image.Url == imageName {
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
	CompatibilityMatrix  []*CompatibilityMatrix `json:"compatibilitymatrixList"` // compatability-matrix-supported-features needed for vsx.
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
	for _, v := range product.Versions {
		if v.Number == version {
			return v
		}
	}
	return nil
}

func (product *Product) HasVersion(version string) bool {
	return product.GetVersion(version) != nil
}

func (product *Product) GetDockerImagesForVersion(version string) *DockerVersionList {
	for _, dockerVersionLink := range product.DockerLinkVersions {
		if dockerVersionLink.AppVersion == version {
			return dockerVersionLink
		}
	}
	return nil
}

func (product *Product) GetChartsForVersion(version string) []*ChartVersion {
	var charts []*ChartVersion
	for _, chart := range product.ChartVersions {
		if chart.AppVersion == version {
			charts = append(charts, chart)
		}
	}
	return charts
}

func (product *Product) PrepForUpdate() {
	product.Encryption = &ProductEncryption{List: map[string]bool{}}
	if product.EncryptionDetails != nil {
		for _, key := range product.EncryptionDetails.List {
			product.Encryption.List[key] = true
		}
	}
}

func (product *Product) SetDeploymentType(deploymentType string) {
	for _, depType := range product.DeploymentTypes {
		if depType == deploymentType {
			return
		}
	}
	product.DeploymentTypes = append(product.DeploymentTypes, deploymentType)
}
