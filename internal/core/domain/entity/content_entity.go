package entity

import "time"

type ContentEntity struct {
	ID          int64
	Title       string
	CategoryID  int64
	CreatedByID int64
	Excerpt     string
	Description string
	Image       string
	Status      string
	Tags        []string
	Category    CategoryEntity
	User        UserEntity
	CreatedAt   time.Time
}

type QueryString struct {
	Limit      int
	Page       int
	Order      string
	OrderType  string
	Search     string
	CategoryID int64
}
