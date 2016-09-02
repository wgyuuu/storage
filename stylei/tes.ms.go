package stylei

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/wgyuuu/storage/stylei/pb"
	"github.com/wgyuuu/storage_key"
)

// 结构体方法
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

func (this *Tes) GetName() string {
	return this.Name
}
func (this *Tes) SetName(name string) {
	this.Name = name
}

func (this *Tes) GetGold() int32 {
	return this.Gold
}
func (this *Tes) SetGold(gold int32) {
	this.Gold = gold
}

func (this *Tes) GetActor() string {
	return this.Actor
}
func (this *Tes) SetActor(actor string) {
	this.Actor = actor
}

func (this *Tes) Serial() ([]byte, error) {
	tes := pb.Tes{
		UserId: this.UserId,
		Level:  this.Level,
		Name:   this.Name,
		Gold:   this.Gold,
		Actor:  this.Actor,
	}
	return tes.Marshal()
}
func (this *Tes) UnSerial(bytes []byte) error {
	var tes pb.Tes
	if err := tes.Unmarshal(bytes); err != nil {
		return err
	}
	this.UserId = tes.UserId
	this.Level = tes.Level
	this.Name = tes.Name
	this.Gold = tes.Gold
	this.Actor = tes.Actor
	return nil
}

// encoding 方法
type TesEncoding struct {}

func (this TesEncoding) Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return []byte{}, errors.New("obj nil")
	}
	tes := obj.(Tes)
	return tes.Serial()
}
func (this TesEncoding) Unmarshal(bytes []byte) (interface{}, error) {
	tes := Tes{}
	err := tes.UnSerial(bytes)
	return tes, err
}

func (this TesEncoding) GetKey(obj interface{}) storage_key.Key {
	tes := obj.(Tes)
	return storage_key.NewKeyList(storage_key.Uint64(tes.UserId), storage_key.String(tes.Name), storage_key.Int32(tes.Level))
}

func (this TesEncoding) Get(key storage_key.Key) string {
	keyList := key.ToStringList()
	return fmt.Sprintf("select user_id, level, name, gold, actor from tes where user_id=%s, name='%s', level=%s", keyList[0], keyList[1], keyList[2])
}
func (this TesEncoding) Set(obj interface{}) string {
	tes := obj.(Tes)
	return fmt.Sprintf("update tes set gold=%d, actor='%s' where user_id=%d and level=%d and name='%s'", tes.GetGold(), tes.GetActor(), tes.GetUserId(), tes.GetLevel(), tes.GetName())
}
func (this TesEncoding) Add(obj interface{}) string {
	tes := obj.(Tes)
	return fmt.Sprintf("insert into tes (user_id, level, name, gold, actor) values(%d, %d, '%s', %d, '%s')", tes.GetUserId(), tes.GetLevel(), tes.GetName(), tes.GetGold(), tes.GetActor())
}
func (this TesEncoding) MultiGet(keyList []storage_key.Key) string {
	return ""
}
func (this TesEncoding) MultiSet(mapObj map[storage_key.Key]interface{}) string {
	return ""
}
func (this TesEncoding) Delete(key storage_key.Key) string {
	return ""
}
func (this TesEncoding) ReadRow(resultSet *sql.Rows) (interface{}, error) {
	tes := Tes{}
	err := resultSet.Scan(
		&tes.UserId,
		&tes.Level,
		&tes.Name,
		&tes.Gold,
		&tes.Actor,
	)
	return tes, err
}

// Complex
func (this TesEncoding) GetKeyList(key storage_key.Key) string {
	return fmt.Sprintf("select user_id, name, level from tes where user_id=%s", key.ToString())
}
func (this TesEncoding) ReadKeyRow(resultSet *sql.Rows) (interface{}, error) {
	var userId uint64
	var name string
	var level int32
	err := resultSet.Scan(
		&userId,
		&name,
		&level,
	)
	return storage_key.NewKeyList(storage_key.Uint64(userId), storage_key.String(name), storage_key.Int32(level)), err
}
