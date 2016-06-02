package test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/dropbox/godropbox/memcache"
	"github.com/gogo/protobuf/proto"
	"github.com/wgyuuu/storage"
	"github.com/wgyuuu/storage_key"
)

func TestMysqlMemc(t *testing.T) {
	playerEncoding := PlayerEncoding{}
	mcStorage := storage.NewMemcStorage(testNewMClient(), "player", 86400, playerEncoding)
	mysqlStorage := storage.NewMysqlStorage(testDBClient(), playerEncoding)
	myStorage := storage.NewStorageProxy(mcStorage, mysqlStorage)

	err := myStorage.Add(storage_key.Uint64(1), Player{UserId: 1, Name: "test"})
	if err != nil {
		t.Error(err)
	}
	obj, err := myStorage.Get(storage_key.Uint64(1))
	if err != nil {
		t.Error(err)
	}
	t.Log(obj)
}

func printInfo(v ...interface{}) {
	log.Println(v)
}

func printError(err error) {
	log.Printf("[ERROR]:%s.\n", err.Error())
}

func testNewMClient() memcache.Client {
	memConfig := storage.MemcacheConfig{
		AddrList:             []string{"localhost:12000"},
		MaxActiveConnections: 5,
		MaxIdleConnections:   10,
		ReadTimeOutMS:        3000,
		WriteTimeOutMS:       500,
	}
	return storage.GetClient(memConfig, printError, printInfo)
}

func testDBClient() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/storage_test?charset=utf8mb4&parseTime=true&loc=Local&tls=false&timeout=1m")
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(5)
	return db
}

type PlayerEncoding struct {
}

func (this PlayerEncoding) Marshal(obj interface{}) ([]byte, error) {
	player := obj.(Player)
	return proto.Marshal(&player)
}

func (this PlayerEncoding) Unmarshal(bytes []byte) (interface{}, error) {
	var player Player
	err := proto.Unmarshal(bytes, &player)
	return player, err
}

func (this PlayerEncoding) GetKey(obj interface{}) storage_key.Key {
	player := obj.(Player)
	return storage_key.Uint64(player.UserId)
}

func (this PlayerEncoding) GetString(key storage_key.Key) string {
	return fmt.Sprintf(`select user_id, name from player where user_id = %s`, key.ToString())
}

// 可以继承Player重构函数
func (this PlayerEncoding) SetString(key storage_key.Key, obj interface{}) (sql string) {
	switch data := obj.(type) {
	case Player:
		sql = fmt.Sprintf(`update player set name = '%s' where user_id = %s`, data.Name, key.ToString())
	default:
		sql = ""
	}
	return
}

func (this PlayerEncoding) AddString(key storage_key.Key, obj interface{}) (sql string) {
	switch data := obj.(type) {
	case Player:
		sql = fmt.Sprintf(`insert into player (user_id, name) values(%d, '%s')`, data.UserId, data.Name)
	default:
		sql = ""
	}
	return
}

func (this PlayerEncoding) MultiGetString(keys []storage_key.Key) string {
	return ""
}

func (this PlayerEncoding) MultiSetString(objMap map[storage_key.Key]interface{}) string {
	return ""
}

func (this PlayerEncoding) DeleteString(key storage_key.Key) string {
	return ""
}

func (this PlayerEncoding) ReadRow(resultSet *sql.Rows) (interface{}, error) {
	obj := Player{}
	err := resultSet.Scan(
		&obj.UserId,
		&obj.Name,
	)
	return obj, err
}
