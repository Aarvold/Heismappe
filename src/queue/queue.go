package queue

import (
	def "config"
	"fmt"
	"helpFunc"
	//"math"
	"sort"
)

func Append_and_sort_list(orderlist []int, newOrder int) []int {
	orderlist = append(orderlist, newOrder)
	sort.Ints(orderlist)
	return orderlist
}

func Update_orderlist(orderlist []int, newOrder int, costfunction bool) []int {

	//fmt.Printf("pre orders in append = %v \n", def.Orders)

	for j := 0; j < len(orderlist); j++ {
		if orderlist[j] == newOrder {
			fmt.Printf("Info from Update_orderlist: Orderto floor %d already ordered \n", newOrder)
			return orderlist
		}
	}

	tempOrderlist := Append_and_sort_list(orderlist, newOrder)

	var index = Get_element_index(tempOrderlist)
	temp1 := tempOrderlist[:index]
	temp2 := tempOrderlist[index:]

	/*if def.CurDir == -1 { //&& (def.CurFloor < int(math.Abs(float64(newOrder)))) {
		temp1 = tempOrderlist[index:]
		temp2 = tempOrderlist[:index]
		fmt.Printf("temp 1 = %v temp2 = %v tempOrderlist = %v \n", temp1, temp2, tempOrderlist)
	} else {
		temp1 = tempOrderlist[:index]
		temp2 = tempOrderlist[index:]
		fmt.Printf("temp 1 = %v temp2 = %v tempOrderlist = %v \n", temp1, temp2, tempOrderlist)
	}
	*/

	i := 0
	for i < len(temp2) {
		temp1 = append(temp1, temp2[i])
		i++
	}
	fmt.Printf("temp1 etter for loop = %v \n ", temp1)
	//orderlist = temp1

	//fmt.Printf("post orders in append = %v \n", def.Orders)

	//midletidig fix fordi den fucker opp
	if !costfunction {
		fmt.Printf("!!!!!!!!!!!!!!!!!! Orders !!!!!!!!!!!!!!!!!!\n")
		def.Orders = temp1
	}
	return temp1
}

func Get_element_index(orderlist []int) int {
	nextFloor := def.CurFloor
	if def.CurDir == -1 {
		//fmt.Printf("asfefaef\n")
		nextFloor = def.CurFloor - 1
	}

	//fmt.Printf("%sCurrent dir = %d %s \n", def.ColB, def.CurDir, def.ColN)
	orderNumber := def.CurDir * nextFloor
	var index = 0
	//må sjekkes
	for {
		if orderlist[index] > orderNumber {
			return index
		}
		if index == len(orderlist)-1 {
			return index
		}
		index++
	}
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
	fmt.Print("Error in Get_index")
	return -1
}

func Cost(orderlist []int, newOrder int) int {
	//fmt.Printf("pre orders in cost = %v \n", def.Orders)
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
	//fmt.Printf("post orders in cost = %v \n", def.Orders)
	return int(cost) + 2*len(def.Orders)
}
