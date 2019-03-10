#!/usr/bin/env bash

go install



confbuild \
-excel=./example/ConfData.xlsm \
-sheets="TableLevelMaterial,ChefBasic" \
-package=example