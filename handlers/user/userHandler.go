package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	db2 "relation/db"
	"relation/dto"
	"relation/logger"
	"relation/logic/userLogic"
	"relation/models/userDao"
	"relation/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"net/url"
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
	resultUser, err := userDao.CreateUser(db, &user)
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
	phone, err := strconv.Atoi(c.DefaultQuery("phone", "0"))
	user, err = userDao.GetUserByPhone(db, phone)
	if err != nil {
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
	} else {
		c.JSON(http.StatusOK, util.SuccessResponse(user))
	}
}

// 根据分页获取用户数据
func GetUsersSlice(c *gin.Context) {

	log := logger.GetLogger()
	db := db2.GetDB()
	var err error
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil || page < 1 {
		logger.LogWarn("page is Error:%d", page)
		c.JSON(http.StatusOK, util.SuccessResponse(nil))
		return
	}

	dtoReq := dto.ReqGetUserBySlice{page}
	log.Info(dtoReq)
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, userLogic.GetUsersByPage(db, &dtoReq))

}

func Login(c *gin.Context) {

	type request struct {
		Phone    int    `json:"phone" form:"phone" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	var req request
	if err := c.ShouldBind(&req); err != nil {
		logger.LogError("Request Param Error req : %V", req)
		c.JSON(http.StatusBadRequest, util.FailResponse(102, "login fail", nil))
		return
	}

	log := logger.GetLogger()
	cLog := log.WithFields(logrus.Fields{"Handler": "user"}) //定制化log
	cLog.Info("RequestParam:", req)
	db := db2.GetDB()

	dtoLogin := dto.ReqLogin{req.Phone, req.Password}
	response := userLogic.Login(db, &dtoLogin)
	if !isSuccess(response) {
		c.JSON(http.StatusBadRequest, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func isSuccess(response util.Response) bool {
	return response.ErrNo == util.SUCCESS
}

func ToCheck(c *gin.Context) {
	buf := make([]byte,4096)
	n,_ := c.Request.Body.Read(buf)
	log := logger.GetLogger()
	cLog := log.WithFields(logrus.Fields{"Handler": "user"}) //定制化log
	s,_ := url.QueryUnescape(string(buf[0:n]))

	cLog.Infof("RequestParam:%s", s)
	c.JSON(http.StatusOK,buf[0:n])
}
