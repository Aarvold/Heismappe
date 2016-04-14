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

	go elevRun.Run_elev(outgoingMsg)
	go elevRun.Update_lights_orders(outgoingMsg)
	go network.Init(outgoingMsg, incomingMsg)
	go assign_external_order(costMsg)

	go func() {
		for {
			msg := <-incomingMsg
			handle_msg(msg, outgoingMsg, costMsg)
		}
	}()

	quit := make(chan int)
	go Quit_program(quit)

	<-quit
}

func handle_msg(msg def.Message, outgoingMsg, costMsg chan def.Message) {
	const aliveTimeout = 2 * time.Second

	switch msg.Category {
	case def.Alive:
		//if connection exists
		if connection, exist := onlineElevs[msg.Addr]; exist{
			connection.Timer.Reset(aliveTimeout)
		} else {
			addNewConnection(msg.Addr,aliveTimeout)
		}
	case def.NewOrder:
		fmt.Printf("%sNew order external recieved to floor %d %s \n", def.ColM, msg.Floor, def.ColN)
		driver.Set_button_lamp(msg.Button, msg.Floor, 1)
		def.Mutex.Lock()
		temp := def.Orders
		def.Mutex.Unlock()
		costMsg := def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: queue.Cost(temp, helpFunc.Order_dir(msg.Floor, msg.Button))}
		outgoingMsg <- costMsg

	case def.CompleteOrder:
		driver.Set_button_lamp(msg.Button, msg.Floor, 0)

	case def.Cost:
		//see assign_external_order
		costMsg <- msg
	}
}



func addNewConnection(addr string, aliveTimeout time.Duration){
	newConnection := network.UdpConnection{addr, time.NewTimer(aliveTimeout)}
	onlineElevs[addr] = newConnection
	numOfOnlineElevs = len(onlineElevs)
	go connectionTimer(&newConnection)
	fmt.Printf("%sConnection to IP %s established!%s\n", def.ColG, addr[0:15], def.ColN)
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

			//Sjekker om det det finnes en lik ordre i rcvList
			
			for oldOrder := range rcvList {
				if (newOrder.floor == oldOrder.floor) && (newOrder.button == oldOrder.button) {
					newOrder = oldOrder
				}
			}
			
			if costList, exist := rcvList[newOrder]; exist {
				if !costAlreadyRecieved(newCost,costList) {
					//fmt.Printf("Her blir det lagt tin en cost\n")
					rcvList[newOrder] = append(rcvList[newOrder], newCost)
				}
			} else {
				newOrder.timer = time.NewTimer(timeoutDuration)
				//fmt.Printf("Her blir en ordre med cost lagt til for forste gang\n")
				rcvList[newOrder] = []rcvCost{newCost}
				go costTimer(newOrder, overtime)
			}
			//fmt.Printf("%sLen rcvlst = %d and numOnlElev = %d %s\n", def.ColM, len(rcvList[newOrder]), numOfOnlineElevs, def.ColN)
			if  allCostsRecieved(rcvList,newOrder) {
				addNewOrder(rcvList,newOrder)
				delete(rcvList, newOrder)
				newOrder.timer.Stop()
				//fjerne ordren fra listen
			}
		case newOrder := <-overtime:
			fmt.Print("Assign order timeout: Did not recieve all replies before timeout\n")
			addNewOrder(rcvList,newOrder)
			delete(rcvList, newOrder)
		}
	}
}

func addNewOrder(rcvList map[rcvOrder][]rcvCost, newOrder rcvOrder){
	if this_elevator_has_the_lowest_cost(rcvList[newOrder]) {
		def.Mutex.Lock()
		def.Orders = queue.Update_orderlist(def.Orders, helpFunc.Order_dir(newOrder.floor, newOrder.button), false)
		fmt.Printf("%sExternal: Order list is updated to %v %s \n", def.ColR, def.Orders, def.ColN)
		def.Mutex.Unlock()
	}
}

func allCostsRecieved(rcvList map[rcvOrder][]rcvCost,newOrder rcvOrder)bool{
	return len(rcvList[newOrder]) == numOfOnlineElevs
}

func costAlreadyRecieved(newCost rcvCost,costList []rcvCost)bool{
	for _, adr := range costList {
		if newCost.elevAddr == adr.elevAddr {
			return  true
		}
	}
	return false
}

func costTimer(newOrder rcvOrder, overtime chan<- rcvOrder) {
	<-newOrder.timer.C
	overtime <- newOrder
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

func connectionTimer(connection *network.UdpConnection) {
	<-connection.Timer.C
	deadChan <- *connection
}

func Quit_program(quit chan int) {
	for {
		time.Sleep(time.Second)
		if driver.Get_stop_signal() == 1 {
			quit <- 1
			driver.Set_motor_dir(0)
		}
	}
}
