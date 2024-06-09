package main

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbonce sync.Once
	db     *gorm.DB
)

type Image struct {
	gorm.Model
	data   []byte
	format string
}

func GetDB() *gorm.DB {
	dbonce.Do(func() {
		gormconn, err := ConnectDB()
		if err != nil {
			panic(err)
		}
		db = gormconn
	})
	return db
}

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	var (
		DbName   = os.Getenv("POSTGRES_DB")
		User     = os.Getenv("POSTGRES_USER")
		Password = os.Getenv("POSTGRES_PASSWORD")
		Host     = os.Getenv("POSTGRES_HOST")
		Port     = os.Getenv("POSTGRES_PORT")
	)
	dsn := "host=" + Host + " user=" + User + " password=" + Password + " dbname=" + DbName + " port=" + Port + " sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
