#!/usr/bin/env bash

#go build -o=${GOPATH}/bin/confbuild main.go  struct_parser.go data_parser.go tpl.go sheet_tpl.go



go run main.go  struct_parser.go data_parser.go tpl.go sheet_tpl.go \
-excel=./example/ConfData.xlsm \
-sheets="TableLevelMaterial,ChefBasic" \
-package=example \
-outpath=./example