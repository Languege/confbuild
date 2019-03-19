package main

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"strings"
	"strconv"
	"os"
	"bytes"
	"encoding/json"
	"fmt"
)


type DataMeta struct {
	NameType	string	//required optional repeated optional_struct
	Name 	string	//属性名称
	DataType string //string, int32, float32
	Comment  	string
}


//数据解析单表
func Data_SheetParse(rows [][]string, sheet string)(data  []interface{}, err error){
	//提前创建好底层数组，避免复制
	data = make([]interface{}, 0, len(rows))
	if len(rows) == 0 {
		return
	}

	columnNum := len(rows[0])
	rowNum := len(rows)

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


	//4行及以后数据
	for i := 4;i < rowNum;i++ {
		//一行数据
		var item map[string]interface{}
		item, err = Data_RowParse(rows[i], metaList, columnNum)
		if err != nil {
			return
		}

		data = append(data, item)

	}


	return
}

//单个cell解析
func Data_CellParse(meta *DataMeta, value string)(cell interface{}) {
	var err error
	switch meta.DataType { //标量逗号分隔
	case "int32", "uint32","uint64","int64":
		if value == "" {
			cell = 0
			return
		}
		cell, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
	case "string":
		cell = value
	case "float32","float64","float":
		if value == "" {
			cell = 0.0
			return
		}
		cell, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
	case "bool":
		if value == "" {
			cell = false
			return
		}
		cell, err = strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("未识别的cell类型 %s", meta.DataType))
	}

	return
}

//解析optional_struct
func Data_StructParse(rows []string, metaSlice []*DataMeta)(data map[string]interface{})  {
	data = make(map[string]interface{})
	fmt.Printf("解析结构体 rows len:%d  metaSlice len:%d \n", len(rows), len(metaSlice))
	for k, v := range  rows {
		if meta := metaSlice[k]; meta != nil {
			if meta.NameType == "repeated" {
				var subItems []interface{}

				if v != "" {
					splits := strings.Split(v, ";")

					for   _, v := range splits{
						var item interface{}
						item = Data_CellParse(meta, v)

						subItems = append(subItems, item)
					}
				}

				data[meta.Name] = subItems
			}else{
				data[meta.Name] = Data_CellParse(meta, v)
			}
		}else{
			panic("meta undefined")
		}
	}

	return
}


//解析repeated数据
func Data_RepeatedParse(row []string, metaList []*DataMeta)(slice []interface{}, offset int, key string) {
	if repeatedTimes, err := strconv.ParseInt(metaList[0].DataType, 10 ,64); err == nil {//数字型，表名后面接结构体
		key = metaList[1].Name
		fmt.Printf("解析repeated 重复元素数量：%d  ", repeatedTimes)
		//获取结构体属性数量
		structPropertyNum, err := strconv.Atoi(metaList[1].DataType)
		if err != nil {
			panic(err)
		}

		fmt.Printf("元素属性数量：%d ", structPropertyNum)

		for t := 0; t < int(repeatedTimes); t ++ {
			start :=   t * (structPropertyNum+1) + 2
			end :=  (t + 1) * (structPropertyNum+1)
			fmt.Printf("t:%d   start:%d  end:%d row[start]:%s metaList[start]:%v ", t, start, end, row[start], metaList[start])
			item := Data_StructParse(row[start:end+1], metaList[start:end+1])

			slice = append(slice, item)
		}

		offset = (structPropertyNum+1) * int(repeatedTimes)

		fmt.Printf("列偏移量：%d \n", offset)
	}else{//标量，";"分隔
		key = metaList[0].Name

		if row[0] == "" {
			return
		}

		subItems := strings.Split(row[0], ";")
		fmt.Printf("解析repeated ;分隔 重复元素数量：%d\n", len(subItems))

		for _, si := range subItems {
			r := Data_CellParse(metaList[0], si)
			slice = append(slice, r)
		}
	}

	return
}

//一行数据解析
func Data_RowParse(row []string, metaList []*DataMeta, columnNum int)(item map[string]interface{}, err error)  {
	item = make(map[string]interface{})
	item["Comment"] = row[0]

	for c := 1; c < columnNum;c++ {
		fmt.Printf("c-->%d column %+v cell:%s\n ", c, metaList[c], row[c])

		meta :=  metaList[c]
		//不合法列，超出表头
		if meta == nil || (meta.Name == "" && meta.DataType == "") {
			continue
		}

		switch meta.NameType {
		case "repeated": //数组
			var slice []interface{}
			var offset int
			var key string

			slice, offset, key = Data_RepeatedParse(row[c:], metaList[c:])

			c += offset
			item[key] = slice
		case "optional_struct":
			fmt.Println("optional_struct_case")
			var structPropertyNum int
			structPropertyNum, err = strconv.Atoi(metaList[c].DataType)
			if err != nil {
				panic(err)
			}

			end := c + structPropertyNum + 1
			var embeddedStruct map[string]interface{}
			embeddedStruct = Data_StructParse(row[c+1:end], metaList[c+1:end])

			offset := structPropertyNum + 1
			c += offset
			item[meta.Name] = embeddedStruct
			fmt.Printf("embeddedStruct:%v 偏移量%d\n", embeddedStruct, offset)
		case "required", "optional"://required, optional
			var scalar interface{}
			scalar = Data_CellParse(meta, row[c])
			item[meta.Name] = scalar
		default:
			fmt.Println("空列")
		}
	}

	return
}


func Data_Parse(sheetSlice []string, xlsx *excelize.File) {
	mapData := make(map[string]interface{})
	for _, sheet := range sheetSlice {
		rows := xlsx.GetRows(sheet)
		if len(rows) <= 0 {
			panic("表不存在或者为空")
		}

		sheetSlice, err := Data_SheetParse(rows, sheet)
		if err != nil {
			panic(err)
		}

		mapData[sheet] = sheetSlice

	}

	jsonBytes, err := json.Marshal(mapData)
	if err != nil {
		panic(err.Error())
	}

	var out bytes.Buffer
	err = json.Indent(&out, jsonBytes, "", "\t")

	if err != nil {
		panic(err.Error())
	}


	strarr := strings.Split(excel, ".")
	outFilename := strings.Join(strarr[:len(strarr)-1], ".") + ".json"

	outFile, err := os.Create(outFilename)
	if err != nil {
		panic(err.Error())
	}

	_, err = outFile.Write(out.Bytes())
	if err != nil {
		panic(err.Error())
	}

	outFile.Close()
}