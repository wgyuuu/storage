package stylei

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/wgyuuu/storage/stylei/pb"
	"github.com/wgyuuu/storage_key"
)

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

func (this TesEncoding) GetString(key storage_key.Key) string {
	return fmt.Sprintf(`select user_id, level, gold from tes where user_id = %s`, key.ToString())
}

// 可以继承Tes重构函数
func (this TesEncoding) SetString(key storage_key.Key, obj interface{}) (string) {
    tes := obj.(Tes)
    return fmt.Sprintf(`update tes set level = '%d', gold = '%d' where user_id = %s`, tes.GetLevel(), tes.GetGold(), tes.GetUserId())
}

func (this TesEncoding) AddString(key storage_key.Key, obj interface{}) (sql string) {
	tes := obj.(Tes)
	return fmt.Sprintf(`insert into tes (user_id, level, gold) values(%d, %d, %d)`, tes.GetUserId(), tes.GetLevel(), tes.GetGold())
}

func (this TesEncoding) MultiGetString(keys []storage_key.Key) string {
	return ""
}

func (this TesEncoding) MultiSetString(objMap map[storage_key.Key]interface{}) string {
	return ""
}

func (this TesEncoding) DeleteString(key storage_key.Key) string {
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
