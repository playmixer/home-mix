package database

import "gorm.io/gorm"

type Ping struct {
	IP     string
	Name   string
	Mac    string
	Online bool
	Look   bool
	gorm.Model
}

type Device struct {
	gorm.Model
	Mac   string
	Owner string
}
