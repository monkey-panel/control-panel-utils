package database

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct{ *gorm.DB }

func NewDB(filename string, setupModel func(db *DB)) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})

	new_db := &DB{db}

	setupModel(new_db)

	return new_db, err
}
