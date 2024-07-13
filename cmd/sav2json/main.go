package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ErikKalkoken/stellaris-tool/internal/parser"
)

// Current version need to be injected via ldflags
var Version = "?"

func main() {
	flag.Usage = myUsage
	destFlag := flag.String("d", ".", "destination directory for output files")
	keepFlag := flag.Bool("k", false, "keep original data files")
	sameFlag := flag.Bool("s", false, "create output files in same directory as source files")
	versionFlag := flag.Bool("v", false, "show the current version")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("sav2json %s\n", Version)
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	source := flag.Arg(0)
	var dest string
	if *sameFlag {
		dest = filepath.Dir(source)
	} else {
		dest = *destFlag
	}
	if err := processSaveFile(source, dest, *keepFlag); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

// myUsage writes a custom usage message to configured output stream.
func myUsage() {
	s := "Usage: sav2json [options] <inputfile>:\n\n" +
		"A tool for converting Stellaris save games into JSON.\n" +
		"For more information please see: https://github.com/ErikKalkoken/stellaris-tool\n\n" +
		"Options:\n"
	fmt.Fprint(flag.CommandLine.Output(), s)
	flag.PrintDefaults()
}

// processSaveFile writes the contents of a Stellaris safe game file in JSON format to disk.
// It will optionally also write the raw data files to disk, when keepDataFiles is true.
func processSaveFile(source string, dest string, keepDataFiles bool) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()
	var hasErrors bool
	fmt.Printf("Processing save file: %s\n", source)
	for _, f := range r.File {
		if keepDataFiles {
			if err := writeData(dest, f); err != nil {
				fmt.Printf("ERROR: Failed to write data file for %s: %s\n", f.Name, err)
				hasErrors = true
				continue
			}

		}
		data, err := parseFile(f)
		if err != nil {
			fmt.Printf("ERROR: Failed to parse %s: %s\n", f.Name, err)
			hasErrors = true
			continue
		}
		if err := writeJSON(dest, f.Name, data); err != nil {
			fmt.Printf("ERROR: Failed to write JSON for %s: %s\n", f.Name, err)
			hasErrors = true
			continue
		}
	}
	if hasErrors {
		return errors.New("processing failed with errors")
	}
	return nil
}

// parseFile parses a zip file and returns it's contents.
func parseFile(f *zip.File) (map[string][]any, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	fmt.Printf("Parsing file: %s\n", f.Name)
	p := parser.NewParser(r)
	data, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// writeData writes a a zip file raw to disk.
func writeData(dir string, f *zip.File) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	p := fmt.Sprintf("%s/%s", dir, f.Name)
	fmt.Printf("Writing data file: %s\n", p)
	w, err := os.Create(p)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

// writeJSON writes the given data to disk.
func writeJSON(dir string, name string, data map[string][]any) error {
	y, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	p := fmt.Sprintf("%s/%s.json", dir, name)
	fmt.Printf("Writing JSON: %s\n", p)
	if err := os.WriteFile(p, y, 0644); err != nil {
		return err
	}
	return nil
}
