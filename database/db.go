package database

import (
	"log"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetDb() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/db_pilar?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func SetDbTest() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/db_pilar_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
