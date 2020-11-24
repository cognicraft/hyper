package hyper

import (
	"encoding/json"
)

// Item has properties, links, actions, (sub-)items and errors.
type Item struct {
	Label       string          `json:"label,omitempty"`
	Description string          `json:"description,omitempty"`
	Render      []string        `json:"render,omitempty"`
	Rel         string          `json:"rel,omitempty"`
	ID          string          `json:"id,omitempty"`
	Type        string          `json:"type,omitempty"`
	Properties  Properties      `json:"properties,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
	Links       Links           `json:"links,omitempty"`
	Actions     Actions         `json:"actions,omitempty"`
	Items       Items           `json:"items,omitempty"`
	Errors      Errors          `json:"errors,omitempty"`
}

// AddProperty add a Property to this Item
func (i *Item) AddProperty(p Property) {
	i.Properties = append(i.Properties, p)
}

// AddProperties adds many Properties to this Item
func (i *Item) AddProperties(ps Properties) {
	i.Properties = append(i.Properties, ps...)
}

// AddItem adds a (sub-)Item to this Item
func (i *Item) AddItem(sub Item) {
	i.Items = append(i.Items, sub)
}

// AddItems adds many (sub-)Items to this Item
func (i *Item) AddItems(subs []Item) {
	i.Items = append(i.Items, subs...)
}

// AddLink adds a link to this Item
func (i *Item) AddLink(l Link) {
	i.Links = append(i.Links, l)
}

// AddLinks adds many links to this Item
func (i *Item) AddLinks(ls Links) {
	i.Links = append(i.Links, ls...)
}

// AddAction adds an Action to this Item
func (i *Item) AddAction(a Action) {
	i.Actions = append(i.Actions, a)
}

// AddActions adds many Actions to this Item
func (i *Item) AddActions(as Actions) {
	i.Actions = append(i.Actions, as...)
}

func (i *Item) EncodeData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	i.Data = json.RawMessage(data)
	return nil
}

func (i *Item) DecodeData(v interface{}) error {
	return json.Unmarshal(i.Data, v)
}

// Items represents a collection of Item
type Items []Item

// Find returns an Item that satisfies the specification
func (is Items) Find(accept func(Item) bool) (Item, bool) {
	for _, i := range is {
		if accept(i) {
			return i, true
		}
	}
	return Item{}, false
}

// FindByID returns an Item that has a specific id
func (is Items) FindByID(id string) (Item, bool) {
	return is.Find(ItemIDEquals(id))
}

// FindByRel returns an Item that has a specific rel
func (is Items) FindByRel(rel string) (Item, bool) {
	return is.Find(ItemRelEquals(rel))
}

// KeyBy calculates a map keyed by the result of the extractKey funktion.
func (is Items) KeyBy(extractKey func(Item) string) map[string]Item {
	if len(is) == 0 {
		return nil
	}
	m := map[string]Item{}
	for _, i := range is {
		key := extractKey(i)
		m[key] = i
	}
	return m
}

// KeyByID returns a map of Items keyed by the Item ids
func (is Items) KeyByID() map[string]Item {
	return is.KeyBy(func(i Item) string {
		return i.ID
	})
}

// KeyByRel returns a map of Items keyed by the Item rel
func (is Items) KeyByRel() map[string]Item {
	return is.KeyBy(func(i Item) string {
		return i.Rel
	})
}

// Filter returns a collection of Items that conform to the profided specification
func (is Items) Filter(accept func(Item) bool) Items {
	var res Items
	for _, i := range is {
		if accept(i) {
			res = append(res, i)
		}
	}
	return res
}

// ItemIDEquals is used to Filter a collection of Items by id
func ItemIDEquals(id string) func(Item) bool {
	return func(i Item) bool {
		return id == i.ID
	}
}

// ItemRelEquals is used to Filter a collection of Items by rel
func ItemRelEquals(rel string) func(Item) bool {
	return func(i Item) bool {
		return rel == i.Rel
	}
}
