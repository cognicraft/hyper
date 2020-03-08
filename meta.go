package hyper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/cognicraft/uri"
)

type Meta struct {
	Filter Filter `json:"filter,omitemty"`
	Sort   Sort   `json:"sort,omitempty"`
	Search string `json:"search,omitempty"`
	Skip   uint64 `json:"skip,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
	After  string `json:"after,omitempty"`
	Before string `json:"before,omitempty"`
}

func (m Meta) Clone() Meta {
	clone := Meta{
		Search: m.Search,
		Skip:   m.Skip,
		Limit:  m.Limit,
		After:  m.After,
		Before: m.Before,
	}
	clone.Filter = make(Filter, len(m.Filter))
	copy(clone.Filter, m.Filter)

	clone.Sort = make(Sort, len(m.Sort))
	copy(clone.Sort, m.Sort)
	return clone
}

func (m Meta) MinLimit(limit uint64) Meta {
	if m.Limit < limit {
		m.Limit = limit
	}
	return m
}

func (m Meta) DefaultSort(sort Sort) Meta {
	if m.Sort.IsZero() {
		m.Sort = sort
	}
	return m
}

func (m Meta) IsZero() bool {
	return len(m.Filter) == 0 &&
		len(m.Sort) == 0 &&
		m.Search == "" &&
		m.Skip == 0 &&
		m.Limit == 0 &&
		m.After == "" &&
		m.Before == ""
}

func (m Meta) MarshalJSON() ([]byte, error) {
	res := map[string]interface{}{}
	if !m.Filter.IsZero() {
		res["filter"] = m.Filter
	}
	if !m.Sort.IsZero() {
		res["sort"] = m.Sort
	}
	if m.Search != "" {
		res["search"] = m.Search
	}
	if m.Skip > 0 {
		res["skip"] = m.Skip
	}
	if m.Limit > 0 {
		res["limit"] = m.Limit
	}
	if m.After != "" {
		res["after"] = m.After
	}
	if m.Before != "" {
		res["before"] = m.Before
	}
	return json.Marshal(res)
}

func (m Meta) Next() Meta {
	m.Skip += m.Limit
	return m
}

func (m Meta) Previous() Meta {
	if m.Skip < m.Limit {
		m.Skip = 0
		return m
	}
	m.Skip -= m.Limit
	return m
}

func (m Meta) WithAfter(id string) Meta {
	m.After = id
	return m
}

func (m Meta) WithBefore(id string) Meta {
	m.Before = id
	return m
}

func (m Meta) Query() string {
	qt, err := uri.Parse("{?filter*,sort*,search,skip,limit,after,before}")
	if err != nil {
		log.Printf("uri parse: %s", err)
	}

	vs := map[string]interface{}{}
	if !m.Filter.IsZero() {
		vs["filter"] = m.currentFilter()
	}
	if !m.Sort.IsZero() {
		vs["sort"] = m.currentSort()
	}
	if m.Search != "" {
		vs["search"] = m.Search
	}
	if m.Skip > 0 {
		vs["skip"] = m.Skip
	}
	if m.Limit > 0 {
		vs["limit"] = m.Limit
	}
	if m.After != "" {
		vs["after"] = m.After
	}
	if m.Before != "" {
		vs["before"] = m.Before
	}

	q, err := qt.Expand(vs)
	if err != nil {
		log.Printf("uri expand: %s", err)
	}
	return q
}

func (m Meta) FilterTemplate() string {
	if m.Sort.IsZero() && m.Search == "" && m.Limit == 0 {
		return "{?filter*}"
	}
	qt, err := uri.Parse("{?sort*,search,limit}")
	if err != nil {
		log.Printf("uri parse: %s", err)
	}
	vs := map[string]interface{}{}
	if !m.Sort.IsZero() {
		vs["sort"] = m.currentSort()
	}
	if m.Search != "" {
		vs["search"] = m.Search
	}
	if m.Limit > 0 {
		vs["limit"] = m.Limit
	}
	q, err := qt.Expand(vs)
	if err != nil {
		log.Printf("uri expand: %s", err)
	}
	return q + "{&filter*}"
}

func (m Meta) SearchTemplate() string {
	if m.Filter.IsZero() && m.Sort.IsZero() && m.Limit == 0 {
		return "{?search}"
	}
	qt, err := uri.Parse("{?filter*,sort*,limit}")
	if err != nil {
		log.Printf("uri parse: %s", err)
	}
	vs := map[string]interface{}{}
	if !m.Sort.IsZero() {
		vs["sort"] = m.currentSort()
	}
	if !m.Filter.IsZero() {
		vs["filter"] = m.currentFilter()
	}
	if m.Limit > 0 {
		vs["limit"] = m.Limit
	}
	q, err := qt.Expand(vs)
	if err != nil {
		log.Printf("uri expand: %s", err)
	}
	return q + "{&search}"
}

func (m Meta) SortTemplate() string {
	if m.Filter.IsZero() && m.Search == "" && m.Limit == 0 {
		return "{?sort*}"
	}
	qt, err := uri.Parse("{?filter*,search,limit}")
	if err != nil {
		log.Printf("uri parse: %s", err)
	}
	vs := map[string]interface{}{}
	if !m.Filter.IsZero() {
		vs["filter"] = m.currentFilter()
	}
	if m.Search != "" {
		vs["search"] = m.Search
	}
	if m.Limit > 0 {
		vs["limit"] = m.Limit
	}
	q, err := qt.Expand(vs)
	if err != nil {
		log.Printf("uri expand: %s", err)
	}
	return q + "{&sort*}"
}

func (m Meta) currentFilter() []interface{} {
	fcs := []interface{}{}
	for _, fc := range m.Filter {
		switch fv := fc.Value.(type) {
		case []string:
			fcs = append(fcs, fmt.Sprintf("%s,%s,%s", fc.Name, fc.Operator, strings.Join(fv, ",")))
		case []interface{}:
			sv := make([]string, len(fv))
			for i, v := range fv {
				sv[i] = fmt.Sprintf("%v", v)
			}
			fcs = append(fcs, fmt.Sprintf("%s,%s,%s", fc.Name, fc.Operator, strings.Join(sv, ",")))
		default:
			fcs = append(fcs, fmt.Sprintf("%s,%s,%v", fc.Name, fc.Operator, fc.Value))
		}
	}
	return fcs
}

func (m Meta) currentSort() []interface{} {
	scs := []interface{}{}
	for _, sc := range m.Sort {
		scs = append(scs, fmt.Sprintf("%s,%s", sc.Name, sc.Order))
	}
	return scs
}

func Parse(url *url.URL) (Meta, error) {
	filter, err := ParseFilter(url)
	if err != nil {
		return Meta{}, fmt.Errorf("parse filter: %s", err)
	}
	sort, err := ParseSort(url)
	if err != nil {
		return Meta{}, fmt.Errorf("parse sort: %s", err)
	}
	skip, err := ParseSkip(url)
	if err != nil {
		return Meta{}, fmt.Errorf("parse skip: %s", err)
	}
	limit, err := ParseLimit(url)
	if err != nil {
		return Meta{}, fmt.Errorf("parse limit: %s", err)
	}
	return Meta{
		Filter: filter,
		Sort:   sort,
		Search: url.Query().Get("search"),
		Skip:   skip,
		Limit:  limit,
		After:  url.Query().Get("after"),
		Before: url.Query().Get("before"),
	}, nil
}

func ParseSkip(url *url.URL) (uint64, error) {
	if skip, ok := url.Query()["skip"]; ok {
		s, err := strconv.ParseUint(skip[0], 10, 64)
		if err != nil {
			return 0, err
		}
		return s, nil
	}
	return 0, nil
}

func ParseLimit(url *url.URL) (uint64, error) {
	if limit, ok := url.Query()["limit"]; ok {
		l, err := strconv.ParseUint(limit[0], 10, 64)
		if err != nil {
			return 0, err
		}
		return l, nil
	}
	return 0, nil
}
