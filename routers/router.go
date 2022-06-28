package routers

import (
	"ctp/controller"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitRouter(msqlDb *gorm.DB) (ginRouter *gin.Engine,err error){
	ginRouter = gin.Default()

	controller.Resister(ginRouter, msqlDb)
	return
}
