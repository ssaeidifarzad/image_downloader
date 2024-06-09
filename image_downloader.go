package main

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbonce   sync.Once
	database *gorm.DB
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
		database = gormconn
	})
	return database
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

func ResizeImage(data []byte, width int, height int) (*Image, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, resizedImg, nil)
	case "png":
		err = png.Encode(&buf, resizedImg)
	default:
		return nil, errors.New("unsupported image format")
	}
	if err != nil {
		return nil, err
	}
	return &Image{data: buf.Bytes(), format: format}, nil
}

func StoreImage(img *Image) error {
	db := GetDB()
	return db.Create(img).Error
}
