package location

// labels.go — Location label structs and the location dashboard labels that the
// location view module owns.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// location sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: LocationLabels -> Labels, Location<Xxx>Labels -> <Xxx>Labels,
// LocationDashboardLabels -> DashboardLabels.

// Labels holds all translatable strings for the location module.
// JSON tags match the "location" wrapper key in retail/location.json.
type Labels struct {
	Page      PageLabels      `json:"page"`
	Buttons   ButtonLabels    `json:"buttons"`
	Columns   ColumnLabels    `json:"columns"`
	Empty     EmptyLabels     `json:"empty"`
	Form      FormLabels      `json:"form"`
	Actions   ActionLabels    `json:"actions"`
	Detail    DetailLabels    `json:"detail"`
	Dashboard DashboardLabels `json:"dashboard"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddLocation string `json:"addLocation"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Timezone    string `json:"timezone"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"addressPlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Timezone               string `json:"timezone"`
	Area                   string `json:"area"`
	AreaPlaceholder        string `json:"areaPlaceholder"`
	Active                 string `json:"active"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	AddressInfo     string `json:"addressInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	TimezoneInfo    string `json:"timezoneInfo"`
	AreaInfo        string `json:"areaInfo"`
	ActiveInfo      string `json:"activeInfo"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type DetailLabels struct {
	BasicInfo   DetailBasicInfoLabels `json:"basicInfo"`
	Tabs        DetailTabLabels       `json:"tabs"`
	EmptyStates DetailEmptyLabels     `json:"emptyStates"`
	// Inline feedback messages
	UpdateSuccess string `json:"updateSuccess"`
	UpdateError   string `json:"updateError"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachmentsTab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"auditHistoryTab"`
}

type DetailBasicInfoLabels struct {
	Title                  string `json:"title"`
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"addressPlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Active                 string `json:"active"`
	Save                   string `json:"save"`
}

type DetailTabLabels struct {
	Info       string `json:"info"`
	Users      string `json:"users"`
	PriceLists string `json:"priceLists"`
	AuditTrail string `json:"auditTrail"`
}

type DetailEmptyLabels struct {
	UsersTitle      string `json:"usersTitle"`
	UsersDesc       string `json:"usersDesc"`
	PriceListsTitle string `json:"priceListsTitle"`
	PriceListsDesc  string `json:"priceListsDesc"`
	AuditTitle      string `json:"auditTitle"`
	AuditDesc       string `json:"auditDesc"`
}

// DashboardLabels holds translatable strings for the location dashboard.
type DashboardLabels struct {
	// Stats (4): Total / Active / Regions / Areas Count
	TotalLocations string `json:"totalLocations"`
	Active         string `json:"active"`
	Regions        string `json:"regions"`
	AreasCount     string `json:"areasCount"`

	// Widget titles
	LocationsByRegion  string `json:"locationsByRegion"`
	TopLocationsByArea string `json:"topLocationsByArea"`
	RecentAdditions    string `json:"recentAdditions"`
	ViewAll            string `json:"viewAll"`

	// Chart filter labels
	FilterWeek  string `json:"filterWeek"`
	FilterMonth string `json:"filterMonth"`
	FilterYear  string `json:"filterYear"`

	// Quick action labels
	QuickNewLocation string `json:"quickNewLocation"`
	QuickNewArea     string `json:"quickNewArea"`

	// Activity / table column labels
	ColumnLocation string `json:"columnLocation"`
	ColumnAreas    string `json:"columnAreas"`
	LocationAdded  string `json:"locationAdded"`
}
