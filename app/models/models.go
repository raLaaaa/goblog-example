package models

import "time"

// The blog entry
type BlogEntry struct {
	ID          int       `storm:"id,increment"`
	Name        string    `storm:"index"`
	Description string    `storm:"index"`
	CreatedAt   time.Time `storm:"index"`
}

// The user for login in at the backend.
// No auto increment since it is designed to only have one user currently which is hardcoded into the DB.
type User struct {
	ID             int       `storm:"id"`
	Name           string    `storm:"index"`
	HashedPassword []byte    `storm:"index"`
	CreatedAt      time.Time `storm:"index"`
}
