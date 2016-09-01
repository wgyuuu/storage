package stylei

import (
	"database/sql"

	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage"
)

func NewTesStorage(db *sql.DB, mc memcache.Client, prefereExpireTime int) storage.ComplexStorage {
	encoding := TesEncoding{}
	msStorage := storage.NewComplexMysqlStorage(db, encoding)
	mcStorage := storage.NewMemcStorage(mc, "tes", prefereExpireTime, encoding)
	return storage.NewComplexStorageProxy(mcStorage, msStorage)
}
