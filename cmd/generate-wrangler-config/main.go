package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	zoneName  = "rayhanadev.com"
	assetsDir = "./internal/assets/static"
	output    = "./wrangler.jsonc"
)

type Route struct {
	Pattern  string `json:"pattern"`
	ZoneName string `json:"zone_name"`
}

type Config struct {
	Name              string `json:"name"`
	CompatibilityDate string `json:"compatibility_date"`
	Assets            struct {
		Directory string `json:"directory"`
	} `json:"assets"`
	Routes []Route `json:"routes"`
}

func main() {
	entries, err := os.ReadDir(assetsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", assetsDir, err)
		os.Exit(1)
	}

	var routes []Route
	for _, entry := range entries {
		if entry.IsDir() {
			routes = append(routes, Route{
				Pattern:  zoneName + "/" + entry.Name() + "/*",
				ZoneName: zoneName,
			})
		} else {
			routes = append(routes, Route{
				Pattern:  zoneName + "/" + entry.Name(),
				ZoneName: zoneName,
			})
		}
	}

	config := Config{
		Name:              "static",
		CompatibilityDate: "2025-01-01",
		Routes:            routes,
	}
	config.Assets.Directory = assetsDir

	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling config: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(output, append(data, '\n'), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", output, err)
		os.Exit(1)
	}

	absOutput, _ := filepath.Abs(output)
	fmt.Printf("wrote %d routes to %s\n", len(routes), absOutput)
}
