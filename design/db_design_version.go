package design

/*
记录  策划数据的版本md5， 信息， 以避免每次都load sql 到design db
*/

import (
	"database/sql"
	"fmt"
	"github.com/0studio/databasetemplate"
	_ "github.com/go-sql-driver/mysql"
)

type DesignTblVersionDB struct {
	Tablename string
	// *sql.DB
	*databasetemplate.GenericDaoImpl
}

func (db *DesignTblVersionDB) DropTable() {
	sql := fmt.Sprintf(`drop table if exists %s ;`, db.Tablename)
	err := db.GenericDaoImpl.DatabaseTemplate.Exec(sql)
	if err != nil {
		fmt.Println(db.Tablename, err)
	}

}

// 创建表
func (db *DesignTblVersionDB) CreateTable(mode string) {
	sql := fmt.Sprintf(
		`CREATE TABLE if not exists %s (
         id int  default 1,
         tbl_version char(32) default '',
         PRIMARY KEY (id)
        ) ENGINE=innodb DEFAULT CHARSET=utf8;`, db.Tablename)
	err := db.GenericDaoImpl.DatabaseTemplate.Exec(sql)
	if err != nil {
		fmt.Println(db.Tablename, err)
	}
}

// 加入一条新数据 ，或者 在原值基础上加Cnt
func (db *DesignTblVersionDB) UpdateTblVersion(id int, TblVersion string) bool {
	sql := fmt.Sprintf("insert into %s (id,tbl_version) values (?,?) on DUPLICATE key update tbl_version =?", db.Tablename)
	err := db.GenericDaoImpl.DatabaseTemplate.Exec(sql, id, TblVersion, TblVersion)
	if print_err(sql, err) {
		return false
	}
	return true
}
func (db *DesignTblVersionDB) GetCurrentTblVersion(id int) (tblVersion string) {
	sql := fmt.Sprintf("select tbl_version from %s where id=?", db.Tablename)
	verionInterface, err := db.GenericDaoImpl.DatabaseTemplate.Query(sql, db.mapRow, id)
	if err != nil {
		return ""
	}
	if verionInterface == nil {
		return ""
	}

	return verionInterface.(string)
}
func (db *DesignTblVersionDB) mapRow(resultSet *sql.Rows) (obj interface{}, err error) {
	var version string
	err = resultSet.Scan(
		&version,
	)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (db *DesignTblVersionDB) Clear() bool {
	sql := fmt.Sprintf("delete from %s ", db.Tablename)
	err := db.GenericDaoImpl.DatabaseTemplate.Exec(sql)
	if print_err(sql, err) {
		return false
	}
	return true
}
