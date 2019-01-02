package postgres

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Open starts a gorm postgres database instance
func Open(host, port, user, dbName, password string) (db *gorm.DB, err error) {
	dbArgs := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)
	db, err = gorm.Open("postgres", dbArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	return db, nil
}
