package model

import "github.com/jinzhu/gorm"

type RegisterParam struct {
	gorm.Model
	Email         string `json:"email"`
	Passwd        string `json:"passwd"`
	ConfirmPasswd string `json:"confirmPasswd"`
	PhoneNumber   string `json:"phoneNumber"`
	UserType      string `json:"userType"`
}
