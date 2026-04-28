package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CatalogObject struct {
	ObjectID            string  `json:"object_id"`
	Name                string  `json:"name"`
	Line1               string  `json:"line1"`
	Line2               string  `json:"line2"`
	InclinationDeg      float64 `json:"inclination_deg"`
	RAANDeg             float64 `json:"raan_deg"`
	Eccentricity        float64 `json:"eccentricity"`
	MeanAnomalyDeg      float64 `json:"mean_anomaly_deg"`
	MeanMotionRevPerDay float64 `json:"mean_motion_rev_per_day"`
}

func ParseCatalog(raw string) ([]CatalogObject, error) {
	lines := splitLines(raw)
	if len(lines)%3 != 0 {
		return nil, errors.New("catalog must be divisible into 3-line entries")
	}

	objects := make([]CatalogObject, 0, len(lines)/3)
	for index := 0; index < len(lines); index += 3 {
		entry, err := parseEntry(lines[index], lines[index+1], lines[index+2])
		if err != nil {
			return nil, fmt.Errorf("parse entry %d: %w", index/3, err)
		}
		objects = append(objects, entry)
	}
	return objects, nil
}

func parseEntry(title string, line1 string, line2 string) (CatalogObject, error) {
	line1Fields := strings.Fields(line1)
	line2Fields := strings.Fields(line2)
	if len(line1Fields) < 2 || len(line2Fields) < 8 {
		return CatalogObject{}, errors.New("tle entry missing expected fields")
	}

	parse := func(value string) (float64, error) { return strconv.ParseFloat(value, 64) }
	inc, err := parse(line2Fields[2])
	if err != nil {
		return CatalogObject{}, err
	}
	raan, err := parse(line2Fields[3])
	if err != nil {
		return CatalogObject{}, err
	}
	ecc, err := parse("0." + line2Fields[4])
	if err != nil {
		return CatalogObject{}, err
	}
	meanAnomaly, err := parse(line2Fields[6])
	if err != nil {
		return CatalogObject{}, err
	}
	meanMotion, err := parse(line2Fields[7])
	if err != nil {
		return CatalogObject{}, err
	}

	return CatalogObject{
		ObjectID:            strings.TrimSuffix(line1Fields[1], "U"),
		Name:                title,
		Line1:               line1,
		Line2:               line2,
		InclinationDeg:      inc,
		RAANDeg:             raan,
		Eccentricity:        ecc,
		MeanAnomalyDeg:      meanAnomaly,
		MeanMotionRevPerDay: meanMotion,
	}, nil
}

func splitLines(raw string) []string {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	return filtered
}
