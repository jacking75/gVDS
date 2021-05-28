package scommon

import (
	"fmt"
	"os"
)

const (
	LOG_LEVEL_TRACE = 0
	LOG_LEVEL_DEBUG = 1
	LOG_LEVEL_INFO = 2
	LOG_LEVEL_WARN = 3
	LOG_LEVEL_ERROR = 4
	LOG_LEVEL_FATAL = 5
)

var logLevelStr = [6]string {"trace", "debug", "info", "warn", "error", "fatal"}

var (
	OutPutLog = _emptyExportLog
)

func InitLog(level int, logFunc func(int, string)) {
	_logLevel = level

	if logFunc != nil {
		OutPutLog = logFunc
	}
}

func LogTrace(msg string) {
	OutPutLog(LOG_LEVEL_TRACE, msg)
}
func LogDebug(msg string) {
	OutPutLog(LOG_LEVEL_DEBUG, msg)
}
func LogInfo(msg string) {
	OutPutLog(LOG_LEVEL_INFO, msg)
}
func LogError(msg string) {
	OutPutLog(LOG_LEVEL_ERROR, msg)
}




// 비공개 함수
func _emptyExportLog(level int, msg string) {
	if level < _logLevel {
		return
	}

	fmt.Fprintf(os.Stdout,"[ %s ] %s\n", logLevelStr[level], msg)
}

var _logLevel int = LOG_LEVEL_DEBUG