#need package
go get -u github.com/wgyuuu/storage_key
go get -u github.com/0studio/databasetemplate
go get -u github.com/0studio/redisapi


#test
cd stylei
generator ./entity.go

#
暂时不支持time.Time作为key
storage文件生成就不会再替换，所以可以根据具体需求修改