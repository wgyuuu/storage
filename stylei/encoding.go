package stylei

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/wgyuuu/storage/stylei/pb"
	"github.com/wgyuuu/storage_key"
)

// encoding 方法
type TesEncoding struct {
}

func (this TesEncoding) Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return []byte{}, errors.New("obj is nil.")
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
	return storage_key.Uint64(tes.UserId)
}

func (this TesEncoding) Get(key storage_key.Key) string {
	return fmt.Sprintf(`select user_id, level, gold from tes where user_id = %s`, key.ToString())
}

// 可以继承Tes重构函数
func (this TesEncoding) Set(obj interface{}) string {
    tes := obj.(Tes)
    return fmt.Sprintf(`update tes set level = '%d', gold = '%d' where user_id = %s`, tes.GetLevel(), tes.GetGold(), tes.GetUserId())
}

func (this TesEncoding) Add(obj interface{}) (sql string) {
	tes := obj.(Tes)
	return fmt.Sprintf(`insert into tes (user_id, level, gold) values(%d, %d, %d)`, tes.GetUserId(), tes.GetLevel(), tes.GetGold())
}

func (this TesEncoding) MultiGet(keys []storage_key.Key) string {
	return ""
}

func (this TesEncoding) MultiSet(objMap map[storage_key.Key]interface{}) string {
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
		&tes.Gold,
	)
	return tes, err
}

func (this TestEncoding) ReadKeyRow(resultSet *sql.Rows) (interface{}, error) {
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
