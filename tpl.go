package main

/**
 *@author LanguageY++2013
 *2019/3/10 11:07 AM
 **/
const tpl = `
package {{.Name}}
import(
	"github.com/spf13/viper"
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)


func UpdateConfAll() {
{{range .List}}
	{{.Name}}_ListUpdate()
{{end}}
}

var ErrTableNotExit = errors.New("config table not define")

func UpdateConf(table string) error {
	switch table {
	{{range .List}}case "{{.Name}}":
		{{.Name}}_ListUpdate()
	{{end}}
	default:
		return ErrTableNotExit
	}

	return nil
}


{{range .List}}
type {{.Name}} struct { {{range .Field}}
	{{.Name}}	{{NameTypeFunc .}}  {{if .IsAnonymStruct}} { {{range .AnonymStruct.Field}}
			 {{.Name}} {{NameTypeFunc .}} 	// {{.NameType}} {{.Comment}}    {{end}}
	} {{else}} // {{.NameType}} {{.Comment}} {{end}} {{end}}
}


var(
	{{.Name}}List = map[{{.PrimaryKey.DataType}}]*{{.Name}}{}
	i{{.Name}}Mutex 	sync.RWMutex
	i{{.Name}}Size  uint32
)

//从文件读取数据到内存
func {{.Name}}_ListUpdate(){
	{{.Name}}Data := viper.Get("{{.Name}}")
	{{.Name}}DataTmp, ok := {{.Name}}Data.([]interface{})
	
	if ok == false {
		panic("{{.Name}} Conf Update failed,reason get json data failed")
		return
	}

	
	i{{.Name}}Mutex.Lock()
	defer i{{.Name}}Mutex.Unlock()

	for _, item := range {{.Name}}DataTmp {
		itemTmp, ok := item.(map[string]interface{})
		if ok == true {
			ele := &{{.Name}}{}
			bytesJson, _ := json.Marshal(itemTmp)
			json.Unmarshal(bytesJson, ele)

			{{.Name}}List[ele.{{.PrimaryKey.Name}}] = ele
		}
	}


	atomic.StoreUint32(&i{{.Name}}Size, uint32(len({{.Name}}List)))
}

//唯一主键查找
func {{.Name}}_FindByPk(ID {{.PrimaryKey.DataType}}) ({{StrFirstToLower .Name}} *{{.Name}}, err error){
	i{{.Name}}Mutex.RLock()
	defer i{{.Name}}Mutex.RUnlock()

	var ok bool
	{{StrFirstToLower .Name}}, ok = {{.Name}}List[ID]
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

	for k, v := range {{.Name}}List {
		m[k] = v
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func {{.Name}}_ListRange(f func(k {{.PrimaryKey.DataType}}, v *{{.Name}}) bool) {
	i{{.Name}}Mutex.RLock()
	defer i{{.Name}}Mutex.RUnlock()


	for k, v := range {{.Name}}List {
		flag := f(k, v)
		if flag == false {
			return
		}
	}
}

{{end}}
`