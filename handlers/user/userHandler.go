package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	db2 "relation/db"
	"relation/logger"
	"relation/models/userDao"
	"relation/util"
	"strconv"
)

func GetUserById(c *gin.Context) {

	var user userDao.User
	var err error
	var cacheKey string
	sID := c.DefaultQuery("id", "0")
	id, _ := strconv.Atoi(sID)

	cacheKey = fmt.Sprintf("C_USER_%d", id)
	redisClient := db2.GetRedisClient()
	defer redisClient.Close()

	data, _ := redis.String(redisClient.Do("GET", cacheKey))

	if data == "" {
		db := db2.GetDB()
		user, err = userDao.GetUserById(id, db)
		if err != nil {
			c.JSON(http.StatusOK, util.SuccessResponse(nil))
			return
		}
		b, _ := json.Marshal(user)
		result, _ := redis.String(redisClient.Do("SETEX", cacheKey, 10, string(b)))
		fmt.Println("Set data to Redis Return " + result)
	} else {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			c.JSON(http.StatusOK, util.SuccessResponse(nil))
		} else {

		}
	}
	c.JSON(http.StatusOK, util.SuccessResponse(user))

}

func CreateUserById(c *gin.Context) {
	type request struct {
		UserName string `json:"user_name" form:"user_name"`
		Password string `json:"password" form:"password"`
		Phone    int64  `json:"phone" form:"phone"`
		Wx       string `json:"wx" form:"wx"`
		Email    string `json:"email" form:"email"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	md5Password, _ := util.Md5(req.Password)
	user := userDao.User{
		UserName: req.UserName,
		Password: md5Password,
		Phone:    req.Phone,
		Wx:       req.Wx,
		Email:    req.Email,
	}

	db := db2.GetDB()
	resultUser, err  := userDao.CreateUser(db,&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.SuccessResponse(resultUser))
}

func GetUsers(c *gin.Context) {

	db := db2.GetDB()
	var users []userDao.User
	var err error
	users, err = userDao.GetUsers(db)
	if err != nil {
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
	} else {
		c.JSON(http.StatusOK, util.SuccessResponse(users))
	}
}

func GetUserByPhone(c *gin.Context) {
	db := db2.GetDB()
	var user userDao.User
	var err error
	phone,err := strconv.Atoi(c.DefaultQuery("phone","0"))
	user, err = userDao.GetUserByPhone(db,phone)
	if err != nil {
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
	} else {
		c.JSON(http.StatusOK, util.SuccessResponse(user))
	}
}


// 根据分页获取用户数据
func GetUsersSlice(c *gin.Context) {

	db := db2.GetDB()
	var users []userDao.User
	var err error
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	logger.LogInfo("In GetUsers Slice Page:%d", page)

	if err != nil || page < 1 {
		logger.LogWarn("page is Error:%d", page)
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
		return
	}

	var (
		start int
		count int
	)

	count = 2 //todo
	start = (page - 1) * count
	users, err = userDao.GetUsersSlice(db, start, count)
	if err != nil {
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
	} else {
		c.JSON(http.StatusOK, util.SuccessResponse(users))
	}

}

func Login(c *gin.Context) {

	type request struct {
		Phone int `json:"phone" form:"phone" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		logger.LogError("Request Param Error req : %V",req)
		c.JSON(http.StatusBadRequest, util.FailResponse(102,"login fail",nil))
		return
	}

	db := db2.GetDB()
	var user userDao.User
	var err error


	user, err = userDao.GetUserByPhone(db,req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
	}

	md5Password,_ := util.Md5(req.Password)
	if md5Password != user.Password {
		c.JSON(http.StatusOK, util.FailResponse(102,"login fail",nil))
	} else {
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
	}

}