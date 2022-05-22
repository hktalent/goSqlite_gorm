package main

import (
	"github.com/gin-gonic/gin"
	"goSqlite_gorm/pkg/db"
	mymod "goSqlite_gorm/pkg/models"
	initrt "goSqlite_gorm/pkg/server"
	task "goSqlite_gorm/pkg/task"
	"gorm.io/gorm"
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
	// 包含启动ssh server，方便通过web连接本地的shell
	go task.DoAllTask()

	if nil != dbCC {
		router := gin.Default()
		initrt.InitRoute(router)
		//router.Use(ConnRmtSvsMiddleware())

		// swagger 似乎成了所有例子的路径
		//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		router.Run(":8081")
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
