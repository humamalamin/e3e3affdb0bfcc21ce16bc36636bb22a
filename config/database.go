package config

import (
	"fmt"
	"latihan-portal-news/database/seeds"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresDB.DBUser,
		cfg.PostgresDB.DBPass,
		cfg.PostgresDB.DBHost,
		cfg.PostgresDB.DBPort,
		cfg.PostgresDB.DBName,
	)
	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		log.Error().Err(err).Msg("[Connect-2] failed to connect to database " + cfg.PostgresDB.DBHost)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[Connect-2] failed to connect to database " + cfg.PostgresDB.DBHost)
		return nil, err
	}

	seeds.SeedRoles(db)

	sqlDB.SetMaxOpenConns(cfg.PostgresDB.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.PostgresDB.DBMaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Postgres{
		DB: db,
	}, nil
}
