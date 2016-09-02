package storage

import (
	"database/sql"

	"github.com/wgyuuu/storage_key"
)

type ComplexMysqlEncoding interface {
	MysqlEncoding
	GetKeyList(key storage_key.Key) string
	ReadKeyRow(resultSet *sql.Rows) (interface{}, error)
}

type ComplexMysqlStorage struct {
	MysqlStorage
}

func NewComplexMysqlStorage(db *sql.DB, encoding ComplexMysqlEncoding) ComplexMysqlStorage {
	return ComplexMysqlStorage{
		MysqlStorage: NewMysqlStorage(db, encoding),
	}
}

func (this ComplexMysqlStorage) GetKeyList(key storage_key.Key) (keyList []storage_key.Key, err error) {
	complexEncoding := this.encoding.(ComplexMysqlEncoding)
	execSql := complexEncoding.GetKeyList(key)
	objList, err := this.DatabaseTemplate.QueryArray(nil, execSql, complexEncoding.ReadKeyRow)
	if err != nil {
		return
	}
	for _, obj := range objList {
		keyList = append(keyList, obj.(storage_key.Key))
	}
    return
}

func (this ComplexMysqlStorage) SetKeyList(key storage_key.Key, keyList []storage_key.Key) error {
    return nil
}
