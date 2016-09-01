package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	PACKAGE     string = "package "
	MYSQL       string = "mysql:"
	PRIMARY_KEY string = "primary_key"
	VARCHAR     string = "varchar"
)

type TableInfo struct {
	TableName  string
	ColumnList ColumnInfoList
}
type ColumnInfo struct {
	Name   string
	Typ    string
	SqlTyp string
	Attr   AttrInfo
}
type AttrInfo struct {
	PrimaryKey int
	VarcharLen int
}
type ColumnInfoList []ColumnInfo

func (this ColumnInfoList) MaxNameLen() (maxLen int) {
	for _, column := range this {
		if nameLen := len(column.Name); nameLen > maxLen {
			maxLen = nameLen
		}
	}
	return
}
func (this ColumnInfoList) MaxTypLen() (maxLen int) {
	for _, column := range this {
		if typLen := len(column.Typ); typLen > maxLen {
			maxLen = typLen
		}
	}
	return
}
func (this ColumnInfoList) KeyCount() (n int) {
	for _, column := range this {
		if column.Attr.PrimaryKey > 0 {
			n++
		}
	}
	return
}
func (this ColumnInfoList) CheckPrimary() {
	for _, column := range this {
		if column.Attr.PrimaryKey > 0 {
			return
		}
	}
	this[0].Attr.PrimaryKey = 1
}

var (
	fileFullPath string
	fileDir      string
	pathDir      string
	packageName  string
)

func init() {
	flag.StringVar(&fileFullPath, "f", "", "file for analysis")
	flag.Parse()
}

func main() {
	if len(fileFullPath) == 0 {
		log.Println("[error]:please spease the file.")
		return
	}

	file, err := os.Open(fileFullPath)
	if err != nil {
		log.Println("[error]:file lose.")
		return
	}
	defer file.Close()

	fileDir = filepath.Dir(fileFullPath)
	pathDir = getSpecifyDir(fileDir)

	br := bufio.NewReader(file)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			log.Printf("[warn]:file data error(%s).\n", err.Error())
			return
		}

		index := strings.Index(line, PACKAGE)
		if index == -1 {
			continue
		}
		packageName = line[index+len(PACKAGE) : len(line)-1]
		break
	}

	for {
		ok := analysisStruct(br)
		if !ok {
			break
		}
	}
}

func analysisStruct(br *bufio.Reader) bool {
	var table TableInfo

	// 获取表名
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("[warn]:analysis table name error(%s).\n", err.Error())
			}
			return false
		}
		if len(line) == 0 {
			continue
		}

		reg := regexp.MustCompile(`type [a-z,A-Z]* struct[ ]*{`)
		data := reg.FindString(line)
		if len(data) == 0 {
			continue
		}

		strList := strings.Split(data, " ")
		table.TableName = strList[1]
		break
	}

	// 获取字段
	loop := true
	for loop {
		line, err := br.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("[warn]:analysis column error(%s).\n", err.Error())
				return false
			}
			loop = false
		}
		if len(line) == 0 {
			continue
		}
		// 去掉tab 换行
		line = strings.TrimSpace(line)
		// 去掉注释
		index := strings.Index(line, "//")
		if index != -1 {
			line = line[:index]
		}

		// 结构体完结处
		if line == "}" {
			break
		}

		// 去掉多余空格
		reg := regexp.MustCompile("[ ]{1,}") // 至少一个
		line = reg.ReplaceAllString(line, " ")
		if line[0] == '@' {
			line = line[1:]
		}

		index = strings.Index(line, "`")
		var baseString, attrString string
		if index == -1 {
			baseString = line
			attrString = ""
		} else {
			baseString = line[:index-1]
			attrString = line[index+1 : len(line)-1]
		}
		name, typ := columnBase(baseString)
		attr := columnAttr(attrString)

		column := ColumnInfo{
			Name:   name,
			Typ:    typ,
			SqlTyp: getSqlType(typ),
			Attr:   attr,
		}
		table.ColumnList = append(table.ColumnList, column)
	}

	table.ColumnList.CheckPrimary()
	produceFile(table)
	return true
}

func produceFile(table TableInfo) {
	produceProto(table)
	produceEncoding(table)
	produceStorge(table)
}

func columnBase(data string) (name, typ string) {
	strList := strings.Split(data, " ")
	return strList[0], strList[1]
}

func getSqlType(typ string) (sqlType string) {
	switch typ {
	case "string":
		sqlType = "s"
	default:
		sqlType = "d"
	}
	return
}

func columnAttr(data string) (attr AttrInfo) {
	strList := strings.Split(data, " ")

	var attrString string
	for _, v := range strList {
		if strings.Index(v, MYSQL) == 0 {
			attrString = strings.ToLower(v[len(MYSQL):])
			break
		}
	}
	if len(attrString) == 0 {
		return
	}
	// 去掉“”
	attrString = attrString[1 : len(attrString)-1]

	primaryLen := len(PRIMARY_KEY)
	varcharLen := len(VARCHAR)
	for _, v := range strings.Split(attrString, ",") {
		switch {
		case v[:primaryLen] == PRIMARY_KEY:
			if len(v) == primaryLen {
				attr.PrimaryKey = 1
			} else {
				attr.PrimaryKey, _ = strconv.Atoi(v[primaryLen+1:])
			}
		case v[:varcharLen] == VARCHAR:
			if len(v) > varcharLen {
				attr.VarcharLen, _ = strconv.Atoi(v[varcharLen+1:])
			}
		}
	}
	return
}

func getSpecifyDir(dir string) (pathDir string) {
	currentDir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}

	dir2 := dir
	loop := true
	for loop {
		switch {
		case dir2[0] == '/':
			loop = false
			pathDir = dir2
		case dir2[:2] == "./":
			dir2 = dir2[2:]
		case dir2[:3] == "../":
			dir2 = dir2[3:]
			currentDir = currentDir[:strings.LastIndex(currentDir, "/")]
		default:
			loop = false
			pathDir = currentDir + "/" + dir2
		}
	}

	goPath := os.Getenv("GOPATH") + "/src"
	if strings.Contains(pathDir, goPath) {
		pathDir = pathDir[len(goPath)+1:]
	} else {
		pathDir = dir
	}
	return
}
