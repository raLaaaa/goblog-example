package models

type BlogEntry struct {
	ID          int    `storm:"id,increment"`
	Name        string `storm:"index"`
	Description string `storm:"index"`
}

type User struct {
	ID       int    `storm:"id,increment"`
	Name     string `storm:"index"`
	Password string `storm:"index"`
}
