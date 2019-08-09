package example

import (
	"github.com/Languege/redis_wrapper"
	"time"
	"os"
	"io/ioutil"
	"strings"
	"github.com/pkg/errors"
	"fmt"
)

var(
	confRedis  *redis_wrapper.RedisWrapper
)

type Configure struct {
	Path 		string
	RedisHost	string
	RedisPort 	string
	RedisPassword	string
	RedisMaxIdle	int
	RedisIdleTimeout	time.Duration
	RedisMaxActive 	int
}

func Start(conf Configure) {
	//实例化配置表redis
	confRedis = redis_wrapper.NewRedisWrapper(conf.RedisHost, conf.RedisPort, conf.RedisPassword, conf.RedisMaxIdle, conf.RedisIdleTimeout, conf.RedisMaxActive)
	//加载文件数据至redis
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

			key := strings.Split(fi.Name(), ".")[0]
			if key == "" {
				panic(errors.New(fmt.Sprintf("file name %s get failed", fi.Name())))
			}

			confRedis.SSet(key, content, 0, 0, false, false)
		}
	}

	UpdateConfAll()
}