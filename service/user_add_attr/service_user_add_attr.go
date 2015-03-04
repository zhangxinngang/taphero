package user_add_attr

import (
	"database/sql"
	"github.com/0studio/storage_key"
	"time"
)

// User addition 附加属性
type UserAddAttrService interface {
	Get(uin key.KeyUint64, now time.Time) (userAddAttr UserAddAttr, ok bool)
	Set(userAddAttr *UserAddAttr) bool
}

func GetUserAddAttrService(db *sql.DB) UserAddAttrService {
	return &UserAddAttrServiceImpl{InitDBUserAddStorage(db)}
}

type UserAddAttrServiceImpl struct {
	dbUserAddAttrService *DBUserAddStorage
}

func (impl *UserAddAttrServiceImpl) Get(uin key.KeyUint64, now time.Time) (userAddAttr UserAddAttr, ok bool) {
	userAddAttr, ok = impl.dbUserAddAttrService.Get(uin)
	if !ok {
		userAddAttr = NewDefaultUserAddAttr(uin)
		ok = impl.dbUserAddAttrService.Add(&userAddAttr)
	}
	userAddAttr.SetMaxEnergy(100)
	userAddAttr.RecoveryEnergy(now)
	return
}
func (impl *UserAddAttrServiceImpl) Set(userAddAttr *UserAddAttr) bool {
	return impl.dbUserAddAttrService.Set(userAddAttr)
}
