package main

/**
 *@author LanguageY++2013
 *2019/3/10 11:07 AM
 **/
const sheet_tpl = `
package {{.PkgName}}

import(
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)

type {{.Name}} struct { {{range .Field}}
	{{.Name}}	{{NameTypeFunc .}}  {{if .IsAnonymStruct}} { {{range .AnonymStruct.Field}}
			 {{.Name}} {{NameTypeFunc .}} 	// {{.NameType}} {{.Comment}}    {{end}}
	} {{else}} // {{.NameType}} {{.Comment}} {{end}} {{end}}
}


var(
	i{{.Name}}List = map[{{.PrimaryKey.DataType}}]*{{.Name}}{}
	i{{.Name}}Mutex 	sync.RWMutex
	i{{.Name}}Size  uint32
)

//从文件读取数据到内存
func {{.Name}}_ListUpdate(){
	data, err := confRedis.HGet(GameConfDataKey, "{{.Name}}")
	if err != nil {
		panic(err)
	}

	list := []{{.Name}}{}

	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}

	
	i{{.Name}}Mutex.Lock()
	defer i{{.Name}}Mutex.Unlock()

	for k, item := range list {
		i{{.Name}}List[item.TempID] = &list[k]
	}

	atomic.StoreUint32(&i{{.Name}}Size, uint32(len(i{{.Name}}List)))
}

//唯一主键查找
func {{.Name}}_FindByPk(ID {{.PrimaryKey.DataType}}) ({{StrFirstToLower .Name}} *{{.Name}}, err error){
	i{{.Name}}Mutex.RLock()
	defer i{{.Name}}Mutex.RUnlock()

	var ok bool
	{{StrFirstToLower .Name}}, ok = i{{.Name}}List[ID]
	if ok == false {
		err = errors.New("Not Data Found")
		return
	}


	return
}

//map的数据量大小
func {{.Name}}_ListLen() uint32 {
	return atomic.LoadUint32(&i{{.Name}}Size)
}

//获取完整数据
func {{.Name}}_ListAll() map[{{.PrimaryKey.DataType}}]*{{.Name}}{
	i{{.Name}}Mutex.RLock()
	defer i{{.Name}}Mutex.RUnlock()

	m := map[{{.PrimaryKey.DataType}}]*{{.Name}}{}

	for k, _ := range i{{.Name}}List {
		m[k] = i{{.Name}}List[k]
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func {{.Name}}_ListRange(f func(k {{.PrimaryKey.DataType}}, v *{{.Name}}) bool) {
	i{{.Name}}Mutex.RLock()
	defer i{{.Name}}Mutex.RUnlock()


	for k, _ := range i{{.Name}}List {
		flag := f(k, i{{.Name}}List[k])
		if flag == false {
			return
		}
	}
}

//以下为兼容处理
func {{.Name}}List() map[{{.PrimaryKey.DataType}}]*{{.Name}}{
	return {{.Name}}_ListAll()
}

func FindByPk{{.Name}}(ID {{.PrimaryKey.DataType}}) ({{StrFirstToLower .Name}} *{{.Name}}, err error){
	return {{.Name}}_FindByPk(ID)
}

func {{.Name}}Len() uint32 {
	return {{.Name}}_ListLen()
}
`