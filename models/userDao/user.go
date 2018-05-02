package userDao

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

var _ fmt.Formatter //无用保留fmt
type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Phone    int64  `json:"phone"`
	Wx       string `json:"wx"`
	Email    string `json:"email"`
}

func (User) TableName() string {
	return "user_info"
}

func CreateUser(db *gorm.DB,u *User) (*User,error) {

	err := db.Create(&u).Error
	return u,err
}

func GetUserById(id int, db *gorm.DB) (User, error) {
	var user User
	err := db.First(&user, id).Error
	return user, err
}

func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Limit(10).Error
	return users, err
}

func GetUsersSlice(db *gorm.DB, start,count int)([]User, error){
	var users []User
	err := db.Order("id").Offset(start).Limit(count).Find(&users).Error
	return users, err
}

func GetUserByPhone(db *gorm.DB, phoneNum int) (User,error) {
	var user User
	err := db.Where("phone = ?", phoneNum).First(&user).Error
	return user, err
}
