#!/usr/bin/env bash

go install



go run main.go  struct_parser.go data_parser.go tpl.go \
-excel=./example/ConfData.xlsm \
-sheets="TableLevelMaterial,ChefBasic" \
-package=example