/**
 * @Author: wzx
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/5/17 上午6:28
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package common

import (
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

type TixLogger struct {
	Hostname string
	*logrus.Logger
}

var Logger = MustGetLogger()

func (log *TixLogger) Init(config *Config) {
	log.setupLogging(config)
}

func (log *TixLogger) setupLogging(config *Config) {
	logrus.SetLevel(logrus.AllLevels[config.Level])
	if config.LogToFile == true {
		logPath := config.LogPath
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			log.Fatal("Cannot log to file", err.Error())
		}
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(file)
	}
	log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
}

func MustGetLogger() *TixLogger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	tixLogger := &TixLogger{hostname, logrus.StandardLogger()}
	return tixLogger
}

func (tx *TixLogger) Info(args ...interface{}) {
	fields(tx).Info(args...)
}

func (tx *TixLogger) Infof(format string, args ...interface{}) {
	fields(tx).Infof(format, args...)
}

func (tx *TixLogger) Debug(args ...interface{}) {
	fields(tx).Debug(args...)
}

func (tx *TixLogger) Debugf(format string, args ...interface{}) {
	fields(tx).Debugf(format, args...)
}

func (tx *TixLogger) Warn(args ...interface{}) {
	fields(tx).Warn(args...)
}

func (tx *TixLogger) Warnf(format string, args ...interface{}) {
	fields(tx).Warnf(format, args...)
}

func (tx *TixLogger) Error(args ...interface{}) {
	fields(tx).Error(args...)
}

func (tx *TixLogger) Errorf(format string, args ...interface{}) {
	fields(tx).Errorf(format, args...)
}

// DebugfWithId write formatted debug level log with added log_id field
func (tx *TixLogger) DebugfWithId(id string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).Debugf(format, args...)
}

// InfofWithId write formatted info level log with added log_id field
func (tx *TixLogger) InfofWithId(id string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).Infof(format, args...)
}

// InfoWithId write info level log with added log_id field
func (tx *TixLogger) InfoWithId(id string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).Info(args...)
}

// ErrorfWithId write formatted error level log with added log_id field
func (tx *TixLogger) ErrorfWithId(id string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).Errorf(format, args...)
}

// ErrorWithId write error level log with added log_id field
func (tx *TixLogger) ErrorWithId(id string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).Error(args...)
}

func (tx *TixLogger) DebugfWithIdAndConn(id string, connInfo string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("conn_info", connInfo).Debugf(format, args...)
}

func (tx *TixLogger) InfofWithIdAndConn(id string, connInfo string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("conn_info", connInfo).Infof(format, args...)
}

func (tx *TixLogger) InfoWithIdAndConn(id string, connInfo string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("conn_info", connInfo).Info(args...)
}

func (tx *TixLogger) ErrorfWithIdAndConn(id string, connInfo string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("conn_info", connInfo).Errorf(format, args...)
}

func (tx *TixLogger) ErrorWithIdAndConn(id string, connInfo string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("conn_info", connInfo).Error(args...)
}

func (tx *TixLogger) DebugWithConn(connInfo string, args ...interface{}) {
	fields(tx).WithField("conn_info", connInfo).Debug(args...)
}

func (tx *TixLogger) DebugfWithConn(connInfo string, format string, args ...interface{}) {
	fields(tx).WithField("conn_info", connInfo).Debugf(format, args...)
}

func (tx *TixLogger) InfofWithConn(connInfo string, format string, args ...interface{}) {
	fields(tx).WithField("conn_info", connInfo).Infof(format, args...)
}

func (tx *TixLogger) InfoWithConn(connInfo string, args ...interface{}) {
	fields(tx).WithField("conn_info", connInfo).Info(args...)
}

func (tx *TixLogger) ErrorfWithConn(connInfo string, format string, args ...interface{}) {
	fields(tx).WithField("conn_info", connInfo).Errorf(format, args...)
}

func (tx *TixLogger) ErrorWithConn(connInfo string, args ...interface{}) {
	fields(tx).WithField("conn_info", connInfo).Error(args...)
}

func (tx *TixLogger) DebugfWithIdAndConnAndTypeAndCommand(id string, connInfo string, Type string, Command string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("type", Type).WithField("command", Command).WithField("conn_info", connInfo).Debugf(format, args...)
}

func (tx *TixLogger) InfofWithIdAndConnAndTypeAndCommand(id string, connInfo string, Type string, Command string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("type", Type).WithField("command", Command).WithField("conn_info", connInfo).Infof(format, args...)
}

func (tx *TixLogger) InfoWithIdAndConnAndTypeAndCommand(id string, connInfo string, Type string, Command string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("type", Type).WithField("command", Command).WithField("conn_info", connInfo).Info(args...)
}

func (tx *TixLogger) ErrorfWithIdAndConnAndTypeAndCommand(id string, connInfo string, Type string, Command string, format string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("type", Type).WithField("command", Command).WithField("conn_info", connInfo).Errorf(format, args...)
}

func (tx *TixLogger) ErrorWithIdAndConnAndTypeAndCommand(id string, connInfo string, Type string, Command string, args ...interface{}) {
	fields(tx).WithField("conn_id", "id:"+string(id)).WithField("type", Type).WithField("command", Command).WithField("conn_info", connInfo).Error(args...)
}

func fields(tx *TixLogger) *logrus.Entry {
	file, line, funcName := findCaller(4)
	return tx.Logger.WithField("file", file).WithField("fn", funcName).WithField("line", line)
}

func findCaller(skip int) (string, int, string) {
	file := ""
	line := 0
	var pc uintptr
	// 遍历调用栈的最大索引为第11层.
	for i := 0; i < 11; i++ {
		file, line, pc = getCaller(skip + i)
		// 过滤掉所有logrus包，即可得到生成代码信息
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}

	fullFnName := runtime.FuncForPC(pc)

	fnName := ""
	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		// 取得函数名
		parts := strings.Split(fnNameStr, ".")
		fnName = parts[len(parts)-1]
	}

	return file, line, fnName
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0

	// 获取包名
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}
