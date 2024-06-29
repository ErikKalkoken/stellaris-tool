package main

import (
	"archive/zip"
	"encoding/json"
	"example/stellaris-tool/parser"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("Reading zip file...")
	r, err := zip.OpenReader(".temp/ironman.sav")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	name := "gamestate"
	for _, f := range r.File {
		if f.Name == name {
			r, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Parsing %s...\n", f.Name)
			p := parser.NewParser(r)
			x, err := p.Parse()
			if err != nil {
				log.Fatal(err)
			}
			r.Close()
			y, err := json.MarshalIndent(x, "", "    ")
			if err != nil {
				panic(err)
			}
			fn := fmt.Sprintf(".temp/%s.json", name)
			log.Printf("Writing %s...\n", fn)
			if err := os.WriteFile(fn, y, 0644); err != nil {
				panic(err)
			}
		}
	}

	// f, err := os.Open(".temp/gamestate")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// p := parser.NewParser(f)
	// x, err := p.Parse()
	// if err != nil {
	// 	panic(err)
	// }
	// y, err := json.MarshalIndent(x, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// if err := os.WriteFile(".temp/gamestate.json", y, 0644); err != nil {
	// 	panic(err)
	// }
}

// func loadSaveFile() {
// 	log.Println("Reading zip file...")
// 	r, err := zip.OpenReader("ironman.sav")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer r.Close()

// 	files := make(map[string][]byte)
// 	for _, f := range r.File {
// 		log.Printf("found file %s:\n", f.Name)
// 		rc, err := f.Open()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		dat, err := io.ReadAll(rc)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rc.Close()
// 		files[f.Name] = dat
// 	}
// }
