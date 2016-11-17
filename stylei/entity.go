package stylei

import "time"

type User struct {
	UserId uint64 `mysql:"primary_key=1,comment=玩家ID"`
	Level  int32  `mysql:"primary_key=3"`
	Name   string `mysql:"primary_key=2,varchar=128"`
	Gold   int32
	Actor  string `mysql:"varchar=512"`
	Time   time.Time
}
