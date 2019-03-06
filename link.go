package hyper

// Link .
type Link struct {
	Label          string     `json:"label,omitempty"`
	Description    string     `json:"description,omitempty"`
	Render         string     `json:"render,omitempty"`
	Rel            string     `json:"rel"`
	Href           string     `json:"href,omitempty"`
	Type           string     `json:"type,omitempty"`
	Language       string     `json:"language,omitempty"`
	Template       string     `json:"template,omitempty"`
	Parameters     Parameters `json:"parameters,omitempty"`
	Context        string     `json:"context,omitempty"`
	Accept         string     `json:"accept,omitempty"`
	AcceptLanguage string     `json:"accept-language,omitempty"`
}

// Links .
type Links []Link

// Find .
func (ls Links) Find(f func(Link) bool) (Link, bool) {
	for _, l := range ls {
		if f(l) {
			return l, true
		}
	}
	return Link{}, false
}

// FindByRel .
func (ls Links) FindByRel(rel string) (Link, bool) {
	return ls.Find(LinkRelEquals(rel))
}

// Filter .
func (ls Links) Filter(f func(Link) bool) Links {
	var filtered Links
	for _, l := range ls {
		if f(l) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

//FilterByRel .
func (ls Links) FilterByRel(rel string) Links {
	return ls.Filter(LinkRelEquals(rel))
}

func LinkRelEquals(rel string) func(Link) bool {
	return func(l Link) bool {
		return rel == l.Rel
	}
}
