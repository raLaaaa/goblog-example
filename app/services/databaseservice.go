package services

import (
	"fmt"
	"log"
	"rala-blog/app/models"
	"time"

	"github.com/asdine/storm"
	"golang.org/x/crypto/bcrypt"
)

// Local path to the DB file
const DB_PATH string = "my.db"

// Saves an blogentry to the DB
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

// Dummy method to create a base user
// Username: rala
// Password: 123
func CreateBaseUser() {

	user := models.User{
		ID:        1,
		Name:      "rala",
		CreatedAt: time.Now(),
	}

	//Encryption
	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte("123"), bcrypt.DefaultCost)

	db, err := storm.Open(DB_PATH)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = db.Save(&user)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer db.Close()
}

// Fetches a single user from the DB
func GetSingleUserByName(name string) (models.User, error) {
	db, err := storm.Open(DB_PATH)

	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	var user models.User
	errSingle := db.One("Name", name, &user)

	if errSingle != nil {
		fmt.Println(errSingle)
		return models.User{}, err
	}

	defer db.Close()

	return user, nil
}

// Fetches all blog entries from the DB
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
