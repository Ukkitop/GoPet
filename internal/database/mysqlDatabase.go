package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"realAPI/internal/database/models"
)

func GetDatabaseConnection() (*gorm.DB, error) {
	dsn := "admin:1234@(127.0.0.1:3306)/realapidatabase?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return db, nil
}

func SetupDatabase() error {
	db, conErr := GetDatabaseConnection()

	if conErr != nil {
		log.Error(conErr)
		return conErr
	}

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
