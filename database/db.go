package database

import (
	"log"
	"os"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetDb() *gorm.DB {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "db_pilar"
	}
	dsn := "root:@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func SetDbTest() *gorm.DB {
	dbName := os.Getenv("DB_TEST_NAME")
	if dbName == "" {
		dbName = "db_pilar_test"
	}
	dsn := "root:@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
