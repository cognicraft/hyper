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

//Find
func (ps Parameters) Find(accept func(p Parameter) bool) (Parameter, bool) {
	for _, p := range ps {
		if accept(p) {
			return p, true
		}
	}
	return Parameter{}, false
}

// FindByName .
func (ps Parameters) FindByName(name string) (Parameter, bool) {
	return ps.Find(ParameterNameEquals(name))
}

//Filter
func (ps Parameters) Filter(accept func(p Parameter) bool) Parameters {
	var res Parameters
	for _, p := range ps {
		if accept(p) {
			res = append(res, p)
		}
	}
	return res
}

func ParameterNameEquals(name string) func(Parameter) bool {
	return func(p Parameter) bool {
		return name == p.Name
	}
}

func ParameterTypeEquals(typ string) func(Parameter) bool {
	return func(p Parameter) bool {
		return typ == p.Type
	}
}

// SelectOption .
type SelectOption struct {
	Label       string        `json:"label,omitempty"`
	Description string        `json:"description,omitempty"`
	Value       interface{}   `json:"value,omitempty"`
	Options     SelectOptions `json:"options,omitempty"`
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
