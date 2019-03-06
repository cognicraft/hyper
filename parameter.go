package hyper

// Parameter .
type Parameter struct {
	Label       string        `json:"label,omitempty"`
	Description string        `json:"description,omitempty"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Value       interface{}   `json:"value,omitempty"`
	Placeholder string        `json:"placeholder,omitempty"`
	Options     SelectOptions `json:"options,omitempty"`
	Related     string        `json:"related,omitempty"`
	Components  interface{}   `json:"components,omitempty"`
	Pattern     string        `json:"pattern,omitempty"`
	Min         interface{}   `json:"min,omitempty"`
	Max         interface{}   `json:"max,omitempty"`
	MaxLength   interface{}   `json:"max-length,omitempty"`
	Size        interface{}   `json:"size,omitempty"`
	Step        interface{}   `json:"step,omitempty"`
	Cols        interface{}   `json:"cols,omitempty"`
	Rows        interface{}   `json:"rows,omitempty"`
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
