### 安装
```$xslt
go get github.com/Languege/confbuild
```

### 安装依赖包
```$xslt
bash install.sh
```

### 导出
```$xslt
go run main.go  struct_parser.go data_parser.go tpl.go \
-excel=./example/ConfData.xlsm \
-sheets="TableLevelMaterial,ChefBasic" \
-package=example
```

excel       excel配置表文件路径

sheets      需要导出的表，多张表以英文逗号","分隔

package     指定导出包名


推荐编译安装成系统命名后导出
```$xslt
go install Languege/confbuild
```

```$xslt
confbuild \
-excel=./example/ConfData.xlsm \
-sheets="TableLevelMaterial,ChefBasic" \
-package=example
```