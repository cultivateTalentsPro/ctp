package controller

import (
	"ctp/databases"
	"ctp/middleware"
	"ctp/model"
	"ctp/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"ctp/logger"
	//"go-micro.dev/v4/services/db"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

var (
	log *logrus.Logger
)

type LoginParam struct {
	UserName string
	Passwd   string
}

type VueUserInfo struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	UserType    string `json:"userType"`
}

func IsExistUser(db *gorm.DB, reg model.RegisterParam) bool {
	exist := db.HasTable("register_params")
	if exist {
		fmt.Println(exist)
		//db.Where("phone_number = ?", phoneNumber).First(&reg)
		db.Find(&reg, "phone_number = ? OR email = ?", reg.PhoneNumber, reg.Email)
		if reg.ID != 0 {
			return true
		}
	}
	return false
}

func GetBodyContent(ctx *gin.Context, con interface{}) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	err = json.Unmarshal(body, &con)
	if err != nil {
		log.Error(err)
		return err
	}
	return err
}

func Resister(ctx *gin.Context) {
	var reg model.RegisterParam

	if log == nil {
		log = logger.GetLogger()
	}
	db := databases.GetMysqlDB()
	//获取参数
	err := GetBodyContent(ctx, &reg)
	if err != nil {
		log.Error(err)
		return
	}

	//ctx.Bind(&reg)
	fmt.Println(reg)
	if reg.Email == "" || reg.UserType == "" {
		log.Error("111111111")
		response.Response(ctx,http.StatusOK, 422,nil, "邮箱不能为空")
		return
	}
	if len(reg.PhoneNumber) != 11 {
		log.Error("2222222")
		response.Response(ctx,http.StatusOK, 422,nil, "手机号必须是11位")
		return
	}
	if len(reg.Passwd) < 6 {
		response.Response(ctx,http.StatusOK, 422,nil, "密码不能少于6位")
		return
	}

	if IsExistUser(db, reg) {
		response.Response(ctx,http.StatusOK, 422,nil, "用户已存在")
		return
	}
	//timeStr:=time.Now().Format("2006-01-02 15:04:05")  //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	//
	//fmt.Println(timeStr)    //打印结果：2017-04-11 13:24:04
	//创建用户
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(reg.Passwd), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx,http.StatusOK, 422,nil, "加密错误")
		return
	}
	newUser := model.RegisterParam{
		Email:       reg.Email,
		Passwd:      string(hashPasswd),
		PhoneNumber: reg.PhoneNumber,
		UserType:    reg.UserType,
	}
	db.AutoMigrate(&reg)
	db.Create(&newUser)
	//fmt.Println(newUser.ID)
	reg.ID = newUser.ID

	token, err := middleware.ReleaseToken(reg)
	if err != nil {
		response.Response(ctx,http.StatusOK, 422,nil, "token err")
		return
	}
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func CheckUser(db *gorm.DB, lp LoginParam, ctx *gin.Context) error {
	var regUser model.RegisterParam
	// check user
	db.Find(&regUser, "phone_number = ? OR email = ?", lp.UserName, lp.UserName)
	if regUser.ID == 0 {
		response.Response(ctx,http.StatusOK, 422,nil, "用户不存在")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在", "data": ""})
		return errors.New("用户不存在")
	}
	//check passwd
	err := bcrypt.CompareHashAndPassword([]byte(regUser.Passwd), []byte(lp.Passwd))
	if err != nil {
		response.Response(ctx,http.StatusOK, 422,nil, "密码错误")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码错误", "data": ""})
		return err
	}
	token, err := middleware.ReleaseToken(regUser)
	if err != nil {
		response.Response(ctx,http.StatusOK, 422,nil, "token err")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "token err", "data": ""})
		return err
	}
	//db.Where("phone_number = ?", lp.UserName).First(&regUser)
	//if regUser.ID == 0 {
	//	db.Where("email = ?", lp.UserName).First(&regUser)
	//}
	response.Success(ctx, gin.H{"token": token}, "登录成功")
	return nil
}

func Login(ctx *gin.Context) {
	var loginParam LoginParam
	if log == nil {
		log = logger.GetLogger()
	}
	db := databases.GetMysqlDB()
	//get param
	err := GetBodyContent(ctx, &loginParam)
	if err != nil {
		log.Error(err)
		return
	}
	// check param
	//if len(loginParam.UserName) != 11 {
	//	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须是11位"})
	//	return
	//}
	if len(loginParam.Passwd) < 6 {
		response.Response(ctx,http.StatusOK, 422,nil, "密码不能少于6位")
		return
	}
	//check username and pwd
	err = CheckUser(db, loginParam, ctx)
	if err != nil {
		log.Error(err)
		return
	}
}

func ToVueUser(reg model.RegisterParam) VueUserInfo {
	return VueUserInfo{
		Email: reg.Email,
		PhoneNumber: reg.PhoneNumber,
		UserType: reg.UserType,
	}
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": ToVueUser(user.(model.RegisterParam))}, "获取成功")
}
