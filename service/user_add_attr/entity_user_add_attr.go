package user_add_attr

import (
	key "github.com/0studio/storage_key"
	"time"
)

type UserAddAttr struct {
	uin         key.KeyUint64
	energy      int32
	energyTime  time.Time
	lastOffTime time.Time
	maxEnergy   int32
}

func (this *UserAddAttr) SetUin(value key.KeyUint64) {
	this.uin = value
}
func (this UserAddAttr) GetUin() key.KeyUint64 {
	return this.uin
}
func (this *UserAddAttr) SetEnergy(value int32) {
	this.energy = value
}
func (this UserAddAttr) GetEnergy() int32 {
	return this.energy
}
func (this *UserAddAttr) SetEnergyTime(value time.Time) {
	this.energyTime = value
}
func (this UserAddAttr) GetEnergyTime() time.Time {
	return this.energyTime
}
func (this *UserAddAttr) SetLastOffTime(value time.Time) {
	this.lastOffTime = value
}
func (this UserAddAttr) GetLastOffTime() time.Time {
	return this.lastOffTime
}
func (this *UserAddAttr) SetMaxEnergy(value int32) {
	this.maxEnergy = value
}
func (this UserAddAttr) GetMaxEnergy() int32 {
	return this.maxEnergy
}

func (this *UserAddAttr) ClearFlag() {
}
func NewDefaultUserAddAttr(uin key.KeyUint64) (attr UserAddAttr) {
	attr.uin = uin
	now := time.Now()
	attr.energyTime = now
	attr.lastOffTime = now
	return
}

// 返回  下一次体力恢复 时间  多少秒
func (this *UserAddAttr) GetNextRecoveryEnergy(now time.Time) int32 {
	this.RecoveryEnergy(now)
	if this.GetEnergy() == this.GetMaxEnergy() {
		return 0
	} else {
		sec := this.GetEnergyTime().Add(time.Minute * time.Duration(RECOVERY_ENERGY_MINUTE_UNIT)).Sub(now).Seconds()

		return int32(sec)
	}

	return 0
}

// 尝试回体力
const (
	RECOVERY_ENERGY_MINUTE_UNIT = 5
)

// 每次访问主界面  计算一下 达没达到10分钟 ，需不需要回体力
func (this *UserAddAttr) RecoveryEnergy(now time.Time) {
	var diff time.Duration = now.Sub(this.GetEnergyTime())
	var intMinutes int32 = int32(diff.Minutes())
	addEnergy := (intMinutes / RECOVERY_ENERGY_MINUTE_UNIT) // 加多少体力

	if addEnergy > 0 {
		if (this.GetEnergy() + addEnergy) <= (this.GetMaxEnergy()) {
			newEnergyTime := this.GetEnergyTime().Add(time.Minute * time.Duration(intMinutes))
			this.SetEnergy(addEnergy + this.GetEnergy())
			this.SetEnergyTime(newEnergyTime)
		} else {
			if this.GetEnergy() <= this.GetMaxEnergy() {
				this.SetEnergy(this.GetMaxEnergy())
			}
			this.SetEnergyTime(now)
		}
	}
}
