package services

import (
	"fmt"
	"log"
	"rala-blog/app/models"

	"github.com/asdine/storm"
)

const DB_PATH string = "my.db"

func SaveToDatabase(entry models.BlogEntry) {
	db, err := storm.Open(DB_PATH)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = db.Save(&entry)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer db.Close()
}

func GetAllEntries() []models.BlogEntry {
	var entries []models.BlogEntry

	db, err := storm.Open(DB_PATH)

	if err != nil {
		log.Fatal(err)
	}

	errAll := db.All(&entries)

	if errAll != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return entries

}
