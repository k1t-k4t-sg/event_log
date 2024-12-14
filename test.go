package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open("Devices.json")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("stat: %v", err)
	}

	//ModTime() time.Time

	log.Printf("file %q: size: %d, mod. time: %q", stat.Name(), stat.Size(), stat.ModTime()) 
	log.Printf("file sys: %T", stat.Sys())

}
