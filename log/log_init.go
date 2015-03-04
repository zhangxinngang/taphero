package log

import (
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"zerogame.info/taphero/conf"
)

// 初始化logger
var logger seelog.LoggerInterface

func Init() {
	logPath := fmt.Sprintf("%sseelog_%s.xml", conf.GetConfigDir(), conf.GetMode())
	byteContent, err := ioutil.ReadFile(logPath)
	if err != nil {
		fmt.Println("read seelog config errr:", err)
		return
	}
	content := strings.Replace(string(byteContent), "%platformid", strconv.FormatUint(conf.GetPlatform(), 10), -1)
	content = strings.Replace(content, "%serverid", strconv.FormatUint(conf.GetServer(), 10), -1)
	content = strings.Replace(content, "%processidx", strconv.FormatUint(conf.GetProcess(), 10), -1)
	// content = strings.Replace(content, "%server_type", conf.GetServerType(), -1)

	logger, err = seelog.LoggerFromConfigAsString(content)
	// logger, err := seelog.LoggerFromConfigAsFile(logPath)
	if err != nil {
		fmt.Println(err)
		// os.Exit(defs.EXIT_CODE_SEELOG_INIT_FAIL)
	}

	seelog.ReplaceLogger(logger)
	defer seelog.Flush()
}
func ReloadLogger() {
	Close()
	Init()
}
func Close() {
	if logger != nil && !logger.Closed() {
		seelog.Flush()
		logger.Close()
		logger = nil
	}
}

func Debug(v ...interface{}) {
	// 如果debug 的日志太多，会导致如下提示
	// Seelog queue overflow: more than 10000 messages in the queue. Flushing.

	if conf.IsModeDev() {
		// 虽然可以在seelog.xml 里配置 Debug的内容不打印， 但是 似乎 seelog  在处理的时候
		seelog.Debug(GetPathLine(), v)
	}
}
func Debugf(format string, params ...interface{}) {
	if conf.IsModeDev() {
		seelog.Debugf(fmt.Sprintf("%s:%s", GetPathLine(), format), params...)
	}
}

func LogError(err error) {
	Error(err)
}
func Error(v ...interface{}) {
	if !conf.IsModePro() {
		fmt.Println(GetPathLine(), v)
	}
	seelog.Error(GetPathLine(), v)
}
func Flush() {
	seelog.Flush()
}

func Info(v ...interface{}) {
	seelog.Info(v)
}
func Warn(v ...interface{}) {
	seelog.Warn(v)
}
func Warnf(format string, params ...interface{}) {
	seelog.Warnf(format, params...)
}

func Infof(format string, params ...interface{}) {
	seelog.Infof(format, params...)
}
func Errorf(format string, params ...interface{}) {
	seelog.Errorf(fmt.Sprintf("%s:%s", GetPathLine(), format), params...)
}
func GetPathLine() string {
	_, path, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d", path, line)
	}
	return ""

}
