package handleOrders

import (
	def"config"
	"queue"
	"driver"
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

var NumOfOnlineElevs int


func Handle_orders(outgoingMsg chan def.Message) {
	var alreadyPushed [def.NumFloors][def.NumButtons]bool

	for {
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if button_pushed(buttontype,floor) {

					if !alreadyPushed[floor][buttontype] {
						set_order_light(buttontype, floor)
						handle_new_order(buttontype,floor,outgoingMsg)
					}
					alreadyPushed[floor][buttontype] = true
				} else {
					alreadyPushed[floor][buttontype] = false
				}
			}
		}
	}
}

func set_order_light(buttontype, floor int){
	if !external_order(buttontype) {
		driver.Set_button_lamp(buttontype, floor, 1)
	}else if def.ImConnected{
		driver.Set_button_lamp(buttontype, floor, 1)
	}
}


func handle_new_order(buttontype,floor int,outgoingMsg chan def.Message){
	if external_order(buttontype) && def.ImConnected {
		msg := def.Message{Category: def.NewOrder, Floor: floor, Button: buttontype, Cost: -1}
		outgoingMsg <- msg
	} else if !external_order(buttontype){
		//Internal orders: if the desired floor is under the elevator it is set as a order down
		//fmt.Printf("Intern ordre mottat \n")
		if floor == def.CurFloor && def.CurDir == 1{
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor, false))
		}else if floor < def.CurFloor{
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor, false))
		}else if floor == def.NumFloors-1 {
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor, false))
		}else{
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), floor, false))
		}

		queue.Save_backup_to_file()
	}
}

func button_pushed(buttontype,floor int)bool{
	return driver.Get_button_signal(buttontype, floor) == 1
}

func external_order(buttontype int)bool{
	return buttontype == def.BtnUp || buttontype == def.BtnDown
}

//-----------------------------------------------------------------

func Assign_external_order(costMsg, outgoingMsg, orderIsCompleted chan def.Message) {
	rcvList := make(map[rcvOrder][]rcvCost)
	var notCompletedOrders []rcvOrder
	var overtime = make(chan rcvOrder)
	var orderNotHandled = make(chan rcvOrder)
	//var backupOrderComplete = make(chan rcvOrder)
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
				if (newOrder.floor == oldOrder.floor) && (newOrder.button == oldOrder.button) {
					newOrder = oldOrder
				}
			}
			
			// Adds the rcv cost to order in rcvList
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

				backupOrder.timer = time.NewTimer(notHandledTimeout)
				notCompletedOrders = append(notCompletedOrders, backupOrder)
				go handle_timer(backupOrder, orderNotHandled)
				//go handle_backupOrder_complete(backupOrder, backupOrderComplete, orderIsCompleted)
			}

		case newOrder := <-overtime:
			fmt.Print("Assign order timeout: Did not recieve all replies before timeout\n")
			add_new_order(rcvList,newOrder)
			delete(rcvList, newOrder)

		case order := <-orderNotHandled:
			fmt.Print("Order to floor %d was not handled in time, resending\n", order.floor )
			order.timer.Stop()
			if def.ImConnected{
				msgs := def.Message{Category: def.NewOrder, Floor: order.floor, Button: order.button, Cost: -1}
				outgoingMsg <- msgs				
			}
		case orderMsg := <- orderIsCompleted:
			fmt.Printf("BackupOrder sin timer er stoppet\n")
			for index,order := range notCompletedOrders{
				if orderMsg.Floor == order.floor && orderMsg.Button == order.button{
					order.timer.Stop()
					notCompletedOrders = append(notCompletedOrders[:index], notCompletedOrders[(index+1):]...)
				}
			}
		}
	} 
}

func add_new_order(rcvList map[rcvOrder][]rcvCost, newOrder rcvOrder){
	if this_elevator_has_the_lowest_cost(rcvList[newOrder]) {
		
		fmt.Printf("%sNew order added to floor = %d with cost = %d\n%s",def.ColY,helpFunc.Order_dir(newOrder.floor,newOrder.button),queue.Cost(queue.Get_Orders(),helpFunc.Order_dir(newOrder.floor,newOrder.button)),def.ColN)
		
		queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), helpFunc.Order_dir(newOrder.floor, newOrder.button), false))
		
		fmt.Printf("%sUpdated orders = %v\n%s\n",def.ColY,queue.Get_Orders(),def.ColN)
		
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

func handle_timer(order rcvOrder, orderNotHandled chan<- rcvOrder) {
	<-order.timer.C
	orderNotHandled <- order
}

func this_elevator_has_the_lowest_cost(listOfCosts []rcvCost) bool {
	//fmt.Print("Is this the best elev?\n")
	var bestCost = rcvCost{cost: 1000, elevAddr: "999"}

	for _, costStruct := range listOfCosts {
		//fmt.Printf("cost = %d addr %d\n", costStruct.cost, cS)
		if costStruct.cost < bestCost.cost {
			bestCost = costStruct
		} else if costStruct.cost == bestCost.cost {
			costStructAddr, _ := strconv.Atoi(costStruct.elevAddr)
			bestCostAddr, _ := strconv.Atoi(bestCost.elevAddr)
			// if equal cost: choose the minimum of the last three numbers in IP
			if costStructAddr > bestCostAddr {
				bestCost = costStruct
			}
		}
	}

	if bestCost.elevAddr == def.Laddr[12:15] {
		return true
	} else {
		return false
	}
}

