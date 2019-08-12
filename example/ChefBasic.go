
package example

import(
	"sync"
	"errors"
	"encoding/json"
	"sync/atomic"
)

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
	iChefBasicList = map[uint32]*ChefBasic{}
	iChefBasicMutex 	sync.RWMutex
	iChefBasicSize  uint32
)

//从文件读取数据到内存
func ChefBasic_ListUpdate(){
	data, err := confRedis.HGet(GameConfDataKey, "ChefBasic")
	if err != nil {
		panic(err)
	}

	list := []ChefBasic{}

	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}

	
	iChefBasicMutex.Lock()
	defer iChefBasicMutex.Unlock()

	for _, item := range list {
		iChefBasicList[item.TempID] = &item
	}

	atomic.StoreUint32(&iChefBasicSize, uint32(len(iChefBasicList)))
}

//唯一主键查找
func ChefBasic_FindByPk(ID uint32) (chefBasic *ChefBasic, err error){
	iChefBasicMutex.RLock()
	defer iChefBasicMutex.RUnlock()

	var ok bool
	chefBasic, ok = iChefBasicList[ID]
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

	for k, _ := range iChefBasicList {
		m[k] = iChefBasicList[k]
	}

	return m
}

//自定义处理, 返回false, 终止遍历
func ChefBasic_ListRange(f func(k uint32, v *ChefBasic) bool) {
	iChefBasicMutex.RLock()
	defer iChefBasicMutex.RUnlock()


	for k, _ := range iChefBasicList {
		flag := f(k, iChefBasicList[k])
		if flag == false {
			return
		}
	}
}

//以下为兼容处理
func ChefBasicList() map[uint32]*ChefBasic{
	return ChefBasic_ListAll()
}

func FindByPkChefBasic(ID uint32) (chefBasic *ChefBasic, err error){
	return ChefBasic_FindByPk(ID)
}

func ChefBasicLen() uint32 {
	return ChefBasic_ListLen()
}
