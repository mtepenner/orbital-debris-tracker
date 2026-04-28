package spatial

import "testing"

func TestQueryRadiusFindsNearbyPoint(t *testing.T) {
	tree := NewOctree(Box{CenterX: 0, CenterY: 0, CenterZ: 0, Half: 100}, 2)
	tree.Insert(Point{ID: "a", X: 1, Y: 2, Z: 3})
	tree.Insert(Point{ID: "b", X: 60, Y: 60, Z: 60})

	var found []Point
	tree.QueryRadius(Point{X: 0, Y: 0, Z: 0}, 10, &found)
	if len(found) != 1 || found[0].ID != "a" {
		t.Fatalf("unexpected query results: %+v", found)
	}
}
