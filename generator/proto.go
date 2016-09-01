package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func produceProto(table TableInfo) {
	os.MkdirAll(fileDir+"/pb", 0711)
	fileName := fmt.Sprintf("%s/pb/%s.proto", fileDir, strings.ToLower(table.TableName))
	os.Remove(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("[error]:create proto file error(%s).\n", err.Error())
		return
	}
	defer file.Close()
	defer execBuildProto(fileName)

	file.WriteString("// -*- coding:utf-8 -*-\n")
	file.WriteString("syntax = \"proto2\";\n\n")
	file.WriteString("package pb;\n\n")
	file.WriteString("import \"github.com/gogo/protobuf/gogoproto/gogo.proto\";\n\n")

	file.WriteString("option (gogoproto.sizer_all) = true;\n")
	file.WriteString("option (gogoproto.marshaler_all) = true;\n")
	file.WriteString("option (gogoproto.unmarshaler_all) = true;\n\n")

	file.WriteString(fmt.Sprintf("message %s {\n", table.TableName))
	format := fmt.Sprintf("\trequired %%-%ds %%-%ds = %%d [(gogoproto.nullable) = false];\n", table.ColumnList.MaxTypLen(), table.ColumnList.MaxNameLen())
	for k, column := range table.ColumnList {
		file.WriteString(fmt.Sprintf(format, column.Typ, strings.ToLower(column.Name[:1])+column.Name[1:], k+1))
	}
	file.WriteString("}\n")
}

func execBuildProto(fileName string) {
	arg1 := fmt.Sprintf("--proto_path=$GOPATH/src/github.com/gogo/protobuf/protobuf:../../../../:%s", fileDir)
	arg2 := fmt.Sprintf("--gogo_out=%s", fileDir)
	arg3 := fileName
	cmd := exec.Command("protoc", arg1, arg2, arg3)
	err := cmd.Run()
	if err != nil {
		log.Printf("[error]:run protoc error(%s).\n", err.Error())
	}
}
