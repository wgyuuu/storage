package stylei

type Tes struct {
	UserId uint64 `mysql:"primary_key=1"`
	Level  int32  `mysql:"primary_key=3"`
	Name   string `mysql:"primary_key=2,varchar=1024"`
	Gold   int32
	Actor  string `mysql:"varchar=512"`
}
