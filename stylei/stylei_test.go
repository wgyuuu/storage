package stylei

import (
	"database/sql"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/0studio/redisapi"
	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage"
	"github.com/wgyuuu/storage_key"
)

var (
	db *sql.DB
	mc memcache.Client
	rd redisapi.Redis

	encoding TesEncoding
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
	encoding = TesEncoding{}

	newMysql()
	newMemcache()
	newRedis()
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

func newRedis() {
	rd, _ = storage.NewRedisClient("localhost:6379")
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
	err := storage.Set(encoding.GetKey(tes), tes)
	t.Log(err)
	aa, err := storage.Get(encoding.GetKey(tes))
	t.Log(aa, err)
}

func TestMemcache(t *testing.T) {
	myStorage := NewTesStorage(db, mc, 5)

	tes := Tes{
		UserId: 11,
		Level:  1,
		Name:   "wang",
		Gold:   100011,
		Actor:  "aaa",
		Time:   time.Now(),
	}
	err := myStorage.Add(encoding.GetKey(tes), tes)
	t.Log(err)
	myStorage.PushKey(storage_key.Uint64(tes.GetUserId()), encoding.GetKey(tes))
	tes.Level = 2
	myStorage.Add(encoding.GetKey(tes), tes)
	myStorage.PushKey(storage_key.Uint64(tes.GetUserId()), encoding.GetKey(tes))
	keyList, err := myStorage.PreferedStorage.(storage.ComplexStorage).GetKeyList(storage_key.Uint64(tes.GetUserId()))
	t.Log("keylist1:", keyList, err)

	myStorage.FlushAll()
	tes.Name = "wwww"
	tes.Level = 3
	myStorage.Add(encoding.GetKey(tes), tes)
	myStorage.PushKey(storage_key.Uint64(tes.GetUserId()), encoding.GetKey(tes))
	keyList, err = myStorage.PreferedStorage.(storage.ComplexStorage).GetKeyList(storage_key.Uint64(tes.GetUserId()))
	t.Log("keylist2:", keyList, err)
}

func TestGet(t *testing.T) {
	myStorage := NewTesStorage(db, mc, 5)
	keyList, _ := myStorage.GetKeyList(storage_key.Uint64(11))
	for _, key := range keyList {
		obj, err := myStorage.Get(key)
		t.Log(key, "->", obj, err)
	}
}

func TestRedisMysql(t *testing.T) {
	msStorage := storage.NewComplexMysqlStorage(db, encoding)
	rdStorage := storage.NewComplexRedisStorage(rd, "tes", 100, encoding)
	myStorage := storage.NewComplexStorageProxy(rdStorage, msStorage)

	tes := Tes{
		UserId: 11,
		Level:  1,
		Name:   "wang",
		Gold:   100011,
		Actor:  "aaa",
		Time:   time.Now(),
	}
	err := myStorage.Add(encoding.GetKey(tes), tes)
	t.Log(err)
	myStorage.PushKey(storage_key.Uint64(tes.GetUserId()), encoding.GetKey(tes))
	tes.Level = 2
	myStorage.Add(encoding.GetKey(tes), tes)
	myStorage.PushKey(storage_key.Uint64(tes.GetUserId()), encoding.GetKey(tes))
	keyList, err := myStorage.PreferedStorage.(storage.ComplexStorage).GetKeyList(storage_key.Uint64(tes.GetUserId()))
	t.Log("keylist1:", keyList, err)

	tes.Name = "wwww"
	tes.Level = 3
	myStorage.Add(encoding.GetKey(tes), tes)
	myStorage.PushKey(storage_key.Uint64(tes.GetUserId()), encoding.GetKey(tes))
	keyList, err = myStorage.PreferedStorage.(storage.ComplexStorage).GetKeyList(storage_key.Uint64(tes.GetUserId()))
	t.Log("keylist2:", keyList, err)
}
