package conf

type ConfigType struct {
	DB struct {
		Mysql string
		Redis struct{
			Port int
			Auth string
			Timeout int
			MaxIdle int
		}
	}

}

var Config ConfigType


