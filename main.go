package main

import (
	"flag"
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/html"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/projectExport"
	"log"
	"os"
	"regexp"
)

var (
	scopeCleanupPattern = regexp.MustCompile(`\s{2,}|[\t\n\r]+`)
)

func mainInternal() error {
	password := flag.String("password", "", "The password used when exporting the project.")
	filename := flag.String("export-filename", "export.html", "The name of the HTML report file to generate.")
	sensitiveOnly := flag.Bool("sensitive-only", false, "Only show sensitive variables.")

	flag.Parse()

	if *password == "" {
		return fmt.Errorf("missing the password used on export, specify the --password argument")
	}

	key, err := decryptor.KeyFromPassword(*password)
	if err != nil {
		return fmt.Errorf("failed to derive key from password: %w", err)
	}

	// read all known entity types
	log.Printf("Reading files...")

	entities := make(projectExport.EntityMap)

	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed reading directory: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		// read the entity
		entity, err := projectExport.ReadEntity(f.Name())
		if err != nil {
			log.Printf("Ignoring file '%s': %s", f.Name(), err)
			continue
		}

		if entity != nil {
			entity.AddToEntityMap(entities)
		}
	}

	log.Printf("Loaded %d entities.", len(entities))

	// render the html report
	content, err := html.RenderHTML(entities, key, *sensitiveOnly)
	if err != nil {
		return fmt.Errorf("failed render HTML report: %w", err)
	}

	log.Printf("Writing report to '%s'...", *filename)

	err = os.WriteFile(*filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed render HTML report: %w", err)
	}

	return nil
}

func main() {
	log.SetOutput(os.Stdout)

	err := mainInternal()
	if err != nil {
		log.Printf("UNEXPECTED ERROR: %v", err)
		os.Exit(1)
	} else {
		log.Printf("DONE.")
	}
}
