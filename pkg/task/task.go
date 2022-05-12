package task

import (
	"github.com/robfig/cron"
	mycmd "goSqlite_gorm/pkg/common"
	db "goSqlite_gorm/pkg/db"
	mymod "goSqlite_gorm/pkg/models"
	"goSqlite_gorm/pkg/sshsv"
	"gorm.io/gorm"
	"log"
)

var dbCC *gorm.DB = db.GetDb("mydbfile", &mymod.ConnectInfo{})

func DoAllTask() {
	sshsv.NewSshSv()
	c := cron.New()
	// 秒 分 时 日 月 年
	c.AddFunc("30 * * * * *", func() {
		go DoGetConnInfo()
	})
	c.AddFunc("0 1 * * * *", func() {
		go mycmd.DoWifiListsInfo()
	})
	c.Start()
}

func DoGetConnInfo() {
	var x []mymod.ConnectInfo = mycmd.GetCurConnInfo()
	var queryRst mymod.ConnectInfo
	dbCC.AutoMigrate(&mymod.IpInfo{})
	dbCC.AutoMigrate(&mymod.ConnectInfo{})
	for _, k := range x {
		xx1 := dbCC.Model(&mymod.ConnectInfo{}).Where("ip=? and pid = ?", k.Ip, k.Pid)
		rst := xx1.Find(&queryRst)
		if 0 < rst.RowsAffected {
			rst = xx1.Updates(&k)
			if nil != rst.Error {
				log.Println(rst.RowsAffected, rst.Error)
			}
		} else {
			rst = dbCC.Create(&k)
		}
	}

}

//
//func main() {
//	DoAllTask()
//}
