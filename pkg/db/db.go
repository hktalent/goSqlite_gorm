package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbCC *gorm.DB

func GetDb(dbName string, dst ...interface{}) *gorm.DB {
	if nil != dbCC {
		return dbCC
	}
	db, err := gorm.Open(sqlite.Open("file:"+dbName+".db?cache=shared&mode=rwc&_journal_mode=WAL&Synchronous=Off&temp_store=memory&mmap_size=30000000000"), &gorm.Config{})
	if err != nil {
		return nil
	}
	// Migrate the schema
	db.AutoMigrate(dst[0])
	dbCC = db
	return db
}
