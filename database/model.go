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
