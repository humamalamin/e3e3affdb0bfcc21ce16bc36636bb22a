package seeds

import (
	"latihan-portal-news/internal/core/domain/model"
	"latihan-portal-news/lib/conv"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	// Jika belum ada, buat role 'admin'
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Gagal mengenkripsi password: %v", err)
	}
	admin := model.User{
		Name:     "admin",
		Email:    "admin@mail.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@mail.com"}).Error; err != nil {
		log.Fatalf("Gagal membuat role 'admin': %v", err)
	} else {
		log.Println("Role 'admin' berhasil di-seed.")
	}
}
