package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func produceEncoding(table TableInfo) {
	fileName := fmt.Sprintf("%s/%s.ms.go", fileDir, strings.ToLower(table.TableName))
	os.Remove(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("[error]:create encoding file error(%s).\n", err.Error())
		return
	}
	defer file.Close()
	// 表的参数名
	tableParamete := parameter(table.TableName)

	file.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	file.WriteString(fmt.Sprintf("import (\n\t\"database/sql\"\n\t\"errors\"\n\t\"fmt\"\n\n\t\"%s/pb\"\n\t\"github.com/wgyuuu/storage_key\"\n)\n\n", pathDir))

	file.WriteString("// 结构体方法\n")
	for _, column := range table.ColumnList {
		formatGet := "func (this *%s) Get%s() %s {\n\treturn this.%s\n}\n"
		file.WriteString(fmt.Sprintf(formatGet, table.TableName, column.Name, column.Typ, column.Name))

		formatSet := "func (this *%s) Set%s(%s %s) {\n\tthis.%s = %s\n}\n\n"
		paramete := parameter(column.Name)
		file.WriteString(fmt.Sprintf(formatSet, table.TableName, column.Name, paramete, column.Typ, column.Name, paramete))
	}

	var serialString, unSerialString string = "", ""
	serialFormat := fmt.Sprintf("\t\t%%-%ds this.%%s,\n", table.ColumnList.MaxNameLen()+1)
	for _, column := range table.ColumnList {
		serialString += fmt.Sprintf(serialFormat, column.Name+":", column.Name)
		unSerialString += fmt.Sprintf("\tthis.%s = %s.%s\n", column.Name, tableParamete, column.Name)
	}

	file.WriteString(fmt.Sprintf("func (this *%s) Serial() ([]byte, error) {\n", table.TableName))
	file.WriteString(fmt.Sprintf("\t%s := pb.%s{\n", tableParamete, table.TableName))
	file.WriteString(serialString)
	file.WriteString(fmt.Sprintf("\t}\n\treturn %s.Marshal()\n}\n", tableParamete))

	file.WriteString(fmt.Sprintf("func (this *%s) UnSerial(bytes []byte) error {\n", table.TableName))
	file.WriteString(fmt.Sprintf("\tvar %s pb.%s\n", tableParamete, table.TableName))
	file.WriteString(fmt.Sprintf("\tif err := %s.Unmarshal(bytes); err != nil {\n\t\treturn err\n\t}\n", tableParamete))
	file.WriteString(unSerialString)
	file.WriteString("\treturn nil\n}\n\n")

	// encoding方法
	encodingName := fmt.Sprintf("%sEncoding", table.TableName)
	file.WriteString("// encoding 方法\n")
	file.WriteString(fmt.Sprintf("type %s struct {}\n\n", encodingName))
	// Marshal
	file.WriteString(fmt.Sprintf("func (this %s) Marshal(obj interface{}) ([]byte, error) {\n", encodingName))
	file.WriteString("\tif obj == nil {\n\t\treturn []byte{}, errors.New(\"obj nil\")\n\t}\n")
	file.WriteString(fmt.Sprintf("\t%s := obj.(%s)\n\treturn %s.Serial()\n}\n", tableParamete, table.TableName, tableParamete))
	// Unmarshal
	file.WriteString(fmt.Sprintf("func (this %s) Unmarshal(bytes []byte) (interface{}, error) {\n", encodingName))
	file.WriteString(fmt.Sprintf("\t%s := %s{}\n", tableParamete, table.TableName))
	file.WriteString(fmt.Sprintf("\terr := %s.UnSerial(bytes)\n", tableParamete))
	file.WriteString(fmt.Sprintf("\treturn %s, err\n}\n\n", tableParamete))
	// GetKey
	file.WriteString(fmt.Sprintf("func (this %s) GetKey(obj interface{}) storage_key.Key {\n", encodingName))
	file.WriteString(fmt.Sprintf("\t%s := obj.(%s)\n", tableParamete, table.TableName))
	keyList := make([]string, table.ColumnList.KeyCount())
	for _, column := range table.ColumnList {
		if column.Attr.PrimaryKey > 0 {
			keyList[column.Attr.PrimaryKey-1] = fmt.Sprintf("storage_key.%s(%s.%s)", strings.Title(column.Typ), tableParamete, column.Name)
		}
	}
	if len(keyList) == 1 {
		file.WriteString(fmt.Sprintf("\treturn %s\n", keyList[0]))
	} else {
		keyListString := "\treturn storage_key.NewKeyList("
		for k, v := range keyList {
			if k > 0 {
				keyListString += ", "
			}
			keyListString += v
		}
		keyListString += ")\n"
		file.WriteString(keyListString)
	}
	file.WriteString("}\n\n")
	// 处理过的表名
	tableName := splitName(table.TableName)

	// sql method
	// Get
	file.WriteString(fmt.Sprintf("func (this %s) Get(key storage_key.Key) string {\n", encodingName))
	file.WriteString("\tkeyList := key.ToStringList()\n")
	listKeyName := make([]string, table.ColumnList.KeyCount())
	getSqlString := "select "
	for k, column := range table.ColumnList {
		if k > 0 {
			getSqlString += ", "
		}
		columnName := splitName(column.Name)

		getSqlString += columnName
		if column.Attr.PrimaryKey > 0 {
			listKeyName[column.Attr.PrimaryKey-1] = columnName
		}
	}
	getSqlString += fmt.Sprintf(" from %s where ", tableName)
	for k, name := range listKeyName {
		if k > 0 {
			getSqlString += ", "
		}
		if table.ColumnList.IsString(name) {
			getSqlString += fmt.Sprintf("%s='%%s'", name)
		} else {
			getSqlString += fmt.Sprintf("%s=%%s", name)
		}

	}
	var getArgString string
	for k := range listKeyName {
		getArgString += fmt.Sprintf(", keyList[%d]", k)
	}
	file.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s\"%s)\n}\n", getSqlString, getArgString))

	// Set
	file.WriteString(fmt.Sprintf("func (this %s) Set(obj interface{}) string {\n", encodingName))
	file.WriteString(fmt.Sprintf("\t%s := obj.(%s)\n", tableParamete, table.TableName))
	var setSqlString, setWhereString, setValueString1, setValueString2 string
	for _, column := range table.ColumnList {
		if column.Attr.PrimaryKey == 0 {
			if len(setSqlString) > 0 {
				setSqlString += ", "
			}
			setSqlString += fmt.Sprintf("%s=%s", splitName(column.Name), column.SqlTyp)
			if len(setValueString1) > 0 {
				setValueString1 += ", "
			}
			setValueString1 += fmt.Sprintf("%s.Get%s()", tableParamete, column.Name)
		} else {
			if len(setWhereString) > 0 {
				setWhereString += " and "
			}
			setWhereString += fmt.Sprintf("%s=%s", splitName(column.Name), column.SqlTyp)
			if len(setValueString2) > 0 {
				setValueString2 += ", "
			}
			setValueString2 += fmt.Sprintf("%s.Get%s()", tableParamete, column.Name)
		}

	}
	setSqlString = fmt.Sprintf("update %s set %s where %s", tableName, setSqlString, setWhereString)
	file.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s\", %s, %s)\n}\n", setSqlString, setValueString1, setValueString2))

	// Add
	file.WriteString(fmt.Sprintf("func (this %s) Add(obj interface{}) string {\n", encodingName))
	file.WriteString(fmt.Sprintf("\t%s := obj.(%s)\n", tableParamete, table.TableName))
	addSqlString := fmt.Sprintf("insert into %s (", tableName)
	addTypeString := "values("
	addValueString := ""
	for k, column := range table.ColumnList {
		if k > 0 {
			addSqlString += ", "
			addTypeString += ", "
			addValueString += ", "
		}
		addSqlString += splitName(column.Name)
		addTypeString += fmt.Sprintf("%s", column.SqlTyp)
		addValueString += fmt.Sprintf("%s.Get%s()", tableParamete, column.Name)
	}
	file.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s) %s)\", %s)\n}\n", addSqlString, addTypeString, addValueString))

	// 剩下的
	file.WriteString(fmt.Sprintf("func (this %s) MultiGet(keyList []storage_key.Key) string {\n\treturn \"\"\n}\n", encodingName))
	file.WriteString(fmt.Sprintf("func (this %s) MultiSet(mapObj map[storage_key.Key]interface{}) string {\n\treturn \"\"\n}\n", encodingName))
	file.WriteString(fmt.Sprintf("func (this %s) Delete(key storage_key.Key) string {\n\treturn \"\"\n}\n", encodingName))

	// ReadRow
	file.WriteString(fmt.Sprintf("func (this %s) ReadRow(resultSet *sql.Rows) (interface{}, error) {\n", encodingName))
	file.WriteString(fmt.Sprintf("\t%s := %s{}\n", tableParamete, table.TableName))
	file.WriteString("\terr := resultSet.Scan(\n")
	for _, column := range table.ColumnList {
		file.WriteString(fmt.Sprintf("\t\t&%s.%s,\n", tableParamete, column.Name))
	}
	file.WriteString(fmt.Sprintf("\t)\n\treturn %s, err\n}\n\n", tableParamete))

	if len(listKeyName) <= 1 {
		return
	}

	// complex
	file.WriteString("// Complex\n")
	file.WriteString(fmt.Sprintf("func (this %s) GetKeyList(key storage_key.Key) string {\n", encodingName))
	getlistSqlString := "select "
	for k, name := range listKeyName {
		if k > 0 {
			getlistSqlString += ", "
		}
		getlistSqlString += name
	}
	if table.ColumnList.IsString(listKeyName[0]) {
		getlistSqlString += fmt.Sprintf("from %s where %s='%%s'", tableName, listKeyName[0])
	} else {
		getlistSqlString += fmt.Sprintf("from %s where %s=%%s", tableName, listKeyName[0])
	}
	file.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s\", key.ToString())\n}\n", getlistSqlString))

	file.WriteString(fmt.Sprintf("func (this %s) ReadKeyRow(resultSet *sql.Rows) (interface{}, error) {\n", encodingName))
	varKeyList := make([]string, table.ColumnList.KeyCount())
	valueKeyList := make([]string, table.ColumnList.KeyCount())
	for _, column := range table.ColumnList {
		if column.Attr.PrimaryKey > 0 {
			varKeyList[column.Attr.PrimaryKey-1] = fmt.Sprintf("\tvar %s %s\n", parameter(column.Name), column.Typ)
			valueKeyList[column.Attr.PrimaryKey-1] = fmt.Sprintf("\t\t&%s,\n", parameter(column.Name))
		}
	}
	for _, str := range varKeyList {
		file.WriteString(str)
	}
	file.WriteString("\terr := resultSet.Scan(\n")
	for _, str := range valueKeyList {
		file.WriteString(str)
	}
	file.WriteString("\t)\n")
	rowKeyList := make([]string, table.ColumnList.KeyCount())
	for _, column := range table.ColumnList {
		if column.Attr.PrimaryKey > 0 {
			rowKeyList[column.Attr.PrimaryKey-1] = fmt.Sprintf("storage_key.%s(%s)", strings.Title(column.Typ), parameter(column.Name))
		}
	}
	rowKeyListString := "\treturn storage_key.NewKeyList("
	for k, v := range rowKeyList {
		if k > 0 {
			rowKeyListString += ", "
		}
		rowKeyListString += v
	}
	rowKeyListString += "), err\n"
	file.WriteString(rowKeyListString)
	file.WriteString("}\n")
}

func splitName(name string) string {
	var indexList []int
	for k, v := range name {
		if v < 'a' {
			indexList = append(indexList, k)
		}
	}

	var sectionList []string
	start := 0
	name = strings.ToLower(name)
	for _, v := range indexList {
		if v == 0 {
			continue
		}
		sectionList = append(sectionList, name[start:v])
		start = v
	}
	sectionList = append(sectionList, name[start:len(name)])

	return strings.Join(sectionList, "_")
}

func parameter(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}
