#need package
go get -u github.com/wgyuuu/storage_key
go get -u github.com/0studio/databasetemplate
go get -u github.com/0studio/redisapi


#test
cd stylei
generator ./entity.go

#
暂时不支持time.Time作为key