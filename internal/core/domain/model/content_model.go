package model

import "time"

type Content struct {
	ID          int64      `gorm:"id"`
	Title       string     `gorm:"title"`
	CategoryID  int64      `gorm:"category_id"`
	CreatedByID int64      `gorm:"created_by_id"`
	Excerpt     string     `gorm:"excerpt"`
	Description string     `gorm:"description"`
	Image       string     `gorm:"image"`
	Status      string     `gorm:"status"`
	Tags        string     `gorm:"tags"`
	Category    Category   `gorm:"foreignKey:CategoryID"`
	User        User       `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time  `gorm:"created_at"`
	UpdatedAt   *time.Time `gorm:"updated_at"`
}
