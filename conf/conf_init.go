package conf

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"os"
	"runtime"
	"zerogame.info/thserver/utils"
)

func Init(_mode string, _server, _process uint64, _locale string) {
	mode = _mode
	server = _server
	process = _process
	locale = _locale
	// 设置起多少个操作系统线程来处理
	runtime.GOMAXPROCS(4 * runtime.NumCPU())
	if !dosetupConfigDir() {
		panic("get_configdir_error")
	}
	if !dosetupServerConfig() {
		panic("get_config_file_error")
	}

}

func dosetupConfigDir() bool {
	if os.Getenv("HOME") != "" {
		configDir = os.Getenv("HOME") + "/.taphero/"
	}
	if !utils.IsFileExists(configDir) {
		configDir = "/data/taphero/config/"
		if !utils.IsFileExists(configDir) {
			configDir = "./config/"
			if !utils.IsFileExists(configDir) {
				configDir = "../config/"
				if !utils.IsFileExists(configDir) {
					fmt.Println("app config dir doesnot exists")
					return false
				}
			}

		}

	}
	if !IsModeTest() {
		fmt.Printf("read config from dir %s\n", configDir)
	}

	return true
}
func dosetupServerConfig() bool {

	configFile = fmt.Sprintf("%s%s_%d_%d.ini", configDir, GetMode(), GetPlatform(), GetServer())
	if !utils.IsFileExists(configFile) {
		configFile = fmt.Sprintf("%s%s_%d_%d.ini", configDir, GetMode(), GetPlatform(), GetServer())
		if IsModePro() {
			fmt.Printf("app config file %s does NOT exists\n", configFile)
			return false
		}
		configFile = fmt.Sprintf("%s%s.ini", configDir, GetMode())

		if !utils.IsFileExists(configFile) {
			fmt.Printf("app config file %s does NOT exists\n", configFile)
			return false
		}
	}

	if !IsModeTest() {
		fmt.Printf("read config file %s\n", configFile)
	}
	// 加载配置
	config, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		fmt.Println("load config file %s fail", configFile)
		return false
	}
	appConfig = config
	return true
}
