package hyper

// Parameter .
type Parameter struct {
	Label       string        `json:"label,omitempty"`
	Description string        `json:"description,omitempty"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Accept      string        `json:"accept,omitempty"`
	Value       interface{}   `json:"value,omitempty"`
	Options     SelectOptions `json:"options,omitempty"`
	Related     string        `json:"related,omitempty"`
	Components  interface{}   `json:"components,omitempty"`
	Placeholder string        `json:"placeholder,omitempty"`
	Pattern     string        `json:"pattern,omitempty"`    // pattern to be matched by the value
	Min         interface{}   `json:"min,omitempty"`        // minimum value
	MinLength   interface{}   `json:"min-length,omitempty"` // minimum length of value
	Max         interface{}   `json:"max,omitempty"`        // maximum value
	MaxLength   interface{}   `json:"max-length,omitempty"` // maximum length of value
	Step        interface{}   `json:"step,omitempty"`       // granularity to be matched by the parameter's value
	Rows        interface{}   `json:"rows,omitempty"`       // specifies the visible number of lines in a text area
	Cols        interface{}   `json:"cols,omitempty"`       // specifies the visible width of a text area
	Wrap        interface{}   `json:"wrap,omitempty"`       // specifies how the text in a text area is to be wrapped when submitted in a form
	Required    bool          `json:"required,omitempty"`
	ReadOnly    bool          `json:"read-only,omitempty"`
	Multiple    bool          `json:"multiple,omitempty"`
}

// Parameters .
type Parameters []Parameter

// FindByName .
func (as Parameters) FindByName(name string) (Parameter, bool) {
	for _, l := range as {
		if l.Name == name {
			return l, true
		}
	}
	return Parameter{}, false
}

// SelectOption .
type SelectOption struct {
	Label       string         `json:"label,omitempty"`
	Description string         `json:"description,omitempty"`
	Value       interface{}    `json:"value,omitempty"`
	Options     []SelectOption `json:"options,omitempty"`
}

// SelectOptions .
type SelectOptions []SelectOption

func (s SelectOptions) Len() int {
	return len(s)
}

func (s SelectOptions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SelectOptions) Less(i, j int) bool {
	return s[i].Label < s[j].Label
}
