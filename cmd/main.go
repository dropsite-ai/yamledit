package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dropsite-ai/yamledit"
	"gopkg.in/yaml.v3"
)

func main() {
	// Define command-line flags.
	op := flag.String("op", "", "Operation: 'read' or 'update'")
	filePath := flag.String("file", "", "Path to the YAML file")
	dotPath := flag.String("path", "", "Dot-notation path to the value")
	newValue := flag.String("value", "", "New value for update (as YAML literal)")
	outFile := flag.String("out", "", "Optional output file (if not provided, prints to stdout)")

	flag.Parse()

	// Basic argument validation.
	if *op == "" || *filePath == "" || *dotPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Read the YAML file.
	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	switch *op {
	case "read":
		// Use an empty interface to capture any kind of value.
		var result interface{}
		if err := yamledit.Read(data, *dotPath, &result); err != nil {
			log.Fatalf("Failed to read value: %v", err)
		}
		// Marshal the result back into YAML for pretty printing.
		out, err := yaml.Marshal(result)
		if err != nil {
			log.Fatalf("Failed to marshal result: %v", err)
		}
		if *outFile != "" {
			if err := os.WriteFile(*outFile, out, 0644); err != nil {
				log.Fatalf("Failed to write output file: %v", err)
			}
		} else {
			fmt.Print(string(out))
		}

	case "update":
		// The update operation requires a new value.
		if *newValue == "" {
			log.Fatalf("Update operation requires a non-empty -value flag")
		}
		// Parse the newValue as a YAML literal so that numbers, booleans, etc. are handled correctly.
		var parsedValue interface{}
		if err := yaml.Unmarshal([]byte(*newValue), &parsedValue); err != nil {
			log.Fatalf("Failed to parse new value: %v", err)
		}
		updated, err := yamledit.Update(data, *dotPath, parsedValue)
		if err != nil {
			log.Fatalf("Failed to update YAML: %v", err)
		}
		if *outFile != "" {
			if err := os.WriteFile(*outFile, updated, 0644); err != nil {
				log.Fatalf("Failed to write output file: %v", err)
			}
		} else {
			fmt.Print(string(updated))
		}

	default:
		log.Fatalf("Unknown operation: %s", *op)
	}
}
