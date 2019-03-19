
package example

import (
	"github.com/spf13/viper"
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)


func UpdateConfAll() {

	TableLevelMaterial_ListUpdate()

	ChefBasic_ListUpdate()

}


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
	TableLevelMaterialData := viper.Get("TableLevelMaterial")
	TableLevelMaterialDataTmp, ok := TableLevelMaterialData.([]interface{})
	
	if ok == false {
		panic("TableLevelMaterial Conf Update failed,reason get json data failed")
		return
	}

	
	iTableLevelMaterialMutex.Lock()
	defer iTableLevelMaterialMutex.Unlock()

	for _, item := range TableLevelMaterialDataTmp {
		itemTmp, ok := item.(map[string]interface{})
		if ok == true {
			ele := &TableLevelMaterial{}
			bytesJson, _ := json.Marshal(itemTmp)
			json.Unmarshal(bytesJson, ele)

			TableLevelMaterialList[ele.TempID] = ele
		}
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


type ChefBasic struct { 
	Comment	string   // optional 服务端本地化  
	TempID	uint32   // required 模板ID  
	InitStar	int32   // optional 初始星级  
	InitStage	int32   // optional 初始段位  
	InitLevel	int32   // optional 初始等级  
	UnlockNeedCustomer	[]int32   // repeated 解锁需要的顾客id  
	Skills	[]int32   // repeated 技能  
	Score	int32   // optional 基础值  
	Speed	int32   // optional 速度  
	DishTypeScore1	int32   // optional 菜品种类1能力  
	DishTypeScore2	int32   // optional 菜品种类2能力  
	DishTypeScore3	int32   // optional 菜品种类3能力  
	DishTypeScore4	int32   // optional 菜品种类4能力  
	DishTypeScore5	int32   // optional 菜品种类5能力  
	DishTypeScore6	int32   // optional 菜品种类6能力  
	DishTypeScore7	int32   // optional 菜品种类7能力  
	DishTypeScore8	int32   // optional 菜品种类8能力  
	DishTypeScore9	int32   // optional 菜品种类9能力  
	DishTypeScore10	int32   // optional 菜品种类10能力  
	DishTypeScore11	int32   // optional 菜品种类11能力  
	DishTypeScore12	int32   // optional 菜品种类12能力  
	DishTypeScore13	int32   // optional 菜品种类13能力  
	DishTypeScore14	int32   // optional 菜品种类14能力  
	DishTypeScore15	int32   // optional 菜品种类15能力  
	DishTypeScore16	int32   // optional 菜品种类16能力  
	DishTypeScore17	int32   // optional 菜品种类17能力  
	DishTypeScore18	int32   // optional 菜品种类18能力  
	DishTypeScore19	int32   // optional 菜品种类19能力  
	DishTypeScore20	int32   // optional 菜品种类20能力  
}


var(
	ChefBasicList = map[uint32]*ChefBasic{}
	iChefBasicMutex 	sync.RWMutex
	iChefBasicSize  uint32
)

//从文件读取数据到内存
func ChefBasic_ListUpdate(){
	ChefBasicData := viper.Get("ChefBasic")
	ChefBasicDataTmp, ok := ChefBasicData.([]interface{})
	
	if ok == false {
		panic("ChefBasic Conf Update failed,reason get json data failed")
		return
	}

	
	iChefBasicMutex.Lock()
	defer iChefBasicMutex.Unlock()

	for _, item := range ChefBasicDataTmp {
		itemTmp, ok := item.(map[string]interface{})
		if ok == true {
			ele := &ChefBasic{}
			bytesJson, _ := json.Marshal(itemTmp)
			json.Unmarshal(bytesJson, ele)

			ChefBasicList[ele.TempID] = ele
		}
	}


	atomic.StoreUint32(&iChefBasicSize, uint32(len(ChefBasicList)))
}

//唯一主键查找
func ChefBasic_FindByPk(ID uint32) (chefBasic *ChefBasic, err error){
	iChefBasicMutex.RLock()
	defer iChefBasicMutex.RUnlock()

	var ok bool
	chefBasic, ok = ChefBasicList[ID]
	if ok == false {
		err = errors.New("Not Data Found")
		return
	}


	return
}

//map的数据量大小
func ChefBasic_ListLen() uint32 {
	return atomic.LoadUint32(&iChefBasicSize)
}

//获取完整数据
func ChefBasic_ListAll() map[uint32]*ChefBasic{
	iChefBasicMutex.RLock()
	defer iChefBasicMutex.RUnlock()

	m := map[uint32]*ChefBasic{}

	for k, v := range ChefBasicList {
		m[k] = v
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func ChefBasic_ListRange(f func(k uint32, v *ChefBasic) bool) {
	iChefBasicMutex.RLock()
	defer iChefBasicMutex.RUnlock()


	for k, v := range ChefBasicList {
		flag := f(k, v)
		if flag == false {
			return
		}
	}
}


