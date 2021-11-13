package database

import (
	"fmt"
	"strings"

	"delineate.io/customers/src/config"
	"delineate.io/customers/src/logging"
	"delineate.io/customers/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize() {
	if db, err := OpenDB(); err != nil {
		panic(err)
	} else {
		err = db.AutoMigrate(&models.Customer{})
		if err != nil {
			panic(err)
		}
		logging.Info("successfully initialized the database")
	}
}

func OpenDB() (*gorm.DB, error) {
	dsn := getConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Warn(dsn)
		logging.Err(err)
	}
	return db, err
}

func getConnectionString() string {
	var builder strings.Builder
	values := config.GetSection("database")
	for k, v := range values {
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(v)
		builder.WriteString(" ")
	}
	cs := strings.TrimSpace(builder.String())
	logging.Debug(fmt.Sprintf("connection string: '%s'", cs))
	return cs
}
