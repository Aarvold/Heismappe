package handleOrders

import (
	def"config"
	"driver"
	"queue"
	"time"
	"fmt"
	"helpFunc"
	"strconv"
)

type rcvCost struct {
	cost     int
	elevAddr string
}
type rcvOrder struct {
	floor  int
	button int
	timer  *time.Timer
}


func assign_external_order(costMsg, outgoingMsg, orderCompleted chan def.Message) {
	rcvList := make(map[rcvOrder][]rcvCost)
	var notCompletedOrders []rcvOrder
	var overtime = make(chan rcvOrder)
	var orderNotCompleted = make(chan rcvOrder)
	const notHandledTimeout = 9 * def.NumFloors * time.Second 
	const costTimeout = 1000 * time.Millisecond

	for {
		select {
		case msg := <-costMsg:

			order := rcvOrder{floor: msg.Floor, button: msg.Button}
			backupOrder := rcvOrder{floor: msg.Floor, button: msg.Button}
			newCost := rcvCost{cost: msg.Cost, elevAddr: msg.Addr[12:15]}

			// Checks if order allready is in rcvList  
			// i.e. if we already have recieved an order of this kind 
			for oldOrder := range rcvList {
				if equal_rcvOrders(oldOrder, order) {
					order = oldOrder
				}
			}
			
			// Maps the recieved cost to orders in rcvList
			if costList, exist := rcvList[order]; exist {
				if !cost_already_recieved(newCost,costList) {
					rcvList[order] = append(rcvList[order], newCost)
				}
			} else {
				//else: newCost is the first recieved cost for this order
				order.timer = time.NewTimer(costTimeout)
				rcvList[order] = []rcvCost{newCost}
				go cost_timer(order, overtime)
			}

			if  all_costs_recieved(rcvList,order) {
				add_order_to_best_elev(rcvList,order)
				delete(rcvList, order)
				order.timer.Stop()

				//This if makes sure that all external orders are handled
				if !order_exist_in_list(notCompletedOrders, backupOrder){
					backupOrder.timer = time.NewTimer(notHandledTimeout)
					notCompletedOrders = append(notCompletedOrders, backupOrder)
					go not_handled_timer(backupOrder, orderNotCompleted)
				}		
			}

		case order := <-overtime:
			fmt.Printf("%sAssign order timeout: Did not recieve all replies before timeout%s\n",def.ColY,def.ColN)
			add_order_to_best_elev(rcvList,order)
			delete(rcvList, order)

		case order := <-orderNotCompleted:
			order.timer.Stop()
			if ImConnected{
				fmt.Printf("%sOrder %v was not handled in time and is now resent%s\n", def.ColY, queue.Order_direction(order.floor, order.button), def.ColN)
				msgs := def.Message{Category: def.NewOrder, Floor: order.floor, Button: order.button, Cost: -1}
				outgoingMsg <- msgs				
			}

		case orderMsg := <- orderCompleted:
			for index,order := range notCompletedOrders{
				if orderMsg.Floor == order.floor && orderMsg.Button == order.button{
					order.timer.Stop()
					//Delete completed order
					notCompletedOrders = append(notCompletedOrders[:index], notCompletedOrders[(index+1):]...)
				}
			}
		}
	} 
}

func equal_rcvOrders(order1, order2 rcvOrder )bool{
	return (order1.floor == order2.floor) && (order1.button == order2.button)
}

func order_exist_in_list(rcvOrderList []rcvOrder, order rcvOrder) bool {
	for _,orderInList := range rcvOrderList{
		if equal_rcvOrders(orderInList, order){
			return true
		}
	}
	return false
}

func Cost(orderlist []int, newOrder int) int {
	new_orderlist := queue.Update_orderlist(orderlist, newOrder)
	index := helpFunc.Get_index(new_orderlist, newOrder)

	// Cost is initially the distance to first floor in list
	var cost = helpFunc.Difference_abs(driver.Get_cur_floor(), newOrder)

	// For every order in orderlist, the difference between two orders is
	// added to the cost
	if len(orderlist) > 0 {
		cost = helpFunc.Difference_abs(driver.Get_cur_floor(), orderlist[0])
		for i := 0; i < index-1; i++ {
			cost += helpFunc.Difference_abs(orderlist[i], orderlist[i+1])
		}
	}
	// The cost is penalized by the number of orders in the queue
	return int(cost) + 2*len(queue.Get_Orders())
}

func add_order_to_best_elev(rcvList map[rcvOrder][]rcvCost, newOrder rcvOrder){
	if this_elevator_has_the_lowest_cost(rcvList[newOrder]) {
		//fmt.Printf("%sNew order added to floor = %d with cost = %d\n%s",def.ColB,helpFunc.Order_dir(newOrder.floor,newOrder.button),Cost(queue.Get_Orders(),helpFunc.Order_dir(newOrder.floor,newOrder.button)),def.ColN)
		queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), queue.Order_direction(newOrder.floor, newOrder.button)))
		//fmt.Printf("%sUpdated orders = %v\n%s\n",def.ColB,queue.Get_Orders(),def.ColN)
		queue.Save_backup_to_file()		
	}
}

func all_costs_recieved(rcvList map[rcvOrder][]rcvCost,order rcvOrder)bool{
	return len(rcvList[order]) == NumOfOnlineElevs
}

func cost_already_recieved(newCost rcvCost,costList []rcvCost)bool{
	for _, adr := range costList {
		if newCost.elevAddr == adr.elevAddr {
			return  true
		}
	}
	return false
}

func cost_timer(order rcvOrder, overtime chan<- rcvOrder) {
	<-order.timer.C
	overtime <- order
}

func not_handled_timer(order rcvOrder, orderNotCompleted chan<- rcvOrder) {
	<-order.timer.C
	orderNotCompleted <- order
}

func this_elevator_has_the_lowest_cost(listOfCosts []rcvCost) bool {
	var bestCost = rcvCost{cost: 9999, elevAddr: "999"}

	for _, costStruct := range listOfCosts {
		if costStruct.cost < bestCost.cost {
			bestCost = costStruct
		} else if costStruct.cost == bestCost.cost {
			// If costs are equal, select lowest IP value to ensure that only 
			// one elevator handles the order
			costStructAddr, _ := strconv.Atoi(costStruct.elevAddr)
			bestCostAddr, 	_ := strconv.Atoi(bestCost.elevAddr)
			if costStructAddr > bestCostAddr {
				bestCost = costStruct
			}
		}
	}

	//if best cost is my cost, return true
	if bestCost.elevAddr == def.Laddr[12:15] {
		return true
	} else {
		return false
	}
}