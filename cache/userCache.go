package cache

import (
	"fmt"
	"relation/db"
	"github.com/gomodule/redigo/redis"
	"relation/models/userDao"
	"encoding/json"
	"github.com/pkg/errors"
)

func GetUserById(userId int) (*userDao.User,error){
	var user *userDao.User
	cacheKey := fmt.Sprintf("C_USER_%d", userId)
	redisClient := db.GetRedisClient()
	defer redisClient.Close()
	data, _ := redis.String(redisClient.Do("GET", cacheKey))
	if data == "" {
		return nil,nil
	} else {
		err := json.Unmarshal([]byte(data), user)
		if err != nil {
			return user,nil
		} else {
			return nil,errors.WithStack(err)
		}
	}
}
