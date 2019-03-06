package hyper

// Search performs a DFS with the goal to find an item by the specified id
func Search(root Item, id string) (Item, bool) {
	frontier := []Item{root}
	for len(frontier) > 0 {
		var next Item
		next, frontier = frontier[0], frontier[1:]
		if next.ID == id {
			return next, true
		}
		frontier = append(next.Items, frontier...)
	}
	return Item{}, false
}
