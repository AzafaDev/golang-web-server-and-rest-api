package main

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type APIError struct {
	Error string `json:"error"`
}
