package db

import "time"

type Location struct {
	Id       int       `db:"id"`
	Username string    `db:"username"`
	Created  time.Time `db:"created"`
	Lat      float32   `db:"lat"`
	Lng      float32   `db:"lng"`
	Alt      float32   `db:"alt"`
}
