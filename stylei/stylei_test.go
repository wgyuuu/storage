package stylei

import (
	"database/sql"
	"log"
	"sync"
	"testing"

	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage"
	"github.com/wgyuuu/storage_key"
)

var (
	db *sql.DB
	mc memcache.Client
)

func logError(err error) {
	log.Println("[error]", err)
}
func logInfo(args ...interface{}) {
	info := "[info]"
	argList := []interface{}{info}
	for _, arg := range args {
		argList = append(argList, arg)
	}
	log.Println(argList...)
}

var once sync.Once

func init() {
	once.Do(create)
}

func create() {
	newMysql()
	newMemcache()
}

func newMemcache() {
	mcConf := storage.MemcacheConfig{
		AddrList:             []string{"127.0.0.1:12000"},
		MaxActiveConnections: 5,
		MaxIdleConnections:   5,
		ReadTimeOutMS:        3 * 1000,
		WriteTimeOutMS:       5 * 100,
	}
	mc = storage.GetClient(mcConf, logError, logInfo)
}

func newMysql() {
	db, _ = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local&tls=false&timeout=1m")
	db.SetMaxIdleConns(5)
}

func TestStorage(t *testing.T) {
	storage := NewTesStorage(db, mc, 5)
	tes := Tes{
		UserId: 11,
		Level:  1,
		Name:   "wang",
		Gold:   100011,
		Actor:  "aaa",
	}
	err := storage.Set(storage_key.Uint64(11), tes)
	t.Log(err)
    aa, err := storage.Get(storage_key.Uint64(11))
    t.Log(aa, err)
}
