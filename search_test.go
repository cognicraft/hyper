package hyper_test

import (
	"reflect"
	"testing"

	"github.com/cognicraft/hyper"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name        string
		root        hyper.Item
		id          string
		expectItem  hyper.Item
		expectFound bool
	}{
		{
			name:        "empty",
			root:        hyper.Item{},
			id:          "1",
			expectItem:  hyper.Item{},
			expectFound: false,
		},
		{
			name: "root",
			root: hyper.Item{
				ID: "1",
			},
			id: "1",
			expectItem: hyper.Item{
				ID: "1",
			},
			expectFound: true,
		},
		{
			name: "sub-item",
			root: hyper.Item{
				ID: "1",
				Items: hyper.Items{
					{
						ID: "1.1",
					},
					{
						ID: "1.2",
					},
				},
			},
			id: "1.2",
			expectItem: hyper.Item{
				ID: "1.2",
			},
			expectFound: true,
		},
		{
			name: "sub-item",
			root: hyper.Item{
				ID: "1",
				Items: hyper.Items{
					{
						ID: "1.1",
					},
					{
						ID: "1.2",
					},
				},
			},
			id:          "1.3",
			expectItem:  hyper.Item{},
			expectFound: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item, found := hyper.Search(test.root, test.id)
			if test.expectFound != found {
				t.Errorf("want: %v, got: %v", test.expectFound, found)
			}
			if !reflect.DeepEqual(test.expectItem, item) {
				t.Errorf("want: %#v, got: %#v", test.expectItem, item)
			}
		})
	}
}
