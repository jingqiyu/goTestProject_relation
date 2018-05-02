package logger

import (
	"os"
	"log"
	"fmt"
)

type Level int32

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	Info *log.Logger
	Warn *log.Logger
	Error *log.Logger
)


func Init() {
	infoFile,err := os.OpenFile("logs/info.log",os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	if err != nil {
		fmt.Println(err)
	}
	warnFile,_ := os.OpenFile("logs/warn.log",os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	errorFile,_ := os.OpenFile("logs/error.log",os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	Info = log.New(infoFile,"",log.LstdFlags)
	Warn = log.New(warnFile,"",log.LstdFlags)
	Error = log.New(errorFile,"",log.LstdFlags)
}

func LogInfo(data ...interface{}) {
	Info.Println(data)
}

func LogWarn(data ...interface{}) {
	Warn.Println(data)
}

func LogError(data ...interface{}) {
	Error.Println(data)
}
