#need package
go get -u github.com/wgyuuu/storage_key
go get -u github.com/0studio/databasetemplate
go get -u github.com/0studio/redisapi


#test
cd stylei
generator -f ./entity.go

#nocite
暂时不支持time.Time作为key
storage文件生成就不会再替换，所以可以根据具体需求修改

#special1
有时候需要给结构体加锁，建议做法:
type UpgradeTes struct {
    Tes
    mutex sync.Mutex
}

var upTes UpgradeTes
storage := NewTesStorage(db, mc, 5)
storage.Set(storage_key.Uint64(1), upTes.Tes)

#special2
特殊sql语句是否可以单独保留一个mysql_drive变量用来执行特殊需求的sql语句?