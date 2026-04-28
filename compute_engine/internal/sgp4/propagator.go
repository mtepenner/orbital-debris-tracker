package sgp4

import (
	"math"
	"time"
)

const earthMuKm3S2 = 398600.4418

type CatalogObject struct {
	ObjectID            string  `json:"object_id"`
	Name                string  `json:"name"`
	InclinationDeg      float64 `json:"inclination_deg"`
	RAANDeg             float64 `json:"raan_deg"`
	Eccentricity        float64 `json:"eccentricity"`
	MeanAnomalyDeg      float64 `json:"mean_anomaly_deg"`
	MeanMotionRevPerDay float64 `json:"mean_motion_rev_per_day"`
}

type State struct {
	ObjectID string  `json:"object_id"`
	Name     string  `json:"name"`
	XKm      float64 `json:"x_km"`
	YKm      float64 `json:"y_km"`
	ZKm      float64 `json:"z_km"`
	SpeedKmS float64 `json:"speed_km_s"`
}

func Propagate(object CatalogObject, at time.Time) State {
	meanMotionRadS := object.MeanMotionRevPerDay * 2 * math.Pi / 86400.0
	semiMajorAxisKm := math.Cbrt(earthMuKm3S2 / (meanMotionRadS * meanMotionRadS))
	meanAnomaly := degToRad(object.MeanAnomalyDeg) + meanMotionRadS*float64(at.Unix())
	radiusKm := semiMajorAxisKm * (1 - object.Eccentricity*0.35)
	inclination := degToRad(object.InclinationDeg)
	raan := degToRad(object.RAANDeg)

	xOrbital := radiusKm * math.Cos(meanAnomaly)
	yOrbital := radiusKm * math.Sin(meanAnomaly)

	xEquatorial := xOrbital*math.Cos(raan) - yOrbital*math.Sin(raan)*math.Cos(inclination)
	yEquatorial := xOrbital*math.Sin(raan) + yOrbital*math.Cos(raan)*math.Cos(inclination)
	zEquatorial := yOrbital * math.Sin(inclination)

	return State{
		ObjectID: object.ObjectID,
		Name:     object.Name,
		XKm:      xEquatorial,
		YKm:      yEquatorial,
		ZKm:      zEquatorial,
		SpeedKmS: meanMotionRadS * radiusKm,
	}
}

func SampleObjects() []CatalogObject {
	return []CatalogObject{
		{ObjectID: "25544", Name: "ISS (ZARYA)", InclinationDeg: 51.64, RAANDeg: 62.84, Eccentricity: 0.0006703, MeanAnomalyDeg: 39.85, MeanMotionRevPerDay: 15.50003235},
		{ObjectID: "29716", Name: "FENGYUN 1C DEB", InclinationDeg: 98.68, RAANDeg: 124.41, Eccentricity: 0.0104321, MeanAnomalyDeg: 164.70, MeanMotionRevPerDay: 14.17654241},
		{ObjectID: "54237", Name: "COSMOS 1408 DEB", InclinationDeg: 82.55, RAANDeg: 201.12, Eccentricity: 0.0023411, MeanAnomalyDeg: 335.76, MeanMotionRevPerDay: 15.21541021},
	}
}

func degToRad(value float64) float64 {
	return value * math.Pi / 180.0
}
