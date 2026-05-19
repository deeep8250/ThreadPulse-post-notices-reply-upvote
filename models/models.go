package models

import (
	"time"
)

type Register struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type User struct {
	Id           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	HashPassword string    `json:"-" db:"hashed_pass"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CreateThread struct {
	Id        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title" binding:"required"`
	Content   string    `json:"content" db:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UpdateThread struct {
	Title   string `json:"title" db:"title" `
	Content string `json:"content" db:"content" `
}

type Replies struct {
	Id        int       `json:"id" db:"id"`
	PostID    int       `json:"post_id" db:"post_id"`
	UserID    int       `json:"user_id" db:"replied_user_id"`
	Reply     string    `json:"reply" db:"reply" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type HotThread struct {
	Id          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpvoteCount int       `json:"upvote_count" db:"upvote_count"`
}
