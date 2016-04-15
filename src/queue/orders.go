package queue

import(
	//def"config"
	"sync"
	//"fmt"
)

var mutex = &sync.Mutex{}
var orders []int

func Get_Orders()[]int{
	mutex.Lock()
	copyOrders := make([]int,len(orders))
	copy(copyOrders[:],orders)
	mutex.Unlock()
	return copyOrders
}

func Set_Orders(newOrders []int){
	mutex.Lock()
	copyOrders := make([]int,len(newOrders))
	copy(copyOrders[:],newOrders)
	orders = copyOrders
	mutex.Unlock()
	//fmt.Printf("%sOrder list is updated to %v \t current floor = %d \t cur dir = %d%s\n", def.ColR, Get_Orders(),def.CurFloor,def.CurDir, def.ColN)
}