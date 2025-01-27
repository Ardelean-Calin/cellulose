package main

import "github.com/Ardelean-Calin/cellulose/internal/db"

func main() {
	db, err := db.NewDB(db.Config{DatabasePath: "cellulose.db"})
	if err != nil {
		panic(err)
	}
	db.Init()
}
