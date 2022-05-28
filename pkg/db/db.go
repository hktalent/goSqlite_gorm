package db

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"sync"
)

var dbCC *gorm.DB
var DoOnce sync.Once

// 获取Gorm db连接、操作对象
func GetDb(dst ...interface{}) *gorm.DB {
	dbName := "db/mydbfile"
	DoOnce.Do(func() {
		db, err := gorm.Open(sqlite.Open("file:"+dbName+".db?cache=shared&mode=rwc&_journal_mode=WAL&Synchronous=Off&temp_store=memory&mmap_size=30000000000"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err == nil {
			if nil != dst && 0 < len(dst) {
				db.AutoMigrate(dst[0])
			}
			dbCC = db
		}
	})

	return dbCC
}

// 通用
// 获取T类型mod表名
func GetTableName[T any](mod T) string {
	stmt := &gorm.Statement{DB: dbCC}
	stmt.Parse(&mod)
	return stmt.Schema.Table
}

// 通用,update
// 指定id更新T类型mod数据
func Update[T any](mod T, id interface{}) int64 {
	var t1 *T = &mod
	xxxD := dbCC.Table(GetTableName(mod)).Model(&t1)
	xxxD.AutoMigrate(t1)
	rst := xxxD.Where("id = ?", id).Updates(mod)
	if 0 >= rst.RowsAffected {
		log.Println(rst.Error)
	}
	return rst.RowsAffected
}

// 通用,insert
func Create[T any](mod T) int64 {
	var t1 *T = &mod
	xxxD := dbCC.Table(GetTableName(mod)).Model(&t1)
	xxxD.AutoMigrate(t1)
	rst := xxxD.Create(mod)
	if 0 >= rst.RowsAffected {
		log.Println(rst.Error)
	}
	return rst.RowsAffected
}

// 通用
// 求T类型count，支持条件
// 对T表，mod类型表，args 的where求count
func GetCount[T any](mod T, args ...interface{}) int64 {
	var n int64
	x1 := dbCC.Model(&mod)
	if 0 < len(args) {
		x1.Where(args[0], args[1:]...).Count(&n)
	} else {
		x1.Count(&n)
	}
	return n
}

// 通用
// 查询返回T类型、表一条数据
func GetOne[T any](rst *T, args ...interface{}) *T {
	rst1 := dbCC.First(rst, args...)
	if 0 == rst1.RowsAffected && nil != rst1.Error {
		//log.Println(rst1.Error)
		return nil
	}
	return rst
}

// 通用
// 查询模型T1类型 mode，并关联T1类型对子类型T3 preLd
// 设置 nPageSize 和便宜Offset
// 以及其他查询条件conds
func GetSubQueryLists[T1, T2 any](mode T1, preLd string, aRst []T2, nPageSize int, Offset int, conds ...interface{}) []T2 {
	if "" != preLd {
		dbCC.Model(&mode).Preload(preLd).Limit(nPageSize).Offset(Offset*nPageSize).Find(&aRst, conds...)
	} else {
		dbCC.Model(&mode).Limit(nPageSize).Offset(Offset*nPageSize).Find(&aRst, conds...)
	}
	return aRst
}

// 通用
// 查询模型T1类型 mode，并关联T1类型对子类型T3 preLd
// 设置 nPageSize 和便宜Offset
// 以及其他查询条件conds
func GetSubQueryList[T1, T2, T3 any](mode T1, preLd T3, aRst []T2, nPageSize int, Offset int, conds ...interface{}) []T2 {
	return GetSubQueryLists(mode, GetTableName(preLd), aRst, nPageSize, Offset, conds...)
}

// 通用
// 通过泛型调用,支持多个模型调用
// T1 继承了T2，存在包含关系
func GetRmtsvLists[T1, T2 any](g *gin.Context, mode T1, aRst []T2, conds ...interface{}) {
	//rst := dbCC.Model(&mode).Limit(10000).Find(&aRst)
	aRst = GetSubQueryLists(mode, "", aRst, 1000, 0, conds...)
	if nil != aRst && 0 < len(aRst) {
		g.JSON(http.StatusOK, aRst)
		return
	}
	g.JSON(http.StatusBadRequest, gin.H{"msg": "not found", "code": -1})
}
