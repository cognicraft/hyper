package hyper

// Action .
type Action struct {
	Label       string     `json:"label,omitempty"`
	Description string     `json:"description,omitempty"`
	Render      string     `json:"render,omitempty"`
	Rel         string     `json:"rel"`
	Href        string     `json:"href,omitempty"`
	Encoding    string     `json:"encoding,omitempty"`
	Method      string     `json:"method,omitempty"`
	Template    string     `json:"template,omitempty"`
	Parameters  Parameters `json:"parameters,omitempty"`
	Context     string     `json:"context,omitempty"`
	OK          string     `json:"ok,omitempty"`
	Cancel      string     `json:"cancel,omitempty"`
}

// Actions .
type Actions []Action

func (as Actions) Len() int {
	return len(as)
}

func (as Actions) Less(i, j int) bool {
	return as[i].Label < as[j].Label
}

func (as Actions) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

// Find .
func (as Actions) Find(accept func(Action) bool) (Action, bool) {
	for _, a := range as {
		if accept(a) {
			return a, true
		}
	}
	return Action{}, false
}

// FindByRel .
func (as Actions) FindByRel(rel string) (Action, bool) {
	return as.Find(ActionRelEquals(rel))
}

// Filter .
func (as Actions) Filter(accept func(Action) bool) Actions {
	var res Actions
	for _, a := range as {
		if accept(a) {
			res = append(res, a)
		}
	}
	return res
}

// FilterByRel .
func (as Actions) FilterByRel(rel string) Actions {
	return as.Filter(ActionRelEquals(rel))
}

func ActionRelEquals(rel string) func(Action) bool {
	return func(a Action) bool {
		return rel == a.Rel
	}
}

const (
	MethodPOST   = "POST"
	MethodPATCH  = "PATCH"
	MethodDELETE = "DELETE"
)
