package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // this is needed by Gorm
)

// Create creates a database instance using a dialect and database name
func Create(dialect string, dbname string) (*gorm.DB, error) {
	// for now
	return gorm.Open(dialect, dbname)
}
