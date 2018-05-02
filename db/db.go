package db

import(
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"relation/conf"
	"time"
)


var db *gorm.DB

func InitMysql(cfg conf.ConfigType) (err error){
	mysql := cfg.DB.Mysql
	for i:=0;i<5;i++ {
		db,err = gorm.Open("mysql",mysql)
		db.DB().SetMaxIdleConns(10) //数据库连接池配置
		db.DB().SetMaxOpenConns(10)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	if db == nil {
		return
	}
	db.Close()
}