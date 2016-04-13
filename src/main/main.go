package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	def "config"
	"driver"
	"elevRun"
	"fmt"
	"helpFunc"
	"network"
	"queue"
	"strconv"
	"time"
)

var onlineElevs = make(map[string]network.UdpConnection)
var numOfOnlineElevs int

var outgoingMsg = make(chan def.Message, 10)
var incomingMsg = make(chan def.Message, 10)
var deadChan = make(chan network.UdpConnection)
var costMsg = make(chan def.Message, 10)

func main() {
	//def.CurFloor = 4
	//def.CurDir = -1
	var ordrs = []int{-3, -2, 1, 3, 4}

	fmt.Printf("%d \n", queue.Cost(ordrs, 1))
	fmt.Printf("%d \n", helpFunc.Difference_abs(1, -8))

	go elevRun.Run_elev()
	go elevRun.Update_lights_orders(outgoingMsg)
	go network.Init(outgoingMsg, incomingMsg)

	//go fewafear(costMsg)
	go assign_external_order(costMsg)
	go func() {
		for {
			msg := <-incomingMsg
			handle_msg(msg, outgoingMsg, costMsg)
		}
	}()

	//fmt.Printf("%v", queue.Update_orderlist(orderlist, 4))
	//fmt.Print(queue.Cost(orderlist, 4))

	quit := make(chan int)
	go Quit_program(quit)

	<-quit
}

func handle_msg(msg def.Message, outgoingMsg, costMsg chan def.Message) {
	const aliveTimeout = 2 * time.Second

	switch msg.Category {
	case def.Alive:
		if connection, exist := onlineElevs[msg.Addr]; exist {
			connection.Timer.Reset(aliveTimeout)
		} else {
			newConnection := network.UdpConnection{msg.Addr, time.NewTimer(aliveTimeout)}
			onlineElevs[msg.Addr] = newConnection
			numOfOnlineElevs = len(onlineElevs)
			go connectionTimer(&newConnection)
			fmt.Printf("%sConnection to IP %s established!%s\n", def.ColG, msg.Addr[0:15], def.ColN)
		}
	case def.NewOrder:
		fmt.Printf("%sNew order recieved %s \n", def.ColM, def.ColN)
		driver.Set_button_lamp(msg.Button, msg.Floor, 1)
		temp := def.Orders
		costMsg := def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: queue.Cost(temp, helpFunc.Order_dir(msg.Floor, msg.Button))}
		outgoingMsg <- costMsg

	case def.CompleteOrder:
		driver.Set_button_lamp(msg.Button, msg.Floor, 0)

	case def.Cost:
		//see assign_external_order
		costMsg <- msg
	}
}

//----------------------------Dette skal legges et annet sted------------------------------

type rcvCost struct {
	cost     int
	elevAddr string
}
type rcvOrder struct {
	floor  int
	button int
	timer  *time.Timer
}

func assign_external_order(costMsg chan def.Message) {
	rcvList := make(map[rcvOrder][]rcvCost)
	var overtime = make(chan rcvOrder)
	const timeoutDuration = 1000 * time.Millisecond

	for {
		select {
		case msg := <-costMsg:
			newOrder := rcvOrder{floor: msg.Floor, button: msg.Button}
			newCost := rcvCost{cost: msg.Cost, elevAddr: msg.Addr[12:15]}
			duplicate := false

			if costList, exist := rcvList[newOrder]; exist {

				for _, adr := range costList {
					if newCost.elevAddr == adr.elevAddr {
						duplicate = true
					}

				}
				if !duplicate {
					rcvList[newOrder] = append(rcvList[newOrder], newCost)
				}

			} else {
				newOrder.timer = time.NewTimer(timeoutDuration)
				rcvList[newOrder] = []rcvCost{newCost}
				go costTimer(newOrder, overtime)
			}
			if len(rcvList[newOrder]) == numOfOnlineElevs {
				if this_elevator_has_the_lowest_cost(rcvList[newOrder]) {
					temp := def.Orders
					temp = queue.Update_orderlist(temp, helpFunc.Order_dir(newOrder.floor, newOrder.button))
					fmt.Printf("%s Order list is updated to %v %s \n", def.ColR, def.Orders, def.ColN)
				}
				delete(rcvList, newOrder)
				newOrder.timer.Stop()
				//fjerne ordren fra listen
			}
		case newOrder := <-overtime:
			fmt.Print("Assign order timeout: Did not recieve all replies before timeout\n")
			if this_elevator_has_the_lowest_cost(rcvList[newOrder]) {
				def.Orders = queue.Update_orderlist(def.Orders, helpFunc.Order_dir(newOrder.floor, newOrder.button))
			}
			delete(rcvList, newOrder)

		}
	}
}

func costTimer(newOrder rcvOrder, overtime chan<- rcvOrder) {
	<-newOrder.timer.C
	overtime <- newOrder
}

func this_elevator_has_the_lowest_cost(listOfCosts []rcvCost) bool {
	//fmt.Print("Is this the best elev?\n")
	var bestCost = rcvCost{cost: 1000, elevAddr: "999"}

	for _, costStruct := range listOfCosts {
		cS, _ := strconv.Atoi(costStruct.elevAddr)
		fmt.Printf("cost = %d addr %d\n", costStruct.cost, cS)
		if costStruct.cost < bestCost.cost {
			bestCost = costStruct
		} else if costStruct.cost == bestCost.cost {
			cS, _ := strconv.Atoi(costStruct.elevAddr)
			bC, _ := strconv.Atoi(bestCost.elevAddr)
			// if equal cost: choose the minimum of the last three numbers in IP
			if cS <= bC {
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

func connectionTimer(connection *network.UdpConnection) {
	<-connection.Timer.C
	deadChan <- *connection
}

func Quit_program(quit chan int) {
	for {
		time.Sleep(time.Second)
		if driver.Get_stop_signal() == 1 {
			quit <- 1
		}
	}
}
