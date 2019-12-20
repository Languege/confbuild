package example

import (
	"fmt"
	"sync"
)

/**
 *@author LanguageY++2013
 *2019/12/20 4:24 PM
 **/
var(
	TableLevelHookCustomMap = map[string]*TableLevelMaterial{}
	iTableLevelHookCustomMapMutex   sync.RWMutex
)
func init() {
	iTableLevelMaterialHook = func(list map[uint32]*TableLevelMaterial) {
		iTableLevelHookCustomMapMutex.Lock()
		for _, v := range list {
			key := fmt.Sprintf("%d_%d", v.TempID, v.Coin)
			TableLevelHookCustomMap[key] = v
		}
		iTableLevelHookCustomMapMutex.Unlock()
	}
}