
package example

import(
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)

type TableLevelMaterial struct { 
	Comment	string   // optional 服务端本地化  
	TempID	uint32   // required 模板ID  
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
	TableLevelMaterialList = map[uint32]*TableLevelMaterial{}
	iTableLevelMaterialMutex 	sync.RWMutex
	iTableLevelMaterialSize  uint32
)

//从文件读取数据到内存
func TableLevelMaterial_ListUpdate(){
	data, err := confRedis.SGet("TableLevelMaterial")
	if err != nil {
		panic(err)
	}

	list := []TableLevelMaterial{}

	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}

	
	iTableLevelMaterialMutex.Lock()
	defer iTableLevelMaterialMutex.Unlock()

	for _, item := range list {
		TableLevelMaterialList[item.TempID] = &item
	}

	atomic.StoreUint32(&iTableLevelMaterialSize, uint32(len(TableLevelMaterialList)))
}

//唯一主键查找
func TableLevelMaterial_FindByPk(ID uint32) (tableLevelMaterial *TableLevelMaterial, err error){
	iTableLevelMaterialMutex.RLock()
	defer iTableLevelMaterialMutex.RUnlock()

	var ok bool
	tableLevelMaterial, ok = TableLevelMaterialList[ID]
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
func TableLevelMaterial_ListAll() map[uint32]*TableLevelMaterial{
	iTableLevelMaterialMutex.RLock()
	defer iTableLevelMaterialMutex.RUnlock()

	m := map[uint32]*TableLevelMaterial{}

	for k, v := range TableLevelMaterialList {
		m[k] = v
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func TableLevelMaterial_ListRange(f func(k uint32, v *TableLevelMaterial) bool) {
	iTableLevelMaterialMutex.RLock()
	defer iTableLevelMaterialMutex.RUnlock()


	for k, v := range TableLevelMaterialList {
		flag := f(k, v)
		if flag == false {
			return
		}
	}
}
