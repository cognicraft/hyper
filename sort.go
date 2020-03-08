package hyper

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseSort(url *url.URL) (Sort, error) {
	rawSCs, ok := url.Query()["sort"]
	if !ok {
		return Sort{}, nil
	}
	var scs []SortComponent
	for _, rawSC := range rawSCs {
		sc, err := ParseSortComponent(rawSC)
		if err != nil {
			return nil, err
		}
		scs = append(scs, sc)
	}
	return scs, nil
}

type Sort []SortComponent

func (s Sort) IsZero() bool {
	return len(s) == 0
}

func (s Sort) FindOne(name string) (SortComponent, bool) {
	for _, sc := range s {
		if sc.Name == name {
			return sc, true
		}
	}
	return SortComponent{}, false
}

func ParseSortComponent(rawSC string) (SortComponent, error) {
	sExpr := strings.Split(rawSC, ",")
	if len(sExpr) != 2 {
		return SortComponent{}, fmt.Errorf("invalid sort component: %s", sExpr)
	}
	sc := SortComponent{
		Name: sExpr[0],
	}
	switch sExpr[1] {
	case string(SortOrderAscending):
		sc.Order = SortOrderAscending
	case string(SortOrderDescending):
		sc.Order = SortOrderDescending
	default:
		return SortComponent{}, fmt.Errorf("invalid sort order: %s", sExpr[0])
	}
	return sc, nil
}

type SortComponent struct {
	Order SortOrder `json:"order"`
	Name  string    `json:"name"`
}

type SortOrder string

const (
	SortOrderAscending  SortOrder = "ASC"
	SortOrderDescending SortOrder = "DESC"
)

type SortConfiguration []SortComponentConfiguration

type SortComponentConfiguration struct {
	Label       string                   `json:"label,omitempty"`
	Description string                   `json:"description,omitempty"`
	Name        string                   `json:"name"`
	Orders      []SortOrderConfiguration `json:"orders,omitempty"`
}

type SortOrderConfiguration struct {
	Label       string    `json:"label,omitempty"`
	Description string    `json:"description,omitempty"`
	Order       SortOrder `json:"order"`
}

func MakeSortOperatorConfigurations() []SortOrderConfiguration {
	return []SortOrderConfiguration{
		{Label: "ascending", Order: SortOrderAscending},
		{Label: "descending", Order: SortOrderDescending},
	}
}
