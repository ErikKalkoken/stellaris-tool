package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ErikKalkoken/stellaris-tool/internal/parser"
)

func main() {
	flag.Usage = myUsage
	destPtr := flag.String("d", ".", "destination path for output file")
	rawPrt := flag.Bool("k", false, "keep raw data files")
	flag.Parse()
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	source := flag.Arg(0)
	if err := processSaveFile(source, *destPtr, *rawPrt); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func processSaveFile(source string, dest string, keepDataFiles bool) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()
	fmt.Printf("Processing save file: %s\n", source)
	for _, f := range r.File {
		if keepDataFiles {
			if err := writeData(dest, f); err != nil {
				fmt.Printf("ERROR: Failed to write data file for %s: %s\n", f.Name, err)
				continue
			}

		}
		data, err := parseFile(f)
		if err != nil {
			fmt.Printf("ERROR: Failed to parse %s: %s\n", f.Name, err)
			continue
		}
		if err := writeJson(dest, f.Name, data); err != nil {
			fmt.Printf("ERROR: Failed to write JSON for %s: %s\n", f.Name, err)
			continue
		}
	}
	return nil
}

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

func writeData(path string, f *zip.File) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	fmt.Printf("Writing data file: %s\n", f.Name)
	name := fmt.Sprintf("%s/%s", path, f.Name)
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

func writeJson(path string, name string, data map[string][]any) error {
	y, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	fn := fmt.Sprintf("%s/%s.json", path, name)
	fmt.Printf("Writing JSON: %s\n", fn)
	if err := os.WriteFile(fn, y, 0644); err != nil {
		return err
	}
	return nil
}

func myUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: sav2json [options] <inputfile>:\nOptions:\n")
	flag.PrintDefaults()
}
