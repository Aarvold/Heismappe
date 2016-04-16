package handleOrders

import (
	def"config"
	"driver"
	"queue"
	"time"
	//"fmt"
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
	const notHandledTimeout = 7000 * def.NumFloors * time.Millisecond
	const costTimeout = 1000 * time.Millisecond

	for {
		select {
		case msg := <-costMsg:

			newOrder := rcvOrder{floor: msg.Floor, button: msg.Button}
			backupOrder := rcvOrder{floor: msg.Floor, button: msg.Button}
			newCost := rcvCost{cost: msg.Cost, elevAddr: msg.Addr[12:15]}

			// Checks if order allready is in rcvList
			for oldOrder := range rcvList {
				if equal_rcvOrders(oldOrder, newOrder) {
					newOrder = oldOrder
				}
			}
			
			// Maps the recieced cost to orders in rcvList
			if costList, exist := rcvList[newOrder]; exist {
				if !cost_already_recieved(newCost,costList) {
					rcvList[newOrder] = append(rcvList[newOrder], newCost)
				}
			} else {
				newOrder.timer = time.NewTimer(costTimeout)
				rcvList[newOrder] = []rcvCost{newCost}
				go cost_timer(newOrder, overtime)
			}

			if  all_costs_recieved(rcvList,newOrder) {
				add_new_order(rcvList,newOrder)
				delete(rcvList, newOrder)
				newOrder.timer.Stop()

				if order_exist_in_list(notCompletedOrders, backupOrder){
					backupOrder.timer = time.NewTimer(notHandledTimeout)
					notCompletedOrders = append(notCompletedOrders, backupOrder)
					go not_handled_timer(backupOrder, orderNotCompleted)
				}
				
			}

		case newOrder := <-overtime:
			//fmt.Print("Assign order timeout: Did not recieve all replies before timeout\n")
			add_new_order(rcvList,newOrder)
			delete(rcvList, newOrder)

		case order := <-orderNotCompleted:
			order.timer.Stop()
			if ImConnected{
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
	// The cost is penalties by the number of orders in the queue
	return int(cost) + 2*len(queue.Get_Orders())
}

func add_new_order(rcvList map[rcvOrder][]rcvCost, newOrder rcvOrder){
	if this_elevator_has_the_lowest_cost(rcvList[newOrder]) {
		//fmt.Printf("%sNew order added to floor = %d with cost = %d\n%s",def.ColY,helpFunc.Order_dir(newOrder.floor,newOrder.button),Cost(queue.Get_Orders(),helpFunc.Order_dir(newOrder.floor,newOrder.button)),def.ColN)
		queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), helpFunc.Order_dir(newOrder.floor, newOrder.button)))
		//fmt.Printf("%sUpdated orders = %v\n%s\n",def.ColY,queue.Get_Orders(),def.ColN)
		queue.Save_backup_to_file()		
	}
}

func all_costs_recieved(rcvList map[rcvOrder][]rcvCost,newOrder rcvOrder)bool{
	return len(rcvList[newOrder]) == NumOfOnlineElevs
}

func cost_already_recieved(newCost rcvCost,costList []rcvCost)bool{
	for _, adr := range costList {
		if newCost.elevAddr == adr.elevAddr {
			return  true
		}
	}
	return false
}

func cost_timer(newOrder rcvOrder, overtime chan<- rcvOrder) {
	<-newOrder.timer.C
	overtime <- newOrder
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