// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

const (
	DeploymentStatusActive          = "ACTIVE"
	DeploymentStatusInactive        = "INACTIVE"
	DeploymentStatusApprovalPending = "APPROVAL_PENDING"
	DeploymentStatusNotProcessed    = "NOT_PROCESSED"
)

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

const (
	HashAlgoSHA1   = "SHA1"
	HashAlgoSHA256 = "SHA256"
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
	List                  []string `json:"listList"`
	NonstandardEncryption bool     `json:"nonstandardencryption"`
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

type AppProductResources struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	IsDownloadable bool   `json:"isdownloadable"`
}

type AppProductMetaDetails struct {
	OfficialEmail       string `json:"officialemail"`
	OfficialPhoneNumber string `json:"officialphonenumber"`
	WebsiteURL          string `json:"websiteurl"`
	AppID               string `json:"appid"`
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

type PCADetail struct {
	URL          string `json:"url"`
	Version      string `json:"version"`
	CreatedOn    int32  `json:"createdon"`
	UpdatedOn    int32  `json:"updatedon"`
	PresignedURL string `json:"presignedurl"`
}

const (
	SolutionTypeChart  = "HELMCHARTS"
	SolutionTypeImage  = "CONTAINER"
	SolutionTypeISO    = "ISO"
	SolutionTypeOthers = "OTHERS"
	SolutionTypeOVA    = "OVA"
)

type Product struct {
	ProductId                    string                       `json:"productid"`
	PublishedProductId           string                       `json:"publishedproductid"`
	IsParent                     bool                         `json:"isparent"`
	Slug                         string                       `json:"slug,omitempty"`
	DisplayName                  string                       `json:"displayname"`
	IsPublished                  bool                         `json:"ispublished,omitempty"` // Set on product list requests
	Published                    bool                         `json:"published,omitempty"`   // Set on product get requests
	TechSpecs                    *TechSpecs                   `json:"techspecs"`
	Description                  *Description                 `json:"description,omitempty"`
	License                      *License                     `json:"license,omitempty"`
	Categories                   []string                     `json:"categoriesList"`
	SupportAvailable             bool                         `json:"supportavailable"`
	SupportDetails               *SupportDetails              `json:"supportdetails"`
	LoginRequired                bool                         `json:"loginrequired"`
	PublisherDetails             *Publisher                   `json:"publisherdetails"` // Mandatory
	Type                         string                       `json:"type,omitempty"`
	ProductPricing               []*RateCard                  `json:"productpricingList"`
	Resources                    []*AppProductResources       `json:"resourcesList"`
	MetaDetails                  *AppProductMetaDetails       `json:"metadetails"`
	Status                       string                       `json:"status,omitempty"`
	ParentProductId              string                       `json:"parentproductid"`
	Byol                         bool                         `json:"byol,omitempty"`
	EulaDetails                  *EULADetails                 `json:"euladetails"`
	EulaURL                      string                       `json:"eulaurl"`
	EulaTempURL                  string                       `json:"eulatempurl"`
	Highlights                   []string                     `json:"highlightsList"`
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
	Logo                         string                       `json:"logo"`
	ProductLogo                  *Logo                        `json:"productlogo"`
	InLegalReview                bool                         `json:"inlegalreview"`
	IsVSX                        bool                         `json:"isvsx"`
	AllVersions                  []*Version                   `json:"allversiondetailsList"`
	LatestVersion                string                       `json:"latestversion"`
	Version                      *Version                     `json:"version,omitempty"`
	Versions                     []*Version                   `json:"versionsList"`
	CurrentVersion               string                       `json:"currentversion"`
	DeploymentTypes              []string                     `json:"deploymenttypesList"`
	SolutionType                 string                       `json:"solutiontype"`
	FormFactor                   string                       `json:"formfactor"`
	ChartVersions                []*ChartVersion              `json:"chartversionsList"`
	Blueprints                   []*ProductBlueprintDetails   `json:"blueprintsList"`
	DockerURLs                   []*DockerURLDetails          `json:"dockerurlsList"`
	DockerLinkVersions           []*DockerVersionList         `json:"dockerlinkversionsList"`
	RelatedSlugs                 []string                     `json:"relatedslugsList"`
	VSXDetails                   *VSXDetails                  `json:"vsxdetails"`
	IsFeatured                   bool                         `json:"isfeatured"`
	IsPopular                    bool                         `json:"isPopular,omitempty"`
	IsPrivate                    bool                         `json:"isprivate"`
	IsListingProduct             bool                         `json:"islistingproduct"`
	CompatibilityMatrix          []*CompatibilityMatrix       `json:"compatibilitymatrixList"` // compatibility-matrix-supported-features needed for vsx.
	CompatiblePlatformIDList     []string                     `json:"compatibleplatformidList"`
	CompatiblePlatformNameList   []string                     `json:"compatibleplatformnameList"`
	CertificationList            []*Certification             `json:"certificationList"`
	CertificationTypes           []string                     `json:"certificationtypeList"`
	SolutionAreaId               []string                     `json:"solutionareaidList"`
	SolutionAreaName             []string                     `json:"solutionareanameList"`
	SolutionAreaTypeId           []string                     `json:"solutionareatypeidList"`
	SolutionAreaTypeName         []string                     `json:"solutionareatypenameList"`
	Category                     []string                     `json:"categoryList"`
	SubCategories                []string                     `json:"subcategoriesList"`
	SubCategoryId                []string                     `json:"subcategoryidList"`
	DeploymentType               string                       `json:"deploymenttype"`
	DeploymentInstructions       string                       `json:"deploymentinstructions"`
	Unsubscribable               bool                         `json:"unsubscribable,omitempty"`
	CannotDownload               bool                         `json:"cannotdownload"`
	IsDraft                      bool                         `json:"isdraft"`
	IsAutoDraft                  bool                         `json:"isautodraft"`
	DraftId                      string                       `json:"draftid"`
	ChartId                      string                       `json:"chartid"`
	AddOnFiles                   []*AddOnFile                 `json:"addonfilesList"`
	ProductAddOnFiles            []*AddOnFile                 `json:"productaddonfilesList"`
	Tags                         []string                     `json:"tagsList"`
	SKUS                         []*SKUPublisherView          `json:"skusList"`
	MetaFiles                    []*MetaFile                  `json:"metafilesList"`
	PCADetails                   *PCADetail                   `json:"pcadetails"`
	RedirectURL                  string                       `json:"redirecturl"`
}

type VersionSpecificProductDetails struct {
	EncryptionDetails      *ProductEncryptionDetails  `json:"encryptiondetails"`
	EulaDetails            *EULADetails               `json:"euladetails"`
	EulaURL                string                     `json:"eulaurl"`
	EulaTempURL            string                     `json:"eulatempurl"`
	ExportCompliance       *ProductExportCompliance   `json:"exportcompliance"`
	OpenSourceDisclosure   *OpenSourceDisclosureURLS  `json:"opensourcedisclosure"`
	CertificationList      []*Certification           `json:"certificationList"`
	CertificationTypes     []string                   `json:"certificationtypesList"`
	ProductDeploymentFiles []*ProductDeploymentFile   `json:"productdeploymentfilesList"`
	DockerLinkVersions     []*DockerVersionList       `json:"dockerlinkversionsList"`
	ChartVersions          []*ChartVersion            `json:"chartversionsList"`
	Blueprints             []*ProductBlueprintDetails `json:"blueprintsList"`
	AddOnFiles             []*AddOnFile               `json:"addonfilesList"`
	CreationDate           int                        `json:"createdon"`
	UpdatedDate            int                        `json:"updatedon"`
	UpdatedBy              string                     `json:"updatedby"`
	PublishedDate          int                        `json:"publishedon"`
	CompatibilityMatrix    []*CompatibilityMatrix     `json:"compatibilitymatrixList"` // compatibility-matrix-supported-features needed for vsx.
	HasLimitedAccess       bool                       `json:"haslimitedaccess"`
	Tag                    string                     `json:"tag"`
	MetaFiles              []*MetaFile                `json:"metafilesList"`
	PCADetails             *PCADetail                 `json:"pcadetails"`
}

func (product *Product) UpdateWithVersionSpecificDetails(version string, details *VersionSpecificProductDetails) {
	product.EncryptionDetails = details.EncryptionDetails
	product.EulaDetails = details.EulaDetails
	product.EulaURL = details.EulaURL
	product.EulaTempURL = details.EulaTempURL
	product.ExportCompliance = details.ExportCompliance
	product.OpenSourceDisclosure = details.OpenSourceDisclosure
	product.CertificationList = details.CertificationList
	product.CertificationTypes = details.CertificationTypes
	product.ProductDeploymentFiles = details.ProductDeploymentFiles
	product.DockerLinkVersions = details.DockerLinkVersions
	product.ChartVersions = details.ChartVersions
	product.Blueprints = details.Blueprints
	product.AddOnFiles = details.AddOnFiles
	product.CreationDate = details.CreationDate
	product.UpdatedDate = details.UpdatedDate
	product.UpdatedBy = details.UpdatedBy
	product.PublishedDate = details.PublishedDate
	product.CompatibilityMatrix = details.CompatibilityMatrix
	product.GetVersion(version).HasLimitedAccess = details.HasLimitedAccess
	product.GetVersion(version).Tag = details.Tag
	product.MetaFiles = details.MetaFiles
	product.PCADetails = details.PCADetails
}

func (product *Product) PrepForUpdate() {
	// This whole function is a workaround

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
}

func (product *Product) SetPCAFile(version, pcaURL string) {
	product.PCADetails = &PCADetail{
		URL:     pcaURL,
		Version: version,
	}
}
