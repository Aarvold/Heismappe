package queue

import (
	//def"config"
	//"fmt"
	"helpFunc"
	"sort"
	"sync"
	"driver"
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
	//fmt.Printf("%sOrder list is updated to %v \t current floor = %d \t cur dir = %d%s\n", def.ColR, Get_Orders(),driver.Get_cur_floor(),driver.Get_dir(), def.ColN)
}

func append_and_sort_list(orderlist []int, newOrder int) []int {
	newOrderlist := append(orderlist, newOrder)
	sort.Ints(newOrderlist)
	return newOrderlist
}

func Update_orderlist(orderlist []int, newOrder int) []int {

	if order_exists(orderlist,newOrder){
		return orderlist
	}

	tempOrderlist := append_and_sort_list(orderlist, newOrder)

	ordersDown := get_orders_down(tempOrderlist)
	ordersUp := get_orders_up(tempOrderlist)

	var newOrders []int

	if driver.Get_dir() >= 0{
		for _,orderUp := range ordersUp{
			if orderUp > driver.Get_cur_floor(){
				newOrders = append(newOrders,orderUp)
			}else {
				ordersDown = append(ordersDown,orderUp)
			}
		}
		newOrders = helpFunc.Append_list(newOrders,ordersDown)
	}

	if driver.Get_dir() == -1{
		for _,orderDown := range ordersDown{
			if -orderDown < driver.Get_cur_floor(){
				newOrders = append(newOrders,orderDown)
			}else {
				ordersUp = append(ordersUp,orderDown)
			}
		}
		newOrders = helpFunc.Append_list(newOrders,ordersUp)
	}

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


func order_exists(list []int,order int)bool{
	for j := 0; j < len(list); j++ {
		if list[j] == order {
			return true
		}
	}
	return false
}


func Remove_first_element_in_orders_and_save(){
	Set_Orders(append(Get_Orders()[:0], Get_Orders()[1:]...))
	Save_backup_to_file()
}