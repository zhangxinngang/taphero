package design

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/0studio/databasetemplate"
	"github.com/0studio/scheduler"
	_ "github.com/go-sql-driver/mysql"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/utils"
)

var designResourceInstance DesignResource = DesignResource{flag: true}

type DesignResource struct {
	lock sync.Mutex
	flag bool

	databasetemplate.GenericDaoImpl

	bdesignTblVersiondb DesignTblVersionDB
	// bstageDB            BStageDB
}

/*
* 初始化策划数据
 */
func (this *DesignResource) initDesignDBPool() {
	dbConfig, ok := conf.GetDesignDBConfig()
	if !ok {
		panic("read [db_config] error")
	}
	var (
		designDBInstance *sql.DB
	)
	designDBInstance, ok = databasetemplate.NewDBInstance(dbConfig, false)
	if !ok {
		panic("create [_design_db_config] error")
	}
	dbTemplate := &databasetemplate.DatabaseTemplateImpl{designDBInstance}

	if !conf.IsModeTest() {
		this.closeDesignDbInstance()
	}
	this.GenericDaoImpl = databasetemplate.GenericDaoImpl{dbTemplate}
	this.initDataStruct()
}
func (this *DesignResource) initDataStruct() {
	designDBInstance := this.getDesignGenericDaoImpl()
	if !designDBInstance.IsOk() {
		return
	}
	this.bdesignTblVersiondb.Tablename = TableNameTblVersionDb
	this.bdesignTblVersiondb.GenericDaoImpl = designDBInstance
	// this.bstageDB.GenericDaoImpl = designDBInstance
}

const (
	SCHEMA_FILE_VERSION_ID = 2
	DATA_FILE_VERSION_ID   = 1
)

func (this *DesignResource) checkSchemaFileMd5(sqlSchemaFile string) (ok bool, sqlFileMd5 string) {
	tblVersionInDB := this.bdesignTblVersiondb.GetCurrentTblVersion(SCHEMA_FILE_VERSION_ID)
	sqlFileMd5 = utils.GetFileMd5(sqlSchemaFile)
	if sqlFileMd5 == tblVersionInDB { // sql 文件的md5值与库中记录的一致， 说明库中的数据是最新的不需要再次导入一次
		ok = true
		return
	}
	return
}

func (this *DesignResource) checkDataFileMd5(sqlDataFilePath string) (ok bool, sqlFileMd5 string) {
	tblVersionInDB := this.bdesignTblVersiondb.GetCurrentTblVersion(DATA_FILE_VERSION_ID)
	sqlFileMd5 = utils.GetFileMd5(sqlDataFilePath)
	if sqlFileMd5 == tblVersionInDB { // sql 文件的md5值与库中记录的一致， 说明库中的数据是最新的不需要再次导入一次
		ok = true
		return
	}
	return
}
func (this *DesignResource) createTables(sqlSchemaFilePath string, forceDropAndCreate bool) bool {
	designDBInstance := this.getDesignGenericDaoImpl()
	if !designDBInstance.IsOk() {
		return false
	}
	mode := conf.GetMode()
	this.bdesignTblVersiondb.CreateTable(mode)
	sameMd5, sqlFileMD5 := this.checkSchemaFileMd5(sqlSchemaFilePath)
	if sameMd5 && !forceDropAndCreate {
		return true
	}
	dbConfig, _ := conf.GetDesignDBConfig()

	//strUser := fmt.Sprintf("--user=%s --password=%s %s", dbUser, dbPasswd, dbName)
	strUser := fmt.Sprintf("--user=%s", dbConfig.User)
	strPass := fmt.Sprintf("--password=%s", dbConfig.Pass)
	strHost := fmt.Sprintf("--host=%s", dbConfig.Host)
	strPort := fmt.Sprintf("--port=%s", "3306")
	strDatabase := fmt.Sprintf("%s", dbConfig.Name)
	// cat /data/th/config/data.sql| mysql --user=th_dev --password=th_devpass  --host=localhost  --port=3306  th_dev_1_1_design
	cmd := exec.Command("mysql", strUser, strPass, strHost, strPort, strDatabase)
	file, err := os.Open(sqlSchemaFilePath)
	defer func() {
		if file != nil {
			file.Close()
			file = nil
		}
	}()
	if err != nil {
		log.Error("os open schema_*.sql error:", err, sqlSchemaFilePath)
		fmt.Printf("===========[ERROR]open schema_*.sql error %s=====================\n", sqlSchemaFilePath)
		return false
	}

	cmd.Stdin = file

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("run script error:", err, string(output))
		fmt.Printf("===========[ERROR]load schema_*.sql error %s=====================\n%s\n", sqlSchemaFilePath, string(output))
		return false
	}

	return this.bdesignTblVersiondb.UpdateTblVersion(SCHEMA_FILE_VERSION_ID, sqlFileMD5) && this.bdesignTblVersiondb.UpdateTblVersion(DATA_FILE_VERSION_ID, "")
}

func (this *DesignResource) loadSqlDataToDB(sqlDataFilePath string) bool {
	sameMD5, sqlFileMD5 := this.checkDataFileMd5(sqlDataFilePath) //sql 文件的md5值与库中记录的一致， 说明库中的数据是最新的不需要再次导入一次
	if sameMD5 {                                                  //
		return true
	}

	dbConfig, _ := conf.GetDesignDBConfig()

	//strUser := fmt.Sprintf("--user=%s --password=%s %s", dbUser, dbPasswd, dbName)
	strUser := fmt.Sprintf("--user=%s", dbConfig.User)
	strPass := fmt.Sprintf("--password=%s", dbConfig.Pass)
	strHost := fmt.Sprintf("--host=%s", dbConfig.Host)
	strPort := fmt.Sprintf("--port=%s", "3306")
	strDatabase := fmt.Sprintf("%s", dbConfig.Name)
	// cat /data/th/config/data.sql| mysql --user=th_dev --password=th_devpass  --host=localhost  --port=3306  th_dev_1_1_design
	cmd := exec.Command("mysql", strUser, strPass, strHost, strPort, strDatabase)
	file, err := os.Open(sqlDataFilePath)
	defer func() {
		if file != nil {
			file.Close()
			file = nil
		}
	}()
	if err != nil {
		log.Error("os open data.sql error:", err, sqlDataFilePath)
		fmt.Printf("===========[ERROR]open data.sql error %s=====================\n", sqlDataFilePath)
		return false
	}

	cmd.Stdin = file

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("run script error:", err, string(output))
		fmt.Printf("===========[ERROR]load data.sql error %s=====================\n%s\n", sqlDataFilePath, string(output))
		return false
	}

	return this.bdesignTblVersiondb.UpdateTblVersion(DATA_FILE_VERSION_ID, sqlFileMD5)
}

func (this *DesignResource) loadDbDataToGame() {
	// load data
	// 别忘了在clearSet() 加相应清除的语句
	// this.bstageDB.InitSet(this.flag)
	// 别忘了在clearSet() 加相应清除的语句
	if this.flag {
		this.flag = false
	} else {
		this.flag = true
	}

}

func (this *DesignResource) updateDataToGame() {
	this.loadDbDataToGame()
}

func (this *DesignResource) clearUnusedSet() {
	// this.bstageDB.ClearUnusedSet(this.flag)
	log.Info("clear unused design data set")

}
func (this *DesignResource) clearUnusedSetLater() {
	// 10min后再清理不用的set ,以避免 在切换的瞬间仍然有进程使用 将要被清理的set
	timerFun := func(scheduler *scheduler.Scheduler) {
		if scheduler == nil {
			return
		}
		getInstance().lock.Lock()
		defer getInstance().lock.Unlock()
		oldFlag := scheduler.ID.(bool)
		if this.flag == oldFlag {
			getInstance().clearUnusedSet()
		}
	}
	s := scheduler.InitScheduler(this.flag, 60, 0, timerFun)
	s.Start()
}
func (res *DesignResource) closeDesignDbInstance() {
	if !res.GenericDaoImpl.IsOk() {
		return
	}

	if res.GenericDaoImpl.DatabaseTemplate.IsConnOk() {
		res.GenericDaoImpl.DatabaseTemplate.Close()
	}
}

func (this *DesignResource) getDesignGenericDaoImpl() *databasetemplate.GenericDaoImpl {
	return &(this.GenericDaoImpl)
}

func getInstance() *DesignResource {
	return &designResourceInstance
}

func print_err(sql string, err error) bool {
	if err != nil {
		fmt.Println(err, sql)
		return true
	}
	return false
}
