// Package main is a gowsdl-based code-generation wrapper that generates Go
// binding code from a given WSDL file into a target directory.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// main parses command-line flags and invokes gowsdl to generate Go code from
// the given WSDL file.
func main() {
	var (
		wsdlPath    = flag.String("wsdl", "", "Path to WSDL file")
		outputDir   = flag.String("output", "pkg/generated", "Output directory for generated code")
		packageName = flag.String("package", "generated", "Package name for generated code")
	)
	flag.Parse()

	if *wsdlPath == "" {
		log.Fatal("Please provide WSDL file path using -wsdl flag")
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate code using gowsdl
	cmd := exec.Command("gowsdl", "-o", *outputDir, "-p", *packageName, *wsdlPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Generating code from WSDL: %s\n", *wsdlPath)
	fmt.Printf("Output directory: %s\n", *outputDir)
	fmt.Printf("Package name: %s\n", *packageName)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Println("Code generation completed successfully!")
}
