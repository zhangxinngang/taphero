package entity

type BStage struct {
	StageId    int32 // 关卡ID
	StageType  int32
	NeedEnergy int32
}
type BStageSet map[int32]BStage
