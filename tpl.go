package main

/**
 *@author LanguageY++2013
 *2019/3/10 11:07 AM
 **/
const tpl = `
package {{.Name}}

import(
	"errors"
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
`