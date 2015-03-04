package design

import (
	"fmt"

	"zerogame.info/taphero/conf"
	// "zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
)

func SetUp() bool {
	// 同一服务器 不同进程间  避免同时create table truncate table 造成冲突
	// memcache.LockDesignDataLoader()
	// defer memcache.UnlockDesignDataLoader()
	getInstance().lock.Lock()
	defer getInstance().lock.Unlock()

	getInstance().initDesignDBPool()
	if !conf.IsModeTest() {
		defer getInstance().closeDesignDbInstance()
	}

	getInstance().createTables(conf.GetDesignSchemaSqlFile(), false)
	if getInstance().loadSqlDataToDB(conf.GetDesignDataSqlFile()) {
		getInstance().loadDbDataToGame()
		return true
	}
	return false
}

func LoadSqlToDB() {
	fmt.Println("reloading design sql file to db...!! ")
	log.Info("reloading design sql file to db...!! ")

	// 同一服务器 不同进程间  避免同时create table truncate table 造成冲突
	// memcache.LockDesignDataLoader()
	// defer memcache.UnlockDesignDataLoader()

	getInstance().lock.Lock()
	defer getInstance().lock.Unlock()

	getInstance().initDesignDBPool()
	defer getInstance().closeDesignDbInstance()

	getInstance().createTables(conf.GetDesignSchemaSqlFile(), false)
	if getInstance().loadSqlDataToDB(conf.GetDesignDataSqlFile()) {
	}
	fmt.Println("reloading design sql file to db [done]...!! ")
	log.Info("reloading design sql file to db [done]...!! ")
}

func ReloadFromSqlFile() {
	fmt.Println("reloading design from sql...!! ")
	log.Info("reloading design from sql !!! ")

	// 同一服务器 不同进程间  避免同时create table truncate table 造成冲突
	// memcache.LockDesignDataLoader()
	// defer memcache.UnlockDesignDataLoader()

	getInstance().lock.Lock()
	defer getInstance().lock.Unlock()

	getInstance().initDesignDBPool()
	defer getInstance().closeDesignDbInstance()

	getInstance().createTables(conf.GetDesignSchemaSqlFile(), false)
	if getInstance().loadSqlDataToDB(conf.GetDesignDataSqlFile()) {
		getInstance().loadDbDataToGame()
		getInstance().clearUnusedSetLater() //
	}
	fmt.Println("reload design from sql done...!! ")
	log.Info("reload design from sql done!! ")
}

func ReloadFromDB() {
	// 同一服务器 不同进程间  避免同时create table truncate table 造成冲突
	fmt.Println("reload design from db to game...!! ")
	log.Info("reload design from db to game ...!! ")
	// memcache.LockDesignDataLoader()
	// defer memcache.UnlockDesignDataLoader()

	getInstance().lock.Lock()
	defer getInstance().lock.Unlock()

	getInstance().initDesignDBPool()
	defer getInstance().closeDesignDbInstance()
	getInstance().loadDbDataToGame()
	getInstance().clearUnusedSetLater() //
	fmt.Println("reload design from db to game done...!! ")
	log.Info("reload design from db to game done...!! ")
}

func DoShutDown() {
	log.Info("shutdowning before  design.DoShutdown ...")
	getInstance().closeDesignDbInstance()
	log.Info("shutdowning after  design.DoShutdown ...")
}

////////////////////////////////////////////////////////////////////////////////
// get b set
////////////////////////////////////////////////////////////////////////////////

// func GetBStageSet() entity.BStageSet {
// 	if getInstance().flag {
// 		return getInstance().bstageDB.Set2
// 	} else {
// 		return getInstance().bstageDB.Set
// 	}
// }
