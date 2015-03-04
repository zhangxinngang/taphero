package conf

import (
	"fmt"
	"github.com/0studio/databasetemplate"
	"github.com/Unknwon/goconfig"
	"strconv"
	"strings"
	"zerogame.info/taphero/utils"
)

var mode string
var server uint64
var process uint64
var configDir string
var configFile string
var appConfig *goconfig.ConfigFile

func GetConfigDir() string {
	return configDir
}

func GetMode() string {
	return mode
}
func IsModeTest() bool {
	return GetMode() == "test"
}
func IsModePro() bool {
	return GetMode() == "pro"
}
func IsModeDev() bool {
	return GetMode() == "dev"
}
func GetPlatform() uint64 {

	return 1
}

var locale string

func GetLocale() string {
	return locale
}
func GetProcess() uint64 {
	return process
}
func GetServer() uint64 {
	return server
}

func GetMCConfig() (mcConfig string, maxActiveConn int32, maxIdleConn uint32, ok bool) {
	var err error
	mcConfig, err = appConfig.GetValue(CONFIG_HEADER_TAPHERO, CONFIG_TAPHERO_KEY_MEMCACHE_ADDR)
	if err != nil {
		return
	}
	var tmpStr string
	var tmpInt int
	tmpStr, err = appConfig.GetValue(CONFIG_HEADER_TAPHERO, CONFIG_TAPHERO_KEY_MEMCACHE_MAX_ACTIVE_CONNECTIONS)
	if err != nil {
		return
	}

	tmpInt, err = strconv.Atoi(tmpStr)
	if err != nil {
		return
	}
	maxActiveConn = int32(tmpInt)

	tmpStr, err = appConfig.GetValue(CONFIG_HEADER_TAPHERO, CONFIG_TAPHERO_KEY_MEMCACHE_MAX_IDLE_CONNECTIONS)
	if err != nil {
		return
	}

	tmpInt, err = strconv.Atoi(tmpStr)
	if err != nil {
		return
	}
	maxIdleConn = uint32(tmpInt)

	ok = true
	return
}
func GetServerRedisConfig() (addr string, maxActive, maxIdle int) {
	addr, _ = appConfig.GetValue(CONFIG_HEADER_TAPHERO, CONFIG_SERVER_REDIS)
	maxActiveStr, _ := appConfig.GetValue(CONFIG_HEADER_TAPHERO, CONFIG_SERVER_REDIS_MAX_ACIGVE)
	maxActive, _ = strconv.Atoi(maxActiveStr)
	maxIdleStr, _ := appConfig.GetValue(CONFIG_HEADER_TAPHERO, CONFIG_SERVER_REDIS_MAX_IDLE)
	maxIdle, _ = strconv.Atoi(maxIdleStr)
	return
}

func GetDBConfig() (dbConfig databasetemplate.DBConfig, ok bool) {
	var err error
	var err2 error

	dbConfig.User, err = appConfig.GetValue(CONFIG_HEADER_DB, CONFIG_DB_KEY_USER)
	if err != nil {
		ok = false
		return
	}
	dbConfig.User = strings.TrimSpace(dbConfig.User)

	dbConfig.Pass, err = appConfig.GetValue(CONFIG_HEADER_DB, CONFIG_DB_KEY_PASSWD)
	if err != nil {
		ok = false
		return
	}
	dbConfig.Pass = strings.TrimSpace(dbConfig.Pass)

	dbConfig.Name, err = appConfig.GetValue(CONFIG_HEADER_DB, CONFIG_DB_KEY_DATABASE)
	if err != nil {
		ok = false
		return
	}
	dbConfig.Name = strings.TrimSpace(dbConfig.Name)

	dbConfig.Host, err = appConfig.GetValue(CONFIG_HEADER_DB, CONFIG_DB_KEY_HOST)
	dbConfig.Host = strings.TrimSpace(dbConfig.Host)
	if err != nil && err2 != nil {
		ok = false
		return
	}
	ok = true

	return
}

func GetDesignDBConfig() (dbConfig databasetemplate.DBConfig, ok bool) {
	var err error
	var err2 error

	dbConfig.User, err = appConfig.GetValue(CONFIG_HEADER_DESIGN_DB, CONFIG_DB_KEY_USER)
	if err != nil {
		ok = false
		return
	}
	dbConfig.User = strings.TrimSpace(dbConfig.User)

	dbConfig.Pass, err = appConfig.GetValue(CONFIG_HEADER_DESIGN_DB, CONFIG_DB_KEY_PASSWD)
	if err != nil {
		ok = false
		return
	}
	dbConfig.Pass = strings.TrimSpace(dbConfig.Pass)

	dbConfig.Name, err = appConfig.GetValue(CONFIG_HEADER_DESIGN_DB, CONFIG_DB_KEY_DATABASE)
	if err != nil {
		ok = false
		return
	}
	dbConfig.Name = strings.TrimSpace(dbConfig.Name)

	dbConfig.Host, err = appConfig.GetValue(CONFIG_HEADER_DESIGN_DB, CONFIG_DB_KEY_HOST)
	dbConfig.Host = strings.TrimSpace(dbConfig.Host)
	if err != nil && err2 != nil {
		ok = false
		return
	}
	ok = true

	return
}
func GetDesignDataSqlFile() string {
	file := fmt.Sprintf("%sdata_%s.sql", GetConfigDir(), GetLocale())
	if !utils.IsFileExists(file) {
		return ""
	}
	return file
}

func GetDesignSchemaSqlFile() string {
	file := fmt.Sprintf("%sschema_%s.sql", GetConfigDir(), GetLocale())
	if !utils.IsFileExists(file) {
		return ""
	}
	return file
}
func GetDesignDataTestSqlFile() string {
	file := fmt.Sprintf("%stest_schema_%s.sql", GetConfigDir(), GetLocale())
	if !utils.IsFileExists(file) {
		return ""
	}
	return file
}
