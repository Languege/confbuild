package main

import (
	"github.com/Languege/confbuild/example"
	"log"
)

/**
 *@author LanguageY++2013
 *2020/4/21 5:26 PM
 **/
func main() {
	example.Start(example.Configure{
		Path:"./",
		EtcdEndPoints:[]string{"62.234.93.123:2379"},
		PrevKey: "open.confdata",
	})
	m, err := example.TableLevelMaterial_FindByPk(1)
	if err != nil {
		panic(err)
	}

	log.Println(m)
}