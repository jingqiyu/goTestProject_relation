package logger

import(
	"github.com/rifflock/lfshook"
	"github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/pkg/errors"
	"fmt"
)


var Log *log.Logger

func GetLogger() *log.Logger {
	if Log != nil {
		return Log
	} else {
		return InitLocalFileSystemLogger(7*24*time.Hour,time.Hour)
	}
}

func InitLocalFileSystemLogger(maxAge,rotationTime time.Duration ) *log.Logger{
	if Log != nil {
		fmt.Println("Get log from singleton")
		return Log
	}
	Log = log.New()
	InfoLogPath := "logs/info1.log"
	WarnLogPath := "logs/wf.log"
	infoWriter,err := rotatelogs.New(
		InfoLogPath + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(InfoLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		LogWarn("config local file system logger error. %+v", errors.WithStack(err))
	}
	warnWriter,err := rotatelogs.New(
		WarnLogPath + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(InfoLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		LogWarn("config local file system logger error. %+v", errors.WithStack(err))
	}


	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: infoWriter, // 为不同级别设置不同的输出目的
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  warnWriter,
		log.ErrorLevel: warnWriter,
		log.FatalLevel: warnWriter,
		log.PanicLevel: warnWriter,
	},&log.JSONFormatter{})
	Log.AddHook(lfHook)
	return Log
}