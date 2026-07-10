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
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	AddLocation string `json:"add_location"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Timezone    string `json:"timezone"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type FormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"address_placeholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	Timezone               string `json:"timezone"`
	Area                   string `json:"area"`
	AreaPlaceholder        string `json:"area_placeholder"`
	Active                 string `json:"active"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	AddressInfo     string `json:"address_info"`
	DescriptionInfo string `json:"description_info"`
	TimezoneInfo    string `json:"timezone_info"`
	AreaInfo        string `json:"area_info"`
	ActiveInfo      string `json:"active_info"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type DetailLabels struct {
	BasicInfo   DetailBasicInfoLabels `json:"basic_info"`
	Tabs        DetailTabLabels       `json:"tabs"`
	EmptyStates DetailEmptyLabels     `json:"empty_states"`
	// Inline feedback messages
	UpdateSuccess string `json:"update_success"`
	UpdateError   string `json:"update_error"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachments_tab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"audit_history_tab"`
}

type DetailBasicInfoLabels struct {
	Title                  string `json:"title"`
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"address_placeholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	Active                 string `json:"active"`
	Save                   string `json:"save"`
}

type DetailTabLabels struct {
	Info       string `json:"info"`
	Users      string `json:"users"`
	PriceLists string `json:"price_lists"`
	AuditTrail string `json:"audit_trail"`
}

type DetailEmptyLabels struct {
	UsersTitle      string `json:"users_title"`
	UsersDesc       string `json:"users_desc"`
	PriceListsTitle string `json:"price_lists_title"`
	PriceListsDesc  string `json:"price_lists_desc"`
	AuditTitle      string `json:"audit_title"`
	AuditDesc       string `json:"audit_desc"`
}

// DashboardLabels holds translatable strings for the location dashboard.
type DashboardLabels struct {
	// Stats (4): Total / Active / Regions / Areas Count
	TotalLocations string `json:"total_locations"`
	Active         string `json:"active"`
	Regions        string `json:"regions"`
	AreasCount     string `json:"areas_count"`

	// Widget titles
	LocationsByRegion  string `json:"locations_by_region"`
	TopLocationsByArea string `json:"top_locations_by_area"`
	RecentAdditions    string `json:"recent_additions"`
	ViewAll            string `json:"view_all"`

	// Chart filter labels
	FilterWeek  string `json:"filter_week"`
	FilterMonth string `json:"filter_month"`
	FilterYear  string `json:"filter_year"`

	// Quick action labels
	QuickNewLocation string `json:"quick_new_location"`
	QuickNewArea     string `json:"quick_new_area"`

	// Activity / table column labels
	ColumnLocation string `json:"column_location"`
	ColumnAreas    string `json:"column_areas"`
	LocationAdded  string `json:"location_added"`
}
