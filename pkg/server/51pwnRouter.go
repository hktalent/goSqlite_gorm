package server

import "github.com/gin-gonic/gin"

func Init51PwnRoute(router *gin.Engine) {
	xx1 := router.Group("/51pwn")
	InitSubDomainRoute(xx1)
}
