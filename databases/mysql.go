package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type MysqlParam struct {
	DriverName string `ini:"driver_name"`
	Host       string `ini:"host"`
	Port       string `ini:"port"`
	UserName   string `ini:"user_name"`
	Passwd     string `ini:"passwd"`
	Charset    string `ini:"charset"`
	Database   string `ini:"database"`
}

func InitDbMysql(mp MysqlParam) (*gorm.DB, error) {
	fmt.Println(mp.UserName, mp.Passwd, mp.Host, mp.Port, mp.Database, mp.Charset)
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		mp.UserName, mp.Passwd, mp.Host, mp.Port, mp.Database, mp.Charset)
	fmt.Println(args)
	db, err := gorm.Open(mp.DriverName, args)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Duration(28801) * time.Second)
	return db, nil
}
