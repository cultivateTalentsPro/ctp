package logger

import (
	"ctp/config"
	log "github.com/sirupsen/logrus"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"time"
)

var logger *log.Logger
/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	`WithMaxAge` 设置文件清理前的最长保存时间
	`WithRotationCount` 设置文件清理前最多保存的个数
*/
func LoggerInit() {
	logger = log.New()
	logdir := config.GetLogDir()

	writer, _ := rotatelogs.New(
		logdir+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logdir),
		rotatelogs.WithMaxAge(10 * time.Hour),
		rotatelogs.WithRotationTime(time.Duration(time.Hour)),
		)
	logger.SetOutput(writer)
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
}
func GetLogger() *log.Logger {
	return logger

}
func TestLog(log *log.Logger) {
	for {
		log.Info("hello, world!")
		time.Sleep(time.Duration(2) * time.Second)
	}
}