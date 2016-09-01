package main

import (
	"fmt"
	"log"
	"os"
)

func produceStorge(table TableInfo) {
	fileName := fmt.Sprintf("%s/storage.go", fileDir)

	_, err := os.Stat(fileName)
	if err == nil || os.IsExist(err) {
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("[error]:create storage file error(%s).\n", err.Error())
		return
	}
	defer file.Close()

    file.WriteString(fmt.Sprintf("package %s\n\n", packageName))
    file.WriteString("import (\n\t\"database/sql\"\n\n\t\"github.com/dropbox/godropbox/memcache\"\n\t\"github.com/wgyuuu/storage\"\n)\n\n")

    resultObjName, addStorageStr := "StorageProxy", ""
    if table.ColumnList.KeyCount() > 1 {
        resultObjName = "ComplexStorage"
        addStorageStr = "Complex"
    }
    file.WriteString(fmt.Sprintf("func New%sStorage(db *sql.DB, mc memcache.Client, prefereExpireTime int) storage.%s {\n", table.TableName, resultObjName))
    file.WriteString(fmt.Sprintf("\tencoding := %sEncoding{}\n", table.TableName))
    file.WriteString(fmt.Sprintf("\tmsStorage := storage.New%sMysqlStorage(db, encoding)\n", addStorageStr))
    file.WriteString(fmt.Sprintf("\tmcStorage := storage.NewMemcStorage(mc, \"%s\", prefereExpireTime, encoding)\n", splitName(table.TableName)))
    file.WriteString(fmt.Sprintf("\treturn storage.New%sStorageProxy(mcStorage, msStorage)\n}\n", addStorageStr))
}
