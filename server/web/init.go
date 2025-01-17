package web

import (
	"github.com/gin-gonic/gin"
	"smServer/db"
	"smServer/server/web/controller"
	"smServer/server/web/middlware"
)

func Init(router *gin.Engine) {
	db.TestDb()
	initRouter(router)
}

func initRouter(router *gin.Engine) {
	router.Use(middlware.Cors()) // 跨域中间件
	router.Any("/account/register", controller.DefaultAccountController.Register)
}
