
package example

import(
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)

type TableLevelMaterial struct { 
	Comment	string   // optional 服务端本地化  
	TempID	string   // required 模板ID  
	Coin	int32   // optional 升级消耗金币  
	AddCoinPro	int32   // optional 金币加成（千分比）  
	UpStarData	[]struct   { 
			 TempID int32 	// optional 物品ID    
			 Num int32 	// optional 数量    
	}  
	UpLevelData	struct   { 
			 TempID int32 	// optional 物品ID    
			 Num []int32 	// repeated 数量    
	}  
}


var(
	iTableLevelMaterialList = map[string]*TableLevelMaterial{}
	iTableLevelMaterialMutex 	sync.RWMutex
	iTableLevelMaterialSize  uint32
	iTableLevelMaterialHook	func(list map[string]*TableLevelMaterial)
)

//从文件读取数据到内存
func TableLevelMaterial_ListUpdate(data []byte){

	list := []TableLevelMaterial{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}

	
	iTableLevelMaterialMutex.Lock()
	for k, item := range list {
		iTableLevelMaterialList[item.TempID] = &list[k]
	}
	iTableLevelMaterialMutex.Unlock()

	if iTableLevelMaterialHook != nil {
		iTableLevelMaterialMutex.RLock()
		iTableLevelMaterialHook(iTableLevelMaterialList)
		iTableLevelMaterialMutex.RUnlock()
	}

	atomic.StoreUint32(&iTableLevelMaterialSize, uint32(len(iTableLevelMaterialList)))
}

//唯一主键查找
func TableLevelMaterial_FindByPk(ID string) (tableLevelMaterial *TableLevelMaterial, err error){
	iTableLevelMaterialMutex.RLock()
	defer iTableLevelMaterialMutex.RUnlock()

	var ok bool
	tableLevelMaterial, ok = iTableLevelMaterialList[ID]
	if ok == false {
		err = errors.New("Not Data Found")
		return
	}


	return
}

//map的数据量大小
func TableLevelMaterial_ListLen() uint32 {
	return atomic.LoadUint32(&iTableLevelMaterialSize)
}

//获取完整数据
func TableLevelMaterial_ListAll() map[string]*TableLevelMaterial{
	iTableLevelMaterialMutex.RLock()
	defer iTableLevelMaterialMutex.RUnlock()

	m := map[string]*TableLevelMaterial{}

	for k, _ := range iTableLevelMaterialList {
		m[k] = iTableLevelMaterialList[k]
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func TableLevelMaterial_ListRange(f func(k string, v *TableLevelMaterial) bool) {
	iTableLevelMaterialMutex.RLock()
	defer iTableLevelMaterialMutex.RUnlock()


	for k, _ := range iTableLevelMaterialList {
		flag := f(k, iTableLevelMaterialList[k])
		if flag == false {
			return
		}
	}
}

//以下为兼容处理
func TableLevelMaterialList() map[string]*TableLevelMaterial{
	return TableLevelMaterial_ListAll()
}

func FindByPkTableLevelMaterial(ID string) (tableLevelMaterial *TableLevelMaterial, err error){
	return TableLevelMaterial_FindByPk(ID)
}

func TableLevelMaterialLen() uint32 {
	return TableLevelMaterial_ListLen()
}
