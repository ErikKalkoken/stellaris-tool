package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ErikKalkoken/stellaris-tool/internal/parser"
)

func main() {
	flag.Usage = myUsage
	destPtr := flag.String("d", ".", "destination path for output file")
	flag.Parse()
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	source := flag.Arg(0)
	data, err := parseFile(source)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	basename := filepath.Base(source)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))
	if err := writeJson(*destPtr, name, data); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func myUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: pdx2json [options] <inputfile>:\nOptions:\n")
	flag.PrintDefaults()
}

func parseFile(path string) (map[string][]any, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	fmt.Printf("Parsing %s...\n", path)
	p := parser.NewParser(r)
	data, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeJson(path string, name string, data map[string][]any) error {
	y, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	fn := fmt.Sprintf("%s/%s.json", path, name)
	fmt.Printf("Writing %s...\n", fn)
	if err := os.WriteFile(fn, y, 0644); err != nil {
		return err
	}
	return nil
}
