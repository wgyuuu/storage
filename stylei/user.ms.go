package stylei

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wgyuuu/storage/stylei/pb"
	"github.com/wgyuuu/storage_key"
)

// 结构体方法
func (this *User) GetUserId() uint64 {
	return this.UserId
}
func (this *User) SetUserId(userId uint64) {
	this.UserId = userId
}

func (this *User) GetLevel() int32 {
	return this.Level
}
func (this *User) SetLevel(level int32) {
	this.Level = level
}

func (this *User) GetName() string {
	return this.Name
}
func (this *User) SetName(name string) {
	this.Name = name
}

func (this *User) GetGold() int32 {
	return this.Gold
}
func (this *User) SetGold(gold int32) {
	this.Gold = gold
}

func (this *User) GetActor() string {
	return this.Actor
}
func (this *User) SetActor(actor string) {
	this.Actor = actor
}

func (this *User) GetTime() time.Time {
	return this.Time
}
func (this *User) SetTime(time time.Time) {
	this.Time = time
}

func (this *User) Serial() ([]byte, error) {
	user := pb.User{
		UserId: this.UserId,
		Level:  this.Level,
		Name:   this.Name,
		Gold:   this.Gold,
		Actor:  this.Actor,
		Time:   this.Time.Unix(),
	}
	return user.Marshal()
}
func (this *User) UnSerial(bytes []byte) error {
	var user pb.User
	if err := user.Unmarshal(bytes); err != nil {
		return err
	}
	this.UserId = user.UserId
	this.Level = user.Level
	this.Name = user.Name
	this.Gold = user.Gold
	this.Actor = user.Actor
	this.Time = time.Unix(user.Time, 0)
	return nil
}

// encoding 方法
type UserEncoding struct {}

func (this UserEncoding) Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return []byte{}, errors.New("obj nil")
	}
	user := obj.(User)
	return user.Serial()
}
func (this UserEncoding) Unmarshal(bytes []byte) (interface{}, error) {
	user := User{}
	err := user.UnSerial(bytes)
	return user, err
}

func (this UserEncoding) GetKey(obj interface{}) storage_key.Key {
	user := obj.(User)
	return storage_key.NewKeyList(storage_key.Uint64(user.UserId), storage_key.String(user.Name), storage_key.Int32(user.Level))
}

/*
create table if not exists user (
	user_id bigint(20) not null default 0 comment '玩家id',
	level int(11) not null default 0,
	name varchar(128) not null default '',
	gold int(11) not null default 0,
	actor varchar(512) not null default '',
	time timestamp not null default current_timestamp,
	primary key(user_id, name, level)
	)engine=InnoDB default charset=utf8;
*/
func (this UserEncoding) Get(key storage_key.Key) string {
	keyList := key.ToStringList()
	return fmt.Sprintf("select user_id, level, name, gold, actor, time from user where user_id=%s and name='%s' and level=%s", keyList[0], keyList[1], keyList[2])
}
func (this UserEncoding) Set(obj interface{}) string {
	user := obj.(User)
	return fmt.Sprintf("update user set gold=%d, actor='%s', time='%.19s' where user_id=%d and level=%d and name='%s'", user.GetGold(), user.GetActor(), user.GetTime(), user.GetUserId(), user.GetLevel(), user.GetName())
}
func (this UserEncoding) Add(obj interface{}) string {
	user := obj.(User)
	return fmt.Sprintf("insert into user (user_id, level, name, gold, actor, time) values(%d, %d, '%s', %d, '%s', '%.19s')", user.GetUserId(), user.GetLevel(), user.GetName(), user.GetGold(), user.GetActor(), user.GetTime())
}
func (this UserEncoding) MultiGet(keyList []storage_key.Key) string {
	return ""
}
func (this UserEncoding) MultiSet(mapObj map[storage_key.Key]interface{}) string {
	return ""
}
func (this UserEncoding) Delete(key storage_key.Key) string {
	return ""
}
func (this UserEncoding) ReadRow(resultSet *sql.Rows) (interface{}, error) {
	user := User{}
	err := resultSet.Scan(
		&user.UserId,
		&user.Level,
		&user.Name,
		&user.Gold,
		&user.Actor,
		&user.Time,
	)
	return user, err
}

// Complex
func (this UserEncoding) GetKeyList(key storage_key.Key) string {
	return fmt.Sprintf("select user_id, name, level from user where user_id=%s", key.ToString())
}
func (this UserEncoding) ReadKeyRow(resultSet *sql.Rows) (interface{}, error) {
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
