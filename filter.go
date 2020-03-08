package hyper

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseFilter(url *url.URL) (Filter, error) {
	rawFCs, ok := url.Query()["filter"]
	if !ok {
		return Filter{}, nil
	}
	var fcs []FilterComponent
	for _, rawFC := range rawFCs {
		fc, err := ParseFilterComponent(rawFC)
		if err != nil {
			return nil, err
		}
		fcs = append(fcs, fc)
	}
	return fcs, nil
}

type Filter []FilterComponent

func (f Filter) IsZero() bool {
	return len(f) == 0
}

func (f Filter) HasComponent(name string) bool {
	for _, fc := range f {
		if fc.Name == name {
			return true
		}
	}
	return false
}

func (f Filter) FindOne(name string) (FilterComponent, bool) {
	for _, fc := range f {
		if fc.Name == name {
			return fc, true
		}
	}
	return FilterComponent{}, false
}

func (f Filter) RemoveAll(name string) Filter {
	var out Filter
	for _, fc := range f {
		if fc.Name != name {
			out = append(out, fc)
		}
	}
	return out
}

func (f Filter) Filter(accept func(c FilterComponent) bool) Filter {
	var res Filter
	for _, fc := range f {
		if accept(fc) {
			res = append(res, fc)
		}
	}
	return res
}

func Named(name string) func(FilterComponent) bool {
	return func(c FilterComponent) bool {
		return name == c.Name
	}
}

func Not(accept func(FilterComponent) bool) func(FilterComponent) bool {
	return func(c FilterComponent) bool {
		return !accept(c)
	}
}

func ParseFilterComponent(rawFC string) (FilterComponent, error) {
	parts := strings.Split(rawFC, ",")
	if len(parts) < 3 {
		return FilterComponent{}, fmt.Errorf("invalid filter component: %s", parts)
	}
	name := parts[0]
	operator := FilterOperator(parts[1])
	values := parts[2:]
	switch operator {
	case FilterOperatorIn, FilterOperatorNotIn:
		return FilterComponent{
			Name:     name,
			Operator: operator,
			Value:    values,
		}, nil
	case FilterOperatorBetween, FilterOperatorNotBetween:
		if len(values) < 2 {
			return FilterComponent{}, fmt.Errorf("invalid size of filter-between values (must >= 2): %d", len(values))
		}
		return FilterComponent{
			Name:     name,
			Operator: operator,
			Value:    values[:2],
		}, nil
	default:
		return FilterComponent{
			Name:     name,
			Operator: operator,
			Value:    values[0],
		}, nil
	}
}

type FilterComponent struct {
	Operator FilterOperator `json:"operator,omitempty"`
	Name     string         `json:"name"`
	Value    interface{}    `json:"value,omitempty"`
}

func (fc FilterComponent) ValueStrings() []string {
	ss := []string{}
	switch v := fc.Value.(type) {
	case []string:
		return v
	}
	return ss
}

func (fc FilterComponent) ValueString() string {
	switch v := fc.Value.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func (fc FilterComponent) ValueBool() bool {
	switch v := fc.Value.(type) {
	case bool:
		return v
	case int:
		switch v {
		case 1:
			return true
		default:
			return false
		}

	case string:
		switch v {
		case "true":
			return true
		default:
			return false
		}
	default:
		return false
	}
}

type FilterOperator string

const (
	FilterOperatorEquals              FilterOperator = "eq"
	FilterOperatorNotEquals           FilterOperator = "neq"
	FilterOperatorLessThen            FilterOperator = "lt"
	FilterOperatorGreaterThen         FilterOperator = "gt"
	FilterOperatorLessThenOrEquals    FilterOperator = "leq"
	FilterOperatorGreaterThenOrEquals FilterOperator = "geq"
	FilterOperatorIn                  FilterOperator = "in"
	FilterOperatorNotIn               FilterOperator = "nin"
	FilterOperatorLike                FilterOperator = "like"
	FilterOperatorNotLike             FilterOperator = "nlike"
	FilterOperatorBetween             FilterOperator = "bet"
	FilterOperatorNotBetween          FilterOperator = "nbet"
)

func MakeFilterLink(fc FilterConfiguration, template string, currentFilter Filter, placeholder string) Link {
	return Link{
		Rel:      "filter",
		Template: template,
		Parameters: []Parameter{
			{
				Name:        "filter",
				Type:        "filter",
				Components:  fc,
				Value:       currentFilter,
				Placeholder: placeholder,
			},
		},
	}
}

type FilterConfiguration []FilterComponentConfiguration

func (c FilterConfiguration) Len() int {
	return len(c)
}

func (c FilterConfiguration) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c FilterConfiguration) Less(i, j int) bool {
	return c[i].Label < c[j].Label
}

type FilterComponentConfiguration struct {
	Label       string                        `json:"label,omitempty"`
	Description string                        `json:"description,omitempty"`
	Name        string                        `json:"name"`
	Operators   []FilterOperatorConfiguration `json:"operators,omitempty"`
	Type        string                        `json:"type,omitempty"`
	Options     SelectOptions                 `json:"options,omitempty"`
	Related     string                        `json:"related,omitempty"`
	Placeholder string                        `json:"placeholder,omitempty"`
	Multiple    bool                          `json:"multiple,omitempty"`
}

type FilterOperatorConfiguration struct {
	Label       string         `json:"label,omitempty"`
	Description string         `json:"description,omitempty"`
	Operator    FilterOperator `json:"operator"`
	Infix       string         `json:"infix,omitempty"`
}

var configs = map[FilterOperator]FilterOperatorConfiguration{
	FilterOperatorEquals:              {Label: "=", Operator: FilterOperatorEquals},
	FilterOperatorNotEquals:           {Label: "!=", Operator: FilterOperatorNotEquals},
	FilterOperatorLessThen:            {Label: "<", Operator: FilterOperatorLessThen},
	FilterOperatorGreaterThen:         {Label: ">", Operator: FilterOperatorGreaterThen},
	FilterOperatorLessThenOrEquals:    {Label: "<=", Operator: FilterOperatorLessThenOrEquals},
	FilterOperatorGreaterThenOrEquals: {Label: ">=", Operator: FilterOperatorGreaterThenOrEquals},
	FilterOperatorBetween:             {Label: "between", Operator: FilterOperatorBetween, Infix: "and"},
	FilterOperatorNotBetween:          {Label: "not between", Operator: FilterOperatorNotBetween, Infix: "and"},
	FilterOperatorIn:                  {Label: "in", Operator: FilterOperatorIn},
	FilterOperatorNotIn:               {Label: "not in", Operator: FilterOperatorNotIn},
	FilterOperatorLike:                {Label: "like", Operator: FilterOperatorLike},
	FilterOperatorNotLike:             {Label: "not like", Operator: FilterOperatorNotLike},
}

func FilterOperatorConfigurationFor(op FilterOperator) FilterOperatorConfiguration {
	c := configs[op]
	return c
}

func FilterOperatorConfigurationsForSelect() []FilterOperatorConfiguration {
	return FilterOperatorConfigurationsForBaseType("select")
}

func FilterOperatorConfigurationsForText() []FilterOperatorConfiguration {
	return FilterOperatorConfigurationsForBaseType(TypeText)
}

func FilterOperatorConfigurationsForTextReduced() []FilterOperatorConfiguration {
	return []FilterOperatorConfiguration{
		// FilterOperatorConfigurationFor(FilterOperatorEquals),
		// FilterOperatorConfigurationFor(FilterOperatorNotEquals),
		// FilterOperatorConfigurationFor(FilterOperatorLessThen),
		// FilterOperatorConfigurationFor(FilterOperatorGreaterThen),
		// FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
		// FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
		// FilterOperatorConfigurationFor(FilterOperatorBetween),
		// FilterOperatorConfigurationFor(FilterOperatorNotBetween),
		FilterOperatorConfigurationFor(FilterOperatorLike),
		FilterOperatorConfigurationFor(FilterOperatorNotLike),
		FilterOperatorConfigurationFor(FilterOperatorIn),
		FilterOperatorConfigurationFor(FilterOperatorNotIn),
	}
}

func FilterOperatorConfigurationsForInteger() []FilterOperatorConfiguration {
	return FilterOperatorConfigurationsForBaseType("integer")
}

func FilterOperatorConfigurationsForIntegerReduced() []FilterOperatorConfiguration {
	return []FilterOperatorConfiguration{
		FilterOperatorConfigurationFor(FilterOperatorEquals),
		FilterOperatorConfigurationFor(FilterOperatorNotEquals),
		FilterOperatorConfigurationFor(FilterOperatorLessThen),
		FilterOperatorConfigurationFor(FilterOperatorGreaterThen),
		FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
		FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
		// FilterOperatorConfigurationFor(FilterOperatorBetween),
		// FilterOperatorConfigurationFor(FilterOperatorNotBetween),
		// FilterOperatorConfigurationFor(FilterOperatorIn),
		// FilterOperatorConfigurationFor(FilterOperatorNotIn),
	}
}

func FilterOperatorConfigurationsForNumber() []FilterOperatorConfiguration {
	return FilterOperatorConfigurationsForBaseType(TypeNumber)
}

func FilterOperatorConfigurationsForNumberReduced() []FilterOperatorConfiguration {
	return []FilterOperatorConfiguration{
		FilterOperatorConfigurationFor(FilterOperatorEquals),
		FilterOperatorConfigurationFor(FilterOperatorNotEquals),
		FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
		FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
	}
}

func FilterOperatorConfigurationsForDate() []FilterOperatorConfiguration {
	return FilterOperatorConfigurationsForBaseType(TypeDate)
}

func FilterOperatorConfigurationsForDatetime() []FilterOperatorConfiguration {
	return FilterOperatorConfigurationsForBaseType(TypeDatetime)
}

func FilterOperatorConfigurationsForDatetimeReduced() []FilterOperatorConfiguration {
	return []FilterOperatorConfiguration{
		FilterOperatorConfigurationFor(FilterOperatorEquals),
		FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
		FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
	}
}

func FilterOperatorConfigurationsForBaseType(t string) []FilterOperatorConfiguration {

	switch t {
	case TypeText:
		return []FilterOperatorConfiguration{
			FilterOperatorConfigurationFor(FilterOperatorEquals),
			FilterOperatorConfigurationFor(FilterOperatorNotEquals),
			FilterOperatorConfigurationFor(FilterOperatorLessThen),
			FilterOperatorConfigurationFor(FilterOperatorGreaterThen),
			FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
			FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
			FilterOperatorConfigurationFor(FilterOperatorBetween),
			FilterOperatorConfigurationFor(FilterOperatorNotBetween),
			FilterOperatorConfigurationFor(FilterOperatorIn),
			FilterOperatorConfigurationFor(FilterOperatorNotIn),
			FilterOperatorConfigurationFor(FilterOperatorLike),
			FilterOperatorConfigurationFor(FilterOperatorNotLike),
		}
	case "integer":
		return []FilterOperatorConfiguration{
			FilterOperatorConfigurationFor(FilterOperatorEquals),
			FilterOperatorConfigurationFor(FilterOperatorNotEquals),
			FilterOperatorConfigurationFor(FilterOperatorLessThen),
			FilterOperatorConfigurationFor(FilterOperatorGreaterThen),
			FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
			FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
			FilterOperatorConfigurationFor(FilterOperatorBetween),
			FilterOperatorConfigurationFor(FilterOperatorNotBetween),
			FilterOperatorConfigurationFor(FilterOperatorIn),
			FilterOperatorConfigurationFor(FilterOperatorNotIn),
		}
	case TypeNumber, TypeDate, TypeDatetime:
		return []FilterOperatorConfiguration{
			FilterOperatorConfigurationFor(FilterOperatorEquals),
			FilterOperatorConfigurationFor(FilterOperatorNotEquals),
			FilterOperatorConfigurationFor(FilterOperatorLessThen),
			FilterOperatorConfigurationFor(FilterOperatorGreaterThen),
			FilterOperatorConfigurationFor(FilterOperatorLessThenOrEquals),
			FilterOperatorConfigurationFor(FilterOperatorGreaterThenOrEquals),
			FilterOperatorConfigurationFor(FilterOperatorBetween),
			FilterOperatorConfigurationFor(FilterOperatorNotBetween),
		}
	case "select":
		return []FilterOperatorConfiguration{
			FilterOperatorConfigurationFor(FilterOperatorEquals),
			FilterOperatorConfigurationFor(FilterOperatorNotEquals),
			FilterOperatorConfigurationFor(FilterOperatorIn),
			FilterOperatorConfigurationFor(FilterOperatorNotIn),
		}
	}
	return []FilterOperatorConfiguration{
		FilterOperatorConfigurationFor(FilterOperatorEquals),
		FilterOperatorConfigurationFor(FilterOperatorNotEquals),
	}
}
