package main

/**
 *@author LanguageY++2013
 *2019/3/10 11:07 AM
 **/
const tpl = `
package {{.Name}}

import(
	"errors"
	"strings"
	"time"
	"context"
	"go.etcd.io/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"io/ioutil"
	"os"
	"fmt"
	"os/signal"
	"syscall"
)



var ErrTableNotExit = errors.New("config table not define")

func UpdateConf(table string, data []byte) error {
	switch table {
	{{range .List}}case "{{.Name}}":
		{{.Name}}_ListUpdate(data)
	{{end}}
	default:
		return ErrTableNotExit
	}

	return nil
}

type Configure struct {
	Path 		string
	EtcdEndPoints	[]string
	PrevKey 		string
}

func Start(conf Configure, rebuild bool) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.EtcdEndPoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}


	if rebuild {
		//加载文件数据至etcd
		dir, err := ioutil.ReadDir(conf.Path)
		if err != nil {
			panic(err)
		}
		PthSep := string(os.PathSeparator)
		for _, fi := range dir {
			if fi.IsDir() { // 忽略目录
				continue
			}
			if strings.HasSuffix(fi.Name(), "json") { //匹配文件
				content, err := ioutil.ReadFile(conf.Path + PthSep + fi.Name())
				if err != nil {
					panic(content)
				}

				table := strings.Split(fi.Name(), ".")[0]
				if table == "" {
					panic(errors.New(fmt.Sprintf("file name %s get failed", fi.Name())))
				}

				key := fmt.Sprintf("%s/%s", conf.PrevKey, table)

				_, err = cli.Put(context.Background(), key, string(content))
				if err != nil {
					panic(err)
				}
			}
		}
	}

	//get
	resp, err := cli.Get(context.Background(), conf.PrevKey, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	if resp.Kvs == nil {
		panic("配置项缺失")
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value;v != nil {
			table := strings.Replace(string(resp.Kvs[i].Key), conf.PrevKey, "", -1)
			UpdateConf(table, resp.Kvs[i].Value)
		}
	}


	//watcher
	go func() {
		rch := cli.Watch(context.Background(), conf.PrevKey, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					log.Printf("PUT Key %+v\n", string(ev.Kv.Key))
					table := strings.Replace(string(ev.Kv.Key), conf.PrevKey, "", -1)
					UpdateConf(table, ev.Kv.Value)
				case mvccpb.DELETE:
					log.Printf("DELETE Key %+v\n", string(ev.Kv.Key))
				}
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-c
		if cli != nil {
			cli.Close()
		}

		os.Exit(0)
	}()
}
`