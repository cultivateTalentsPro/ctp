package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"ctp/model"
	"encoding/json"
)

func IsExistPhoneNumber(db *gorm.DB, phoneNumber string) bool {
	var reg model.RegisterParam
	exist := db.HasTable("register_params")
	if exist {
		fmt.Println(exist)
		db.Where("phone_number = ?", phoneNumber).First(&reg)
		if reg.Email != "" {
			return true
		}
	}
	return false
}

func Resister(ginRouter *gin.Engine, db *gorm.DB) {
	var reg model.RegisterParam
	ginRouter.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(body, &reg)

		if reg.Email == "" || reg.UserType == "" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "邮箱不能为空"})
			return
		}
		if len(reg.PhoneNumber) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须是11位"})
			return
		}
		if len(reg.Passwd) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}
		if reg.Passwd != reg.ConfirmPasswd {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "输入密码不一致"})
			return
		}
		if IsExistPhoneNumber(db, reg.PhoneNumber) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号已存在"})
			return
		}
		//timeStr:=time.Now().Format("2006-01-02 15:04:05")  //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
		//
		//fmt.Println(timeStr)    //打印结果：2017-04-11 13:24:04
		//创建用户
		newUser := model.RegisterParam{
			Email: reg.Email,
			Passwd: reg.Passwd,
			ConfirmPasswd: reg.ConfirmPasswd,
			PhoneNumber: reg.PhoneNumber,
			UserType: reg.UserType,
		}
		db.AutoMigrate(&reg)
		db.Create(&newUser)
		fmt.Println(newUser.ID)

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"mg":   "请求成功",
		})
	})
}
