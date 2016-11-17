package stylei

import (
	"database/sql"

	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage"
)

func NewUserStorage(db *sql.DB, mc memcache.Client, prefereExpireTime int) storage.ComplexStorageProxy {
	encoding := UserEncoding{}
	msStorage := storage.NewComplexMysqlStorage(db, encoding)
	mcStorage := storage.NewMemcStorage(mc, "user", prefereExpireTime, encoding)
	return storage.NewComplexStorageProxy(mcStorage, msStorage)
}
