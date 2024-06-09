package main

import "gorm.io/gorm"

type DBconn struct {
	db *gorm.DB
}

type Image struct {
	gorm.Model
	data   []byte
	format string
}
