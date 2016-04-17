package queue

import (
	def"config"
	//"fmt"
	"helpFunc"
	"sort"
	"sync"
	"driver"
)

// orders contains all orders and is sorted with next floor at orders[0]
// orders down will have a negative value and orders up will have a positive value
var orders []int
var mutex = &sync.Mutex{}

func Get_Orders()[]int{
	mutex.Lock()
	copyOrders := make([]int,len(orders))
	copy(copyOrders[:],orders)
	mutex.Unlock()
	return copyOrders
}

func Set_Orders(orderlist []int){
	mutex.Lock()
	copyOrders := make([]int,len(orderlist))
	copy(copyOrders[:], orderlist)
	orders = copyOrders
	mutex.Unlock()
	//fmt.Printf("%sOrder list is updated to %v \t current floor = %d \t cur dir = %d%s\n", def.ColR, Get_Orders(),driver.Get_cur_floor(),driver.Get_dir(), def.ColN)
}

func Update_orderlist(orderlist []int, newOrder int) []int {

	if order_exists(orderlist,newOrder){
		return orderlist
	}

	tempOrderlist := append_and_sort_list(orderlist, newOrder)
	// Split up the orderlist in orders up and down
	ordersDown := get_orders_down(tempOrderlist)
	ordersUp := get_orders_up(tempOrderlist)

	var updatedOrderlist []int

	// These two conditions updates the orders based on the elevators direction and current floor
	if driver.Get_dir() >= 0{
		for _,orderUp := range ordersUp{
			if orderUp > driver.Get_cur_floor(){
				updatedOrderlist = append(updatedOrderlist, orderUp)
			}else {
				ordersDown = append(ordersDown, orderUp)
			}
		}
		updatedOrderlist = helpFunc.Append_list(updatedOrderlist, ordersDown)
	}

	if driver.Get_dir() == -1{
		for _,orderDown := range ordersDown{
			if -orderDown < driver.Get_cur_floor(){
				updatedOrderlist = append(updatedOrderlist, orderDown)
			}else {
				ordersUp = append(ordersUp, orderDown)
			}
		}
		updatedOrderlist = helpFunc.Append_list(updatedOrderlist, ordersUp)
	}

	return updatedOrderlist
}

func append_and_sort_list(orderlist []int, newOrder int) []int {
	newOrderlist := append(orderlist, newOrder)
	sort.Ints(newOrderlist)
	return newOrderlist
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

func Remove_first_element_in_orders(){
	Set_Orders(append(Get_Orders()[:0], Get_Orders()[1:]...))
}

func Order_direction(floor,button int)int{
	if button == def.BtnDown{
		return -floor
	}else{
		return floor
	}
}