package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/hktalent/goSqlite_gorm/pkg/db"
	mymod "github.com/hktalent/goSqlite_gorm/pkg/models"
	initrt "github.com/hktalent/goSqlite_gorm/pkg/server"
	task "github.com/hktalent/goSqlite_gorm/pkg/task"
	"gorm.io/gorm"
	"strconv"
)

var dbCC *gorm.DB = db.GetDb(&mymod.RemouteServerce{})

// http://localhost:8080/swagger/index.html
// c.Header("Content-Type", "application/json")
// @title 51pwn app API
// @version 1.0
// @description This is 51pwn app api docs.
// @license.name Apache 2.0
// @contact.name go-swagger
// @contact.url https://github.com/hktalent/
// @host localhost:8080
// @BasePath /api/v1
func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "port")
	// 包含启动ssh server，方便通过web连接本地的shell
	go task.DoAllTask()

	if nil != dbCC {
		router := gin.Default()
		initrt.InitRoute(router)
		//router.Use(ConnRmtSvsMiddleware())
		// swagger 似乎成了所有例子的路径
		//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		router.Run(":" + strconv.Itoa(port))
		//var x509 tls.Certificate
		//x509, err := tls.LoadX509KeyPair(SSLCRT, SSLKEY)
		//if err != nil {
		//	return
		//}
		//var server *http.Server
		//server = &http.Server{
		//	Addr:    ":8081",
		//	Handler: router,
		//	TLSConfig: &tls.Config{
		//		Certificates: []tls.Certificate{x509},
		//	},
		//}
		//server.ListenAndServe()
	}
}
