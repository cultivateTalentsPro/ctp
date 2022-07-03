package routers

import (
	"ctp/controller"
	"ctp/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitRouter(msqlDb *gorm.DB) (ginRouter *gin.Engine,err error){
	ginRouter = gin.Default()
	ginRouter.Use(middleware.Cors())
	ginRouter.POST("/api/auth/register", controller.Resister)

	ginRouter.POST("/api/auth/login", controller.Login)
	ginRouter.GET("/api/auth/info", middleware.AuthMiddleWare() ,controller.Info)
	return
}
