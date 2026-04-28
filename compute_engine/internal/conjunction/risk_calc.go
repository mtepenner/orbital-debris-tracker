package conjunction

import (
	"math"
	"sort"
	"time"

	"github.com/mtepenner/orbital-debris-tracker/compute_engine/internal/sgp4"
	"github.com/mtepenner/orbital-debris-tracker/compute_engine/internal/spatial"
)

type Conjunction struct {
	PrimaryID            string  `json:"primary_id"`
	SecondaryID          string  `json:"secondary_id"`
	MissDistanceKm       float64 `json:"miss_distance_km"`
	ProbabilityCollision float64 `json:"probability_collision"`
	TCA                  string  `json:"tca"`
}

func Evaluate(states []sgp4.State, thresholdKm float64, now time.Time) []Conjunction {
	tree := spatial.NewOctree(spatial.Box{CenterX: 0, CenterY: 0, CenterZ: 0, Half: 60000}, 8)
	points := make(map[string]sgp4.State, len(states))
	for _, state := range states {
		point := spatial.Point{ID: state.ObjectID, X: state.XKm, Y: state.YKm, Z: state.ZKm}
		tree.Insert(point)
		points[state.ObjectID] = state
	}

	results := make([]Conjunction, 0)
	seen := map[string]struct{}{}
	for _, state := range states {
		var nearby []spatial.Point
		tree.QueryRadius(spatial.Point{ID: state.ObjectID, X: state.XKm, Y: state.YKm, Z: state.ZKm}, thresholdKm, &nearby)
		for _, candidate := range nearby {
			if candidate.ID == state.ObjectID {
				continue
			}
			pair := orderedPair(state.ObjectID, candidate.ID)
			if _, exists := seen[pair]; exists {
				continue
			}
			seen[pair] = struct{}{}

			other := points[candidate.ID]
			miss := missDistance(state, other)
			results = append(results, Conjunction{
				PrimaryID:            minID(state.ObjectID, other.ObjectID),
				SecondaryID:          maxID(state.ObjectID, other.ObjectID),
				MissDistanceKm:       miss,
				ProbabilityCollision: collisionProbability(miss, thresholdKm),
				TCA:                  now.Add(time.Duration(miss/state.SpeedKmS) * time.Second).Format(time.RFC3339),
			})
		}
	}
	sort.Slice(results, func(i int, j int) bool {
		return results[i].MissDistanceKm < results[j].MissDistanceKm
	})
	return results
}

func missDistance(left sgp4.State, right sgp4.State) float64 {
	dx := left.XKm - right.XKm
	dy := left.YKm - right.YKm
	dz := left.ZKm - right.ZKm
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func collisionProbability(missDistanceKm float64, thresholdKm float64) float64 {
	if thresholdKm == 0 {
		return 0
	}
	return math.Exp(-missDistanceKm/thresholdKm) * 0.01
}

func orderedPair(left string, right string) string {
	if left < right {
		return left + ":" + right
	}
	return right + ":" + left
}

func minID(left string, right string) string {
	if left < right {
		return left
	}
	return right
}

func maxID(left string, right string) string {
	if left > right {
		return left
	}
	return right
}
