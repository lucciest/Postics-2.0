package models

import (
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Anons     string    `json:"anons"`
	FullText  string    `json:"fullText"`
	AuthorID  int       `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewPost(title, anons, fullText string, authorID int) *Post {
	now := time.Now()
	return &Post{
		Title:     title,
		Anons:     anons,
		FullText:  fullText,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
