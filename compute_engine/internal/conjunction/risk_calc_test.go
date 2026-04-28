package conjunction

import (
	"testing"
	"time"

	"github.com/mtepenner/orbital-debris-tracker/compute_engine/internal/sgp4"
)

func TestEvaluateFindsClosePair(t *testing.T) {
	states := []sgp4.State{
		{ObjectID: "a", XKm: 0, YKm: 0, ZKm: 0, SpeedKmS: 7.5},
		{ObjectID: "b", XKm: 1, YKm: 1, ZKm: 1, SpeedKmS: 7.4},
		{ObjectID: "c", XKm: 100, YKm: 100, ZKm: 100, SpeedKmS: 7.3},
	}

	conjunctions := Evaluate(states, 5, time.Unix(0, 0))
	if len(conjunctions) != 1 {
		t.Fatalf("expected one conjunction, got %d", len(conjunctions))
	}
	if conjunctions[0].PrimaryID != "a" || conjunctions[0].SecondaryID != "b" {
		t.Fatalf("unexpected conjunction pair: %+v", conjunctions[0])
	}
}
