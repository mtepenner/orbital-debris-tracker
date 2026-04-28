package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/mtepenner/orbital-debris-tracker/tle_ingestion/internal/parser"
	"github.com/mtepenner/orbital-debris-tracker/tle_ingestion/internal/spacetrack"
)

func main() {
	raw, err := spacetrack.FetchCatalog(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	objects, err := parser.ParseCatalog(raw)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(objects); err != nil {
		log.Fatal(err)
	}
}
