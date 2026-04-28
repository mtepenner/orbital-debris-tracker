package parser

import "testing"

func TestParseCatalog(t *testing.T) {
	raw := `OBJECT ONE
1 25544U 98067A   26118.50109859  .00008064  00000+0  14714-3 0  9991
2 25544  51.6400  62.8406 0006703  74.0625  39.8567 15.50003235 42167`

	objects, err := ParseCatalog(raw)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(objects) != 1 {
		t.Fatalf("expected one object, got %d", len(objects))
	}
	if objects[0].Name != "OBJECT ONE" || objects[0].ObjectID != "25544" {
		t.Fatalf("unexpected object parsed: %+v", objects[0])
	}
}
