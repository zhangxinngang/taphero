package design

// import (
// 	"database/sql"
// 	"fmt"
// 	"github.com/0studio/databasetemplate"
// 	_ "github.com/go-sql-driver/mysql"
// 	"zerogame.info/taphero/entity"
// 	"zerogame.info/taphero/log"
// )

// type BStageDB struct {
// 	*databasetemplate.GenericDaoImpl
// 	Set  entity.BStageSet
// 	Set2 entity.BStageSet
// }

// /* 获取所有关卡 */
// func (db *BStageDB) InitSet(flag bool) {
// 	arr, _ := db.DatabaseTemplate.QueryArray(B_STAGE_SELECT, db.mapRow)

// 	var ele entity.BStage
// 	set := make(entity.BStageSet)
// 	for _, obj := range arr {
// 		ele = obj.(entity.BStage)
// 		set[ele.StageId] = ele
// 	}
// 	if len(set) == 0 {
// 		log.Error("[Err]:error load bstage data...........................")
// 		fmt.Println("[Err]:error load bstage data...........................")
// 	}

// 	if flag {
// 		db.Set = set
// 	} else {
// 		db.Set2 = set
// 	}

// 	return
// }

// func (this *BStageDB) mapRow(resultSet *sql.Rows) (obj interface{}, err error) {
// 	var iBStage entity.BStage
// 	// var _v int32
// 	// var _str string
// 	err = resultSet.Scan(
// 		&iBStage.StageId,
// 		&iBStage.StageType,
// 		&iBStage.NeedEnergy,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return iBStage, nil
// }
// func (db *BStageDB) ClearUnusedSet(flag bool) {
// 	if flag {
// 		db.Set = nil
// 	} else {
// 		db.Set2 = nil
// 	}
// }
