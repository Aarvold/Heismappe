package queue

import (
	def "config"
	"fmt"
	"math"
	"sort"
)

func Append_and_sort_list(orderlist []int, newOrder int) []int {
	orderlist = append(orderlist, newOrder)
	sort.Ints(orderlist)
	return orderlist
}

func Update_orderlist(orderlist []int, newOrder int) []int {
	tempOrderlist := Append_and_sort_list(orderlist, newOrder)

	var index = Get_element_index(tempOrderlist)
	temp1 := tempOrderlist[index:]
	temp2 := tempOrderlist[:index]

	temp2_length := len(temp2)

	i := 0
	for i < temp2_length {
		temp1 = append(temp1, temp2[i])
		i++
	}

	return temp1
}

func Get_element_index(orderlist []int) int {
	orderNumber := def.CurDir * def.CurFloor
	var index = 0
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

//må fikses
func Cost(orderlist []int, new_order int) int {
	new_orderlist := Update_orderlist(orderlist, new_order)
	index := Get_index(new_orderlist, new_order)
	fmt.Print(index)
	var cost float64 = math.Abs(float64(float64(def.CurFloor) - math.Abs(float64(orderlist[0]))))

	for i := 0; i < index-1; i++ {
		cost += math.Abs(float64(math.Abs(float64(orderlist[i])) - math.Abs(float64(orderlist[i+1]))))
	}
	return int(cost)
}
