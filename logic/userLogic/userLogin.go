package userLogic

import (
	"github.com/jinzhu/gorm"
	"relation/dto"
	"relation/models/userDao"
	"relation/util"
)

const(
	PAGE_COUNT = 2
)

func Login(db *gorm.DB,req *dto.ReqLogin) util.Response{

	var user userDao.User
	var err error

	user, err = userDao.GetUserByPhone(db,req.Phone)
	if err != nil {
		return util.SuccessResponse(nil)
	}

	md5Password,_ := util.Md5(req.Password)
	if md5Password != user.Password {
		return util.FailResponse(102,"login fail",nil)
	} else {
		return util.SuccessResponse(nil)
	}
}

func GetUsersByPage(db *gorm.DB,req *dto.ReqGetUserBySlice) util.Response{
	var (
		start int
		count int
	)
	count = PAGE_COUNT
	start = (req.Page - 1) * count
	users, err := userDao.GetUsersSlice(db, start, count)
	if err != nil {
		return util.SuccessResponse(nil)
	} else {
		return util.SuccessResponse(users)
	}
}