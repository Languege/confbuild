package main

import (
	"strings"
	"fmt"
	"text/template"
	"strconv"
	"os"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var funcMap = template.FuncMap{}

func init(){

	funcMap["IsString"] = func(dataType string) bool {
		if dataType == "string" {
			return true
		}

		return false
	}

	funcMap["IsBool"] = func(dataType string) bool {
		if dataType == "bool" {
			return true
		}

		return false
	}

	funcMap["StrFirstToLower"] = func(str string) string {

		var upperStr string
		vv := []rune(str)   // 后文有介绍
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				if vv[i] >= 65 && vv[i] <= 90 {  // 后文有介绍
					vv[i] += 32 // string的码表相差32位
					upperStr += string(vv[i])
				} else {
					fmt.Println("Not begins with upcase letter,")
					return str
				}
			} else {
				upperStr += string(vv[i])
			}
		}
		return upperStr
	}


	funcMap["NameTypeFunc"] = func(meta StructMeta) string {
		dataType := meta.DataType
		nameType := meta.NameType
		if nameType == "repeated" {
			if dataType == "int32" || dataType == "uint32" || dataType == "int64" || dataType == "uint64" || dataType == "string" {
				return "[]" + dataType
			}else if dataType == "bytes"{
				return "[]string"
			}else{
				return "[]struct"
			}
		}else if(nameType == "optional_struct"){
			return "struct"
		}else{
			if dataType == "int32" || dataType == "uint32" || dataType == "int64" || dataType == "uint64" || dataType == "string"{
				return dataType
			}else if dataType == "bytes"{
				return "[]byte"
			}else{
				return "struct"
			}
		}
	}
}

type StructMeta struct {
	Name 	string	//名称
	NameType	string //required,optional,repeated
	DataType string //string, int32, float32
	Comment	string //注释
	IsAnonymStruct	bool	//是否是匿名结构
	AnonymStruct *StructDesc	//匿名结构描述
}

type StructDesc struct {//结构体描述
	Name 	string
	PrimaryKey *StructMeta
	Field 	[]*StructMeta
	PropertyNum	int
	PkgName 	string //包名
}

type StructPkg struct {
	Name 	string 	//包名
	List 	[]*StructDesc //各个描述
}

//结构体描述解析 对应一张表
func StructDescParse(metas []*DataMeta, name string) (desc *StructDesc){
	desc = &StructDesc{
		Name:name,
		Field:[]*StructMeta{{"Comment", "optional", "string","服务端本地化", false, nil}},
		PkgName:pkg,
	}

	for i := 1; i < len(metas); i++ {
		fmt.Printf("i:%d m:%+v\n", i, metas[i])
		m := metas[i]
		sm := &StructMeta{
			Name:m.Name,
			NameType:m.NameType,
			DataType:m.DataType,
		}

		if m.Comment != "" {
			sm.Comment = strings.Replace(m.Comment, "\n", " ", -1)
		}

		switch m.NameType {
		case "required":
			//是否存在主键前缀
			if primaryPrefix != "" {
				sm.DataType = "string"
			}
			desc.PrimaryKey = sm
		case "optional_struct":
			sm.IsAnonymStruct = true
			//子结构
			subDesc, offset := Struct_OptionalStructParse(metas[i:])
			sm.AnonymStruct = subDesc
			sm.Name = subDesc.Name
			i += offset
		case "repeated":
			//是否是optional_struct
			if times, err := strconv.Atoi(m.DataType); err == nil {
				sm.IsAnonymStruct = true
				//子结构
				subDesc, offset := Struct_OptionalStructParse(metas[i+1:])
				//偏移量
				i += offset * times
				sm.AnonymStruct = subDesc
				sm.Name = subDesc.Name
			}
		case "optional":
		default:
			panic(fmt.Sprintf("undefined NameType %s", m.NameType))
		}

		desc.Field = append(desc.Field, sm)
	}

	return
}


func Struct_Print(desc *StructDesc) {
	if desc == nil {
		fmt.Println("desc is nil")
	}

	fmt.Println("============================")

	fmt.Println("Name", desc.Name)
	fmt.Println("PropertyNum", desc.PropertyNum)
	fmt.Printf("PrimaryKey:%+v\n", desc.PrimaryKey)
	fmt.Println("Field:%+v")
	for k, v := range desc.Field {
		fmt.Printf("		k:%d, field:%+v\n", k, v)
	}
	fmt.Println("============================")
}


//单表解析
func Struct_SheetParse(rows [][]string, sheet string) *StructDesc{
	fmt.Printf("Parsing Sheet %s ...\n", sheet)
	columnNum := len(rows[0])

	//元数据列表
	metaList :=make([]*DataMeta, columnNum)

	//前4行结构定义
	for j := 0;j < columnNum;j++ {
		cell := rows[0][j]
		if rows[0][j] == "" && j != 0 {
			break
		}

		switch cell {
		case "required","optional","repeated","optional_struct":
			meta := &DataMeta{
				NameType:rows[0][j],
				DataType:rows[1][j],
				Name:rows[2][j],
				Comment:rows[3][j],
			}

			metaList[j] = meta
		}
	}

	//结构体描述
	desc := StructDescParse(metaList, sheet)


	//Struct_Print(desc)
	//
	//structDescJson, _ := json.Marshal(desc)
	//fmt.Println(string(structDescJson))

	outFilename := fmt.Sprintf("%s/%s.go", outPath, sheet)

	outFile, err := os.Create(outFilename)
	if err != nil {
		panic(err.Error())
	}


	tmpl, err := template.New("struct_gen").Funcs(funcMap).Parse(sheet_tpl)
	if err != nil { panic(err) }
	err = tmpl.Execute(outFile, desc)
	if err != nil { panic(err) }

	outFile.Close()

	return desc
}

//optional_struct解析
func  Struct_OptionalStructParse(metas []*DataMeta) (desc *StructDesc, offset int){
	desc = &StructDesc{}
	propertyNum, err := strconv.Atoi(metas[0].DataType)
	if err != nil {
		panic(err.Error())
	}

	desc.PropertyNum = propertyNum
	desc.Name = metas[0].Name
	offset = propertyNum + 1

	for i := 1; i < offset; i++{
		m := metas[i]
		if m.NameType == "optional" || m.NameType == "repeated" {
			meta := &StructMeta{
				Name:m.Name,
				NameType:m.NameType,
				DataType:m.DataType,
				Comment:strings.Replace(m.Comment, "\n", " ", -1),
			}

			desc.Field = append(desc.Field, meta)
		}else{
			panic(fmt.Sprintf("optional_struct not support nesting i:%d m:%+v", i, m))
		}
	}


	return
}



func Struct_Parse(sheetSlice []string, xlsx *excelize.File) {
	sdlist := []*StructDesc{}
	for _, sheet := range sheetSlice {
		rows, err := xlsx.GetRows(sheet)
		if len(rows) <= 0 || err != nil {
			continue
		}

		desc := Struct_SheetParse(rows, sheet)
		sdlist = append(sdlist, desc)

	}

	strarr := strings.Split(excel, ".")
	outFilename := strings.Join(strarr[:len(strarr)-1], ".") + ".go"
	fmt.Println(outFilename)

	outFile, err := os.Create(outFilename)
	if err != nil {
		panic(err.Error())
	}



	tmpl, err := template.New("struct_gen").Funcs(funcMap).Parse(tpl)
	if err != nil { panic(err) }
	err = tmpl.Execute(outFile, &StructPkg{Name:pkg, List:sdlist})
	if err != nil { panic(err) }

	outFile.Close()
}