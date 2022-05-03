package db

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
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

func GetOne[T1 any](rst *T1, args ...interface{}) *T1 {
	dbCC.First(rst, args...)
	return rst
}
func GetRmtsvLists4List[T1, T2 any](mode T1, aRst []T2) []T2 {
	dbCC.Model(&mode).Limit(10000).Find(&aRst)
	return aRst
}

// 通过泛型调用,支持多个模型调用
func GetRmtsvLists[T1, T2 any](g *gin.Context, mode T1, aRst []T2) {
	//rst := dbCC.Model(&mode).Limit(10000).Find(&aRst)
	aRst = GetRmtsvLists4List(mode, aRst)
	if nil != aRst && 0 < len(aRst) {
		g.JSON(http.StatusOK, aRst)
		return
	}
	g.JSON(http.StatusBadRequest, gin.H{"msg": "not found", "code": -1})
}
