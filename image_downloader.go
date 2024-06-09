package main

import (
	"errors"
	"io"
	"net/http"
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

func Migrate() {
	db := GetDB()
	db.AutoMigrate(&Image{})
}

func DownloadImage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download image")
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
