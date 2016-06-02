package stylei

import "github.com/wgyuuu/storage/stylei/pb"

type Tes struct {
	UserId uint64 `mysql:"primary_key"`
	Level  int32
	Gold   int32
}

func (this *Tes) GetUserId() uint64 {
	return this.UserId
}

func (this *Tes) SetUserId(userId uint64) {
	this.UserId = userId
}

func (this *Tes) GetLevel() int32 {
	return this.Level
}

func (this *Tes) SetLevel(level int32) {
	this.Level = level
}

func (this *Tes) GetGold() int64 {
	return this.Gold
}

func (this *Tes) SetGold(gold int64) {
	this.Gold = gold
}

func (this *Tes) Serial() (byte[], error) {
	var pbTes pb.Tes
	pbTes.UserId = this.UserId
	pbTes.Level = this.Level
	pbTes.Gold = this.Gold
	return pbTes.Marshal()
}

func (this *Tes) UnSerial(bytes byte[]) error {
	var pbTes pb.Tes
	err := pbTes.Unmarshal(bytes)
	if err != nil {
		return err
	}
	this.UserId = pbTes.UserId 
	this.Level = pbTes.Level  
	this.Gold = pbTes.Gold
	return nil
}
