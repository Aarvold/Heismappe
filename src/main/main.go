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

func main() {
	def.CurFloor = 4
	def.CurDir = -1
	var ordrs = []int{-3,-2,1,3,4}
	fmt.Printf("%d \n",queue.Cost(ordrs,1))

	fmt.Printf("%d \n",helpFunc.DifferenceAbs(1,-8))

	go elevRun.Run_elev()
	go elevRun.Update_lights_orders(outgoingMsg)
	go network.Init(outgoingMsg, incomingMsg)

	go func() {
		for {
			msg := <-incomingMsg
			handle_msg(msg,outgoingMsg)
		}
	}()

	
	//fmt.Printf("%v", queue.Update_orderlist(orderlist, 4))
	//fmt.Print(queue.Cost(orderlist, 4))

	quit := make(chan int)
	go Quit_program(quit)

	<-quit
}


func handle_msg(msg def.Message,outgoingMsg chan def.Message){
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