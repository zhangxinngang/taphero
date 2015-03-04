package service

import (
	"github.com/0studio/databasetemplate"
	storage "github.com/0studio/mcstorage"
	"zerogame.info/profile"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/service/user_add_attr"
)

type ServiceManager struct {
	UserService        profile.UserService
	PayOrderService    profile.PayOrderService
	UserAddAttrService user_add_attr.UserAddAttrService
}

var ServiceManagerInstance ServiceManager

func Init() {
	initProfile()
	initUserAddAttr()
}
func GetUserService() profile.UserService {
	return ServiceManagerInstance.UserService
}
func GetPayOrderService() profile.PayOrderService {
	return ServiceManagerInstance.PayOrderService
}

func GetUserAddAttrService() user_add_attr.UserAddAttrService {
	return ServiceManagerInstance.UserAddAttrService
}

func initProfile() {
	dbConfig, ok := conf.GetDBConfig()
	if !ok {
		panic("read db config failed")
		return
	}
	sqlDB, ok := databasetemplate.NewDBInstance(dbConfig, true)
	if !ok {
		panic("conn_db_error")
	}
	ServiceManagerInstance.PayOrderService = profile.InitDBPayOrderStorage(sqlDB)

	mcConfig, mcMaxActiveConns, mcMaxIdleActiveConns, ok := conf.GetMCConfig()
	if !ok {
		panic("read memcache config failed")
		return
	}

	mcClient := storage.GetClient(mcConfig, mcMaxActiveConns, mcMaxIdleActiveConns, log.LogError, log.Info)
	ServiceManagerInstance.UserService = profile.GetUserService(sqlDB, mcClient, conf.GetPlatform(), conf.GetServer(), conf.GetProcess())
}
func initUserAddAttr() {
	dbConfig, ok := conf.GetDBConfig()
	if !ok {
		panic("read db config failed")
		return
	}
	sqlDB, ok := databasetemplate.NewDBInstance(dbConfig, true)
	if !ok {
		panic("conn_db_error")
	}
	ServiceManagerInstance.UserAddAttrService = user_add_attr.GetUserAddAttrService(sqlDB)

}
