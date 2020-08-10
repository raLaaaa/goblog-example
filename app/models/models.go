package models

import "time"

type BlogEntry struct {
	ID          int       `storm:"id,increment"`
	Name        string    `storm:"index"`
	Description string    `storm:"index"`
	CreatedAt   time.Time `storm:"index"`
}

type User struct {
	ID             int       `storm:"id"`
	Name           string    `storm:"index"`
	HashedPassword []byte    `storm:"index"`
	CreatedAt      time.Time `storm:"index"`
}
