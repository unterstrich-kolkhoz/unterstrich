package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Create(dialect string, dbname string) (*gorm.DB, error) {
	db, err := gorm.Open(dialect, dbname)

	if err != nil {
		return nil, err
	}

	return db, nil
}
