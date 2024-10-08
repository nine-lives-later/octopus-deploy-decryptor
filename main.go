package main

import (
	"flag"
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/variableset"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	scopeCleanupPattern = regexp.MustCompile(`\s{2,}|[\t\n\r]+`)
)

func mainInternal() error {
	password := flag.String("password", "", "The password used when exporting the project.")
	sensitiveOnly := flag.Bool("sensitive-only", false, "Only show sensitive variables.")

	flag.Parse()

	if *password == "" {
		return fmt.Errorf("missing the password used on export, specify the --password argument")
	}

	key, err := decryptor.KeyFromPassword(*password)
	if err != nil {
		return fmt.Errorf("failed to derive key from password: %w", err)
	}

	// get all variable set files
	log.Printf("Reading variable set files...")

	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed reading directory: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasPrefix(f.Name(), "variableset-") || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		setFilePath := f.Name()

		log.Printf("File: %v", setFilePath)

		// read the variables
		set, err := variableset.ReadVariables(setFilePath)
		if err != nil {
			return fmt.Errorf("failed to read variable set file '%v': %w", setFilePath, err)
		}

		for _, v := range set {
			if *sensitiveOnly && v.Type != "Sensitive" {
				continue
			}

			val, err := v.DecryptedValue(key)
			if err != nil {
				log.Printf("Failed to decrypt variable '%v': %v", v.Name, err)
			}

			if len(v.Scope) > 0 {
				scope := string(v.Scope)
				scope = strings.ReplaceAll(scope, "\n", " ")
				scope = scopeCleanupPattern.ReplaceAllString(scope, " ")

				log.Printf("  Variable: %v = '%v' [scope: %+v]", v.Name, val, scope)
			} else {
				log.Printf("  Variable: %v = '%v' [scope: none]", v.Name, val)
			}
		}
	}
	return nil
}

func main() {
	err := mainInternal()
	if err != nil {
		log.Printf("UNEXPECTED ERROR: %v", err)
		os.Exit(1)
	} else {
		log.Printf("DONE.")
	}
}
