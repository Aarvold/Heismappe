package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	"elevRun"
	"driver"
	"fmt"
	"time"
	def "config"
	"network"
	"queue"
	"helpFunc"
)

var onlineElevs = make(map[string]network.UdpConnection)
var numOfOnlineElevs int

var outgoingMsg = make(chan def.Message, 10)
var incomingMsg = make(chan def.Message, 10)
var deadChan = make(chan network.UdpConnection)
var costMsg = make(chan def.Message,10)

func main() {
	def.CurFloor = 4
	def.CurDir = -1
	var ordrs = []int{-3,-2,1,3,4}

	fmt.Printf("%d \n",queue.Cost(ordrs,1))

	fmt.Printf("%d \n",helpFunc.DifferenceAbs(1,-8))

	go elevRun.Run_elev()
	go elevRun.Update_lights_orders(outgoingMsg)
	go network.Init(outgoingMsg, incomingMsg)

	//go fewafear(costMsg)

	go func() {
		for {
			msg := <-incomingMsg
			handle_msg(msg,outgoingMsg,costMsg)
		}
	}()

	
	//fmt.Printf("%v", queue.Update_orderlist(orderlist, 4))
	//fmt.Print(queue.Cost(orderlist, 4))

	quit := make(chan int)
	go Quit_program(quit)

	<-quit
}


func handle_msg(msg def.Message, outgoingMsg, costMsg chan def.Message){
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
		fmt.Printf("%s New order recieved %s \n",def.ColM,def.ColN)
		if msg.Button == def.BtnUp {
			driver.Set_button_lamp(msg.Button, msg.Floor, 1)
			costMsg := def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: queue.Cost(def.Orders, msg.Floor) }
			outgoingMsg <- costMsg
		}
		if msg.Button == def.BtnDown {
			driver.Set_button_lamp(msg.Button, msg.Floor, 1)
			costMsg := def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: queue.Cost(def.Orders, -msg.Floor) }
			outgoingMsg <- costMsg
		}

	case def.CompleteOrder:
		driver.Set_button_lamp(msg.Button, msg.Floor, 0)
		
	case def.Cost:
		costMsg<- msg
		/*
		costs = append(costs, msg)
		if len(onlineElevs)==numCostRecieved(order){
			sort(costs)
			if(costs[0]==cost(self)){
				if(costs[0]==costs[1]){
					//fiks det
				}else{
					def.Orders = queue.Update_orderlist

				}
			}
		}
		*/
		
	}
}


//----------------------------Dette skal legges et annet sted------------------------------

type reply struct {
	cost int
	lift string
}
type order struct {
	floor  int
	button int
	timer  *time.Timer
}

func assign_external_order(reply chan def.Message){
	recievedReplys := make(map[order][]reply)
	var overtime = make(chan *order)
	const timeoutDuration = 10 * time.Second

	for{
		select {
			case msg := <-reply:
				newOrder := order{floor: message.Floor, button: message.Button}
				newReply := reply{cost: message.Cost, lift: message.Addr[13:15]}




			case <- overtime:
				fmt.print("Assign order timeout: Did not recieve all replyes before timeout")
		}
	}
}

func costTimer(newOrder *order, timeout chan<- *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}

func equal(o1, o2 order) bool {
	return o1.floor == o2.floor && o1.button == o2.button
}

//-------------------------------------------------------------------------------------------


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