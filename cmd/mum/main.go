package main

import (
	"log"

	"github.com/tylerchambers/mum/internal/db"
	"github.com/tylerchambers/mum/internal/parser"
)

const mySecret = "password1"

func main() {
	myDB, err := db.Create("./secrets")
	if err != nil {
		log.Fatal(err)
	}
	err = db.AddCred(myDB, "password1", "production server root password")
	if err != nil {
		log.Fatal(err)
	}
	myWords := parser.FindWordsWithMapReduce("./cmd/mum/", 20)
	for k := range myWords {
		h, err := db.HashCred(k)
		if err != nil {
			log.Fatal(err)
		}
		exists, err := db.CredExists(myDB, h)
		if exists {
			log.Fatalf("cred: %s\nwith hash: %s\nexists!", k, h)
		} else {
			log.Printf("did not find cred: %s, with hash: %s", k, h)
		}
	}
}
