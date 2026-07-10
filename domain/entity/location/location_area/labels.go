package location_area

// labels.go — LocationArea label structs and DefaultLabels().
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// location sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: LocationAreaLabels -> Labels, LocationArea<Xxx>Labels -> <Xxx>Labels,
// DefaultLocationAreaLabels -> DefaultLabels.

// Labels holds all translatable strings for the location area module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
	Errors  ErrorLabels  `json:"errors"`
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
	AddLocationArea string `json:"add_location_area"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	Active                 string `json:"active"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ErrorLabels struct {
	CannotDeleteInUse string `json:"cannot_delete_in_use"`
}

// DefaultLabels returns sensible English defaults for Labels.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Location Areas",
			HeadingActive:   "Active Location Areas",
			HeadingInactive: "Inactive Location Areas",
			Caption:         "Manage location areas",
			CaptionActive:   "Active location areas",
			CaptionInactive: "Inactive location areas",
		},
		Buttons: ButtonLabels{
			AddLocationArea: "Add Location Area",
		},
		Columns: ColumnLabels{
			Name:        "Name",
			Description: "Description",
			Status:      "Status",
			DateCreated: "Date Created",
		},
		Empty: EmptyLabels{
			ActiveTitle:     "No active location areas",
			ActiveMessage:   "Add your first location area to get started.",
			InactiveTitle:   "No inactive location areas",
			InactiveMessage: "Inactive location areas will appear here.",
		},
		Form: FormLabels{
			Name:                   "Name",
			NamePlaceholder:        "Enter name...",
			Description:            "Description",
			DescriptionPlaceholder: "Enter description...",
			Active:                 "Active",
		},
		Actions: ActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Errors: ErrorLabels{
			CannotDeleteInUse: "Cannot delete — this location area is in use.",
		},
	}
}
