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
	"time"
	"handleOrders"

)

var onlineElevs = make(map[string]network.UdpConnection)


var outgoingMsg = make(chan def.Message, 10)
var incomingMsg = make(chan def.Message, 10)
var costMsg = make(chan def.Message, 10)
var orderIsCompleted = make(chan def.Message, 10)
var quitChan = make(chan int)


func main() {

	elevRun.Elev_init()
	
	go network.Init(outgoingMsg, incomingMsg)

	go elevRun.Run_elev(outgoingMsg)
	go handleOrders.Handle_orders(outgoingMsg, incomingMsg, costMsg, orderIsCompleted)
	go handle_msg(incomingMsg, outgoingMsg, costMsg, orderIsCompleted)

	go Quit_program(quitChan)
	<-quitChan
}

func handle_msg(incomingMsg, outgoingMsg, costMsg, orderIsCompleted chan def.Message) {
	for {
		msg := <-incomingMsg
		const aliveTimeout = 2 * time.Second
		switch msg.Category {
		case def.Alive:
			//if connection exists
			if connection, exist := onlineElevs[msg.Addr]; exist{
				connection.Timer.Reset(aliveTimeout)
			} else {
				add_new_connection(msg.Addr,aliveTimeout)
			}

		case def.NewOrder:
			//fmt.Printf("%sNew external order recieved to floor %d %s \n", def.ColM, helpFunc.Order_dir(msg.Floor, msg.Button), def.ColN)
			driver.Set_button_lamp(msg.Button, msg.Floor, 1)
			costMsg := def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: handleOrders.Cost(queue.Get_Orders(), helpFunc.Order_dir(msg.Floor, msg.Button))}
			outgoingMsg <- costMsg

		case def.CompleteOrder:
			driver.Set_button_lamp(msg.Button, msg.Floor, 0)
			orderIsCompleted <- msg

		case def.Cost:
			//see handleOrders.assign_external_order
			costMsg <- msg
		}
	}
}

func add_new_connection(addr string, aliveTimeout time.Duration){
	newConnection := network.UdpConnection{addr, time.NewTimer(aliveTimeout)}
	onlineElevs[addr] = newConnection
	handleOrders.NumOfOnlineElevs = len(onlineElevs)
	go connection_timer(&newConnection)
	fmt.Printf("%sConnection to IP %s established!%s\n", def.ColG, addr[0:15], def.ColN)
}

func connection_timer(connection *network.UdpConnection) {
	<-connection.Timer.C
}

func Quit_program(quit chan int) {
	for {
		time.Sleep(time.Second)
		if driver.Get_stop_signal() == 1 {
			driver.Set_motor_dir(0)
			quit <- 1
		}
	}
}
