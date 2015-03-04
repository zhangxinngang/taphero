package user_add_attr

import (
	"database/sql"
	"github.com/0studio/databasetemplate"
	key "github.com/0studio/storage_key"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func setupTest() {

}

func getMockDB() (db *sql.DB) {
	db, _ = databasetemplate.NewDBInstance(databasetemplate.DBConfig{
		Host: "127.0.0.1",
		User: "th_dev",
		Pass: "th_devpass",
		Name: "test",
	}, false)

	return
}

func TestDBUserAddAttr(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	now := time.Now()
	userAddAttr := UserAddAttr{}
	userAddAttr.SetUin(Uin)

	store := InitDBUserAddStorage(getMockDB())

	ok := store.Add(&userAddAttr)
	assert.True(t, ok)

	userRet, ok := store.Get(userAddAttr.GetUin())
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), userAddAttr.GetUin())
	userAddAttr.SetEnergy(1)
	userAddAttr.SetEnergyTime(now)
	ok = store.Set(&userAddAttr)
	assert.True(t, ok)

	userRet, ok = store.Get(userAddAttr.GetUin())
	assert.True(t, ok)
	assert.Equal(t, userRet.GetUin(), userAddAttr.GetUin())
	assert.Equal(t, userRet.GetEnergy(), userAddAttr.GetEnergy())
	assert.True(t, userRet.GetEnergy() != 0)

}

var (
	UninitedTime time.Time
)

func TestDBUserAddAttrTime0(t *testing.T) {
	setupTest()
	Uin := key.KeyUint64(rand.Int63())
	userAddAttr := UserAddAttr{}
	userAddAttr.SetUin(Uin)

	store := InitDBUserAddStorage(getMockDB())

	ok := store.Add(&userAddAttr)
	assert.True(t, ok)
	userRet, ok := store.Get(userAddAttr.GetUin())
	assert.True(t, ok)
	assert.Equal(t, userRet.GetLastOffTime(), UninitedTime)

}
