package main

import (
	"archive/zip"
	"encoding/json"
	"example/stellaris-tool/internal/parser"
	"fmt"
	"log"
	"os"
)

func main() {
	meta, gamestate, err := loadSaveFile(".temp/ironman.sav")
	if err != nil {
		log.Fatalf("Failed to load save file: %s", err)
	}
	if err := writeJson("meta", meta); err != nil {
		log.Fatal("Failed to write JSON file")
	}
	if err := writeJson("gamestate", gamestate); err != nil {
		log.Fatal("Failed to write JSON file")
	}
}

func loadSaveFile(path string) (map[string][]any, map[string][]any, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	log.Printf("Opened sav file %s", path)

	var meta, gamestate map[string][]any
	for _, f := range r.File {
		if f.Name == "meta" {
			meta, err = parseFile(f)
			if err != nil {
				return nil, nil, err
			}
		}
		if f.Name == "gamestate" {
			gamestate, err = parseFile(f)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return meta, gamestate, nil
}

func parseFile(f *zip.File) (map[string][]any, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	log.Printf("Parsing %s...\n", f.Name)
	p := parser.NewParser(r)
	data, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeJson(name string, data map[string][]any) error {
	y, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	fn := fmt.Sprintf(".temp/%s.json", name)
	log.Printf("Writing %s...\n", fn)
	if err := os.WriteFile(fn, y, 0644); err != nil {
		return err
	}
	return nil
}
