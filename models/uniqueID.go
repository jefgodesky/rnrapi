package models

import (
	"gorm.io/gorm"
	"math/rand"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateID() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

func IsIDUnique(tx *gorm.DB, model interface{}, id string) bool {
	var count int64
	tx.Model(model).Where("id = ?", id).Count(&count)
	return count == 0
}

func UniqueIDBeforeCreate(tx *gorm.DB, model interface{}, id *string) error {
	for {
		*id = GenerateID()
		if IsIDUnique(tx, model, *id) {
			break
		}
	}
	return nil
}
