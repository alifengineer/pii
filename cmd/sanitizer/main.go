package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alifengineer/pii/internal/analyzer"
	"github.com/alifengineer/pii/internal/bufferpool"
	"github.com/alifengineer/pii/internal/generator"
	"github.com/alifengineer/pii/internal/parser"
)

var (
	destination = flag.String("destination", "", "output file path (default: <source>_sanitize.go)")
	typeFlag    = flag.String("type", "", "comma-separated list of struct names to generate for")
	verbose     = flag.Bool("v", false, "print per-field reasoning to stderr")
	debug       = flag.Bool("debug", false, "dump parsed model to stdout without generating")
	check       = flag.Bool("check", false, "compare output to existing file, exit non-zero if different")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pii: ")
	flag.Usage = usage
	flag.Parse()

	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		log.Fatal("GOFILE environment variable not set (run via go:generate)")
	}

	goPackage := os.Getenv("GOPACKAGE")
	if goPackage == "" {
		log.Fatal("GOPACKAGE environment variable not set (run via go:generate)")
	}

	dest := *destination
	if dest == "" {
		dest = defaultDestination(goFile)
	}

	_, structs, err := parser.Parse(".")
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	if *debug {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(structs); err != nil { //nolint:musttag //best effort
			log.Fatalf("encoding model: %v", err)
		}
		return
	}

	var typeFilter map[string]bool
	if *typeFlag != "" {
		typeFilter = make(map[string]bool)
		for name := range strings.SplitSeq(*typeFlag, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				typeFilter[name] = true
			}
		}
	}

	targets, autoIncluded, err := analyzer.Analyze(structs, typeFilter)
	if err != nil {
		log.Fatalf("analyze: %v", err)
	}

	if *verbose {
		for _, name := range autoIncluded {
			fmt.Fprintf(os.Stderr, "pii: auto-including %s (dependency of a requested type)\n", name)
		}
		for _, t := range targets {
			for _, f := range t.Fields {
				fmt.Fprintf(os.Stderr, "%s.%s: %s\n", t.StructName, f.Name, describeAction(f.Action))
			}
		}
	}

	if len(targets) == 0 {
		return
	}

	buf := bufferpool.Get()

	output, err := generator.Generate(goPackage, targets, buf)
	if err != nil {
		bufferpool.Put(buf)
		log.Fatalf("generate: %v", err) //nolint:gosec //best effort
	}
	bufferpool.Put(buf)

	if *check {
		existing, err := os.ReadFile(dest) //nolint:gosec // path from go:generate, not user input
		if err != nil {
			log.Fatalf("check: cannot read %s: %v", dest, err) //nolint:gosec //best effort
		}
		if !bytes.Equal(existing, output) {
			log.Fatalf("check: %s is out of date — run go generate", dest) //nolint:gosec //best effort
		}
		return
	}

	if !generator.ShouldWrite(output, dest) {
		return
	}

	if err := generator.WriteFile(output, dest); err != nil {
		log.Fatalf("write: %v", err)
	}
}

func describeAction(a analyzer.FieldAction) string {
	switch a {
	case analyzer.ActionSanitize:
		return "sanitize"
	case analyzer.ActionSanitizePtr:
		return "sanitize (pointer)"
	case analyzer.ActionSanitizeSlice:
		return "sanitize (slice)"
	case analyzer.ActionSanitizePtrSlice:
		return "sanitize (pointer slice)"
	default:
		return "copy"
	}
}

func defaultDestination(goFile string) string {
	ext := filepath.Ext(goFile)
	return strings.TrimSuffix(goFile, ext) + "_sanitize.go"
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: pii [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Pii generates Sanitize() methods for structs with sanitizable fields.\n")
	fmt.Fprintf(os.Stderr, "Run via //go:generate pii in your Go source files.\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}
