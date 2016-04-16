package queue

import (
	def "config"
	"fmt"
	"helpFunc"
	//"math"
	"sort"
	"sync"
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

func append_and_sort_list(orderlist []int, newOrder int) []int {
	newOrderlist := append(orderlist, newOrder)
	copyOrderlist := make([]int,len(newOrderlist))
	copy(copyOrderlist[:],newOrderlist)
	sort.Ints(copyOrderlist)
	return copyOrderlist
}

func Update_orderlist(orderlist []int, newOrder int, costfunction bool) []int {
	//fmt.Printf("pre orders in append = %v \n", Get_Orders())
	copyOrderlist := orderlist

	if order_exists(copyOrderlist,newOrder){
		fmt.Printf("Info from Update_orderlist: Order to floor %d already ordered \n", newOrder)
		return copyOrderlist
	}

	tempOrderlist := append_and_sort_list(copyOrderlist, newOrder)

	ordersDown := get_orders_down(tempOrderlist)
	ordersUp := get_orders_up(tempOrderlist)

	var newOrders []int

	if def.CurDir >= 0{
		for _,orderUp := range ordersUp{
			if orderUp > def.CurFloor{
				newOrders = append(newOrders,orderUp)
			}else {
				ordersDown = append(ordersDown,orderUp)
			}
		}
		newOrders = append_list(newOrders,ordersDown)
	}

	if def.CurDir == -1{
		for _,orderDown := range ordersDown{
			if -orderDown < def.CurFloor{
				newOrders = append(newOrders,orderDown)
			}else {
				ordersUp = append(ordersUp,orderDown)
			}
		}
		newOrders = append_list(newOrders,ordersUp)
	}
	//fmt.Printf("newOrders = %v\n",newOrders)

	return newOrders
}

func get_orders_up(orderlist []int)[]int{
	var posOrders []int
	for _,order := range orderlist{
		if order >= 0 {
			posOrders = append(posOrders, order)
		}
	}
	return posOrders
}

func get_orders_down(orderlist []int)[]int{
	var negOrders []int
	for _,order := range orderlist{
		if order < 0 {
			negOrders = append(negOrders, order)
		}
	}
	return negOrders
}

func append_list(temp1,temp2 []int)[]int{
	i := 0
	for i < len(temp2) {
		temp1 = append(temp1, temp2[i])
		i++
	}
	return temp1
}

func order_exists(copyOrderlist []int,newOrder int)bool{
	for j := 0; j < len(copyOrderlist); j++ {
		if copyOrderlist[j] == newOrder {
			return true
		}
	}
	return false
}


func Get_index(orderlist []int, new_order int) int {
	orderlist_length := len(orderlist)

	i := 0
	for i < orderlist_length {
		if orderlist[i] == new_order {
			return i
		}
		i++
	}
	fmt.Print("Error in Get_index\n")
	return -1
}

func Cost(orderlist []int, newOrder int) int {
	new_orderlist := Update_orderlist(orderlist, newOrder, true)
	index := Get_index(new_orderlist, newOrder)
	//fmt.Print(index)

	var cost = helpFunc.Difference_abs(def.CurFloor, newOrder)

	if len(orderlist) > 0 {
		cost = helpFunc.Difference_abs(def.CurFloor, orderlist[0])
		for i := 0; i < index-1; i++ {
			cost += helpFunc.Difference_abs(orderlist[i], orderlist[i+1])
		}
	}
	return int(cost) + 2*len(Get_Orders())
}
