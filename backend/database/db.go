package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect opens the database and runs migrations for the given models.
// Models are passed in by main.go so this package stays domain-agnostic.
func Connect(dsn string, models ...any) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}
	log.Println("database connected and migrated")
	return db
}
