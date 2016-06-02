package stylei

import (
	"database/sql"

	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage"
)

func NewTesStorage(db *sql.DB, mClient memcache.Client, prefereExpireTime int) storage.StorageProxy {
    mysqlStorage := storage.NewMysqlStorage(db, TesEncoding{})
    memcStorage := storage.NewMemcStorage(mClient, "tes", prefereExpireTime, TesEncoding{})
    return storage.NewStorageProxy(memcStorage, mysqlStorage)
}
