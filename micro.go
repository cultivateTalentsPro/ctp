package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"time"

	//"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
)

var (
	//etcdCaPath       string = "F:/gotest/src/ca.pem"
	etcdReg registry.Registry
)

func Init() {
	//tlsInfo := transport.TLSInfo{
	//	TrustedCAFile: etcdCaPath,
	//}
	//tlsConfig, _ := tlsInfo.ClientConfig()
	etcdReg = etcd.NewRegistry(
		registry.Addrs("192.168.72.135:2379"),
		//registry.TLSConfig(tlsConfig),
		//etcd.Auth("root", "123456"),
	)
}

type RegisterParam struct {
	Email         string `json:"email"`
	Passwd        string `json:"passwd"`
	ConfirmPasswd string `json:"confirmPasswd"`
	PhoneNumber   string `json:"phoneNumber"`
	UserType      string `json:"userType"`
	RegTime		  string
}

func InitDbMysql() (*gorm.DB, error){
	driverName := "mysql"
	host := "192.168.72.135"
	port := "3306"
	user := "admin"
	pwd := "#HIKlzz123"
	charset := "utf8"
	database := "yyc"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
							user, pwd, host,port,database,charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Duration(28801)*time.Second)
	return db, nil
}

func IsExistPhoneNumber(db *gorm.DB, phoneNumber string) bool {
	var reg RegisterParam
	db.Where("phone_number = ?", phoneNumber).First(&reg)
	if reg.Email != "" {
		return true
	}
	return false
}

func Resister(ginRouter *gin.Engine, db *gorm.DB) {
	var reg RegisterParam
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
		timeStr:=time.Now().Format("2006-01-02 15:04:05")  //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法

		fmt.Println(timeStr)    //打印结果：2017-04-11 13:24:04
		//创建用户
		newUser := RegisterParam{
			Email: reg.Email,
			Passwd: reg.Passwd,
			ConfirmPasswd: reg.ConfirmPasswd,
			PhoneNumber: reg.PhoneNumber,
			UserType: reg.UserType,
			RegTime: timeStr,
		}
		db.AutoMigrate(&reg)
		db.Create(&newUser)

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"mg":   "请求成功",
		})
	})
}
func InitRouter(msqlDb *gorm.DB) (ginRouter *gin.Engine,err error){
	ginRouter = gin.Default()
	Resister(ginRouter, msqlDb)
	return
}

func main() {
	Init()
	msqlDb, err :=InitDbMysql()
	if err != nil{
		fmt.Println(err)
		return
	}
	router, _ := InitRouter(msqlDb)
	//micro.NewService(
	//	micro.RegisterHandler()
	//	)
	service := web.NewService(
		web.Name("api.miki.com.userserver"),
		web.Address(":8002"),
		web.Handler(router),
		web.Registry(etcdReg),
		web.Version("v1.0.1"),
	)
	// initialise flags
	service.Init()
	// start the service
	service.Run()
}
