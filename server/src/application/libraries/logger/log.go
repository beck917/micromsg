package logger

import (
	"application/libraries/toml"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

var Mylog *log.Entry
var globalName string

func InitLog(name string) {
	// Log as JSON instead of the default ASCII formatter.
	if toml.GlobalTomlConfig.Log.Value != "" {
		log.SetFormatter(&log.JSONFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
		globalName = name
		// Output to stderr instead of stdout, could also be a file.
		logFormat := fmt.Sprintf(toml.GlobalTomlConfig.Log.Value, name, time.Now().Format("2006-01-02"))
		file := getFile(logFormat)
		log.SetOutput(file)
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{})
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
		// Output to stderr instead of stdout, could also be a file.
		//log.SetOutput(os.Stderr)
	}

	Mylog = MyLog()

	//每天重置下
	if toml.GlobalTomlConfig.Log.Value != "" {
		c := cron.New()
		c.AddFunc("1 0 0 * * *", func() {
			file := getFile(fmt.Sprintf(toml.GlobalTomlConfig.Log.Value, name, time.Now().Format("2006-01-02")))
			log.SetOutput(file)
		})
		c.Start()
	}
}

func getFile(filename string) *os.File {
	var f *os.File
	if checkFileIsExist(filename) { //如果文件存在
		f, _ = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0777) //打开文件
	} else {
		f, _ = os.Create(filename) //创建文件
	}
	return f
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Info(args ...interface{}) {
	Mylog.Info(args...)
}

func Warning(args ...interface{}) {
	Mylog.Warning(args...)
}

func Error(args ...interface{}) {
	Mylog.Error(args...)
}

func Debug(args ...interface{}) {
	Mylog.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Mylog.Debugf(format, args...)
}

func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}

func WithField(key string, value interface{}) *log.Entry {
	return log.WithField(key, value)
}

func MyLog() *log.Entry {
	return log.WithFields(log.Fields{
		"prama": "mylog",
	})
}
