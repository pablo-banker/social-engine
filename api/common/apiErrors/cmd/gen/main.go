// Generates error_code.json from the apiErrors.All registry.
// Run via: go generate ./common/apiErrors/...
package main

import (
	"log"
	"os"

	apiErrors "social-engine/common/apiErrors"
)

const outputFile = "error_code.json"

func main() {
	data, err := apiErrors.RenderJSON(apiErrors.All)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(outputFile, data, 0o644); err != nil {
		log.Fatalf("write %s: %v", outputFile, err)
	}
	log.Printf("wrote %d entries to %s", len(apiErrors.All), outputFile)
}
