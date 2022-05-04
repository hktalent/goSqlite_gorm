package db

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
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

// 求count，支持条件
func GetCount[T any](mod T, args ...interface{}) int64 {
	var n int64
	if 0 < len(args) {
		dbCC.Model(&mod).Where(args[0], args[1:]...).Count(&n)
	} else {
		dbCC.Model(&mod).Count(&n)
	}
	return n
}

// 查询返回一条数据
func GetOne[T1 any](rst *T1, args ...interface{}) *T1 {
	rst1 := dbCC.First(rst, args...)
	if 0 == rst1.RowsAffected && nil != rst1.Error {
		log.Println(rst1.Error)
		return nil
	}
	return rst
}
func GetRmtsvLists4List[T1, T2 any](mode T1, preLd string, aRst []T2, nPageSize int, Offset int, conds ...interface{}) []T2 {
	if "" != preLd {
		dbCC.Model(&mode).Preload(preLd).Limit(nPageSize).Offset(Offset*nPageSize).Find(&aRst, conds...)
	} else {
		dbCC.Model(&mode).Limit(nPageSize).Offset(Offset*nPageSize).Find(&aRst, conds...)
	}
	return aRst
}

// 通过泛型调用,支持多个模型调用
func GetRmtsvLists[T1, T2 any](g *gin.Context, mode T1, aRst []T2, conds ...interface{}) {
	//rst := dbCC.Model(&mode).Limit(10000).Find(&aRst)
	aRst = GetRmtsvLists4List(mode, "", aRst, 1000, 0, conds...)
	if nil != aRst && 0 < len(aRst) {
		g.JSON(http.StatusOK, aRst)
		return
	}
	g.JSON(http.StatusBadRequest, gin.H{"msg": "not found", "code": -1})
}
