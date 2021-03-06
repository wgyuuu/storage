package storage

import (
	"database/sql"

	"github.com/0studio/databasetemplate"
	"github.com/wgyuuu/storage_key"
)

type MysqlEncoding interface {
	GetKey(obj interface{}) storage_key.Key
	Get(key storage_key.Key) string
	Add(obj interface{}) string
	Set(obj interface{}) string
    // return "" -> transfer Get
	MultiGet(keys []storage_key.Key) string
    // return "" -> transfer Set
	MultiSet(objMap map[storage_key.Key]interface{}) string
	Delete(key storage_key.Key) string
	ReadRow(resultSet *sql.Rows) (interface{}, error)
}

type MysqlStorage struct {
	databasetemplate.GenericDaoImpl
	encoding MysqlEncoding
}

func NewMysqlStorage(db *sql.DB, encoding MysqlEncoding) MysqlStorage {
	dbTemplate := databasetemplate.DatabaseTemplateImpl{Conn: db}
	return MysqlStorage{
		databasetemplate.GenericDaoImpl{DatabaseTemplate: &dbTemplate},
		encoding,
	}
}

func (this MysqlStorage) Get(key storage_key.Key) (interface{}, error) {
	return this.DatabaseTemplate.QueryObject(nil, this.encoding.Get(key), this.encoding.ReadRow)
}

func (this MysqlStorage) Add(key storage_key.Key, object interface{}) error {
	return this.DatabaseTemplate.Exec(nil, this.encoding.Add(object))
}

func (this MysqlStorage) Set(key storage_key.Key, object interface{}) error {
	return this.DatabaseTemplate.Exec(nil, this.encoding.Set(object))
}

func (this MysqlStorage) MultiGet(keys []storage_key.Key) (map[storage_key.Key]interface{}, error) {
	resultMap := make(map[storage_key.Key]interface{})
	execSql := this.encoding.MultiGet(keys)
	if len(execSql) > 0 {
		objList, err := this.DatabaseTemplate.QueryArray(nil, execSql, this.encoding.ReadRow)
		if err != nil {
			return nil, err
		}
		for _, obj := range objList {
			resultMap[this.encoding.GetKey(obj)] = obj
		}
	} else {
		for _, key := range keys {
			obj, err := this.Get(key)
			if err == nil {
				resultMap[key] = obj
			}
		}
	}
	return resultMap, nil
}

func (this MysqlStorage) MultiSet(objectMap map[storage_key.Key]interface{}) error {
	execSql := this.encoding.MultiSet(objectMap)
	if len(execSql) > 0 {
		return this.DatabaseTemplate.Exec(nil, execSql)
	} else {
		for key, obj := range objectMap {
			this.Set(key, obj)
		}
	}
	return nil
}

func (this MysqlStorage) Delete(key storage_key.Key) error {
	return this.DatabaseTemplate.Exec(nil, this.encoding.Delete(key))
}

func (this MysqlStorage) FlushAll() {
	return
}
