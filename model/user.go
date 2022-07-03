package model

import "github.com/jinzhu/gorm"

type RegisterParam struct {
	gorm.Model
	Email         string `json:"email"`
	Passwd        string `json:"passwd"`
	PhoneNumber   string `json:"phoneNumber"`
	UserType      string `json:"userType"`
}
