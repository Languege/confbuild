package main

import (
	"flag"
	"strings"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var sheets,excel, pkg string

func init(){
	flag.StringVar(&sheets, "sheets", "", "-sheets sheets to export, ',' split multiple sheets")
	flag.StringVar(&excel, "excel", "", "-excel  excel filename to parse")
	flag.StringVar(&pkg, "package", "", "-package  struct package name")
}
/**
 *@author LanguageY++2013
 *2019/3/10 1:01 AM
 **/
func main(){
	flag.Parse()

	if excel == "" {
		panic("excel can not empty")
	}

	if sheets == "" {
		panic("sheets can not empty")
	}

	if pkg == "" {
		panic("package can not empty")
	}

	sheetSlice := strings.Split(sheets, ",")

	xlsx, err := excelize.OpenFile(excel)
	if err != nil {
		panic(err.Error())
	}

	//数据解析
	Data_Parse(sheetSlice, xlsx)

	//结构解析
	Struct_Parse(sheetSlice, xlsx)
}