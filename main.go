package main

func main() {

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
