package user_add_attr

import (
	"database/sql"
	"github.com/0studio/databasetemplate"
	key "github.com/0studio/storage_key"
	log "github.com/cihub/seelog"
	"github.com/dropbox/godropbox/errors"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type DBUserAddStorage struct {
	databasetemplate.GenericDaoImpl
}

var initUserAddDaoOnce sync.Once

func InitDBUserAddStorage(db *sql.DB) *DBUserAddStorage {
	var userAddDao *DBUserAddStorage
	dbTemplate := &databasetemplate.DatabaseTemplateImpl{db}
	userAddDao = &DBUserAddStorage{GenericDaoImpl: databasetemplate.GenericDaoImpl{dbTemplate}}
	userAddDao.CreateTable()
	return userAddDao
}

func (this *DBUserAddStorage) CreateTable() {
	sql := `create table if not exists player_add_attr (
                uin BIGINT NOT NULL,
                lastOffTime timestamp NOT NULL default 0,
                energy int NOT NULL default 0,
                energyTime timestamp NOT NULL default 0,
                PRIMARY KEY (uin))
             ENGINE = innodb DEFAULT CHARACTER SET utf8;`

	err := this.DatabaseTemplate.Exec(sql)
	if err != nil {
		log.Error(errors.Wrap(err, "create table player_add_attr error!!!"))
	}
}

func (this *DBUserAddStorage) mapRow(resultSet *sql.Rows) (interface{}, error) {
	var (
		uin         key.KeyUint64
		energy      int32
		energyTime  time.Time
		lastOffTime time.Time
	)
	userAdd := UserAddAttr{}
	err := resultSet.Scan(
		&uin,
		&lastOffTime,
		&energy,
		&energyTime)
	if err != nil {
		return nil, err
	}
	userAdd.SetUin(uin)
	userAdd.SetEnergyTime(energyTime)
	userAdd.SetEnergy(energy)
	userAdd.SetLastOffTime(lastOffTime)
	userAdd.ClearFlag()

	return userAdd, nil
}

func (this *DBUserAddStorage) Get(uin key.KeyUint64) (userAdd UserAddAttr, ok bool) {
	sql := "select uin,lastOffTime,energy,energyTime from player_add_attr where uin=?"
	var obj interface{}
	var err error
	obj, err = this.DatabaseTemplate.QueryObject(sql, this.mapRow, uin)
	if err != nil {
		return
	}
	if obj == nil {
		return
	}

	userAdd = obj.(UserAddAttr)
	ok = true
	return
}

func (this *DBUserAddStorage) Set(userAdd *UserAddAttr) bool {
	sql := "update player_add_attr set lastOffTime =?,energy=?,energyTime=? where Uin=?"
	err := this.DatabaseTemplate.Exec(sql, userAdd.GetLastOffTime(), userAdd.GetEnergy(), userAdd.GetEnergyTime(), userAdd.GetUin())
	return err == nil
}

func (this *DBUserAddStorage) Add(userAdd *UserAddAttr) bool {
	sql := "insert ignore into  player_add_attr (uin,lastOffTime,energy,energyTime) values (?,?,?,?)"
	err := this.DatabaseTemplate.Exec(sql,
		userAdd.GetUin(),
		userAdd.GetLastOffTime(),
		userAdd.GetEnergy(),
		userAdd.GetEnergyTime())
	return err == nil
}
