package data

import (
	"time"
)

type Document struct {
	Title       string
	Description string
	CreatedAt   time.Time
	Thumbnail   string
	Tags        []Tag
}

type Tag struct {
	ID    int
	Name  string
	Color string
}
