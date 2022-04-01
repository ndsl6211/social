package pkg

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteGormClient() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./db/gorm.db"), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				// LogLevel:      logger.Info,
				Colorful: true,
			},
		),
	})
	if err != nil {
		panic("failed to init db")
	}

	return db
}

func NewMemoryGormClient() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{})

	if err != nil {
		panic("failed to init in-memory db")
	}

	return db
}
