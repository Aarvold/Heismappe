package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	def"config"
	"elevRun"
	"network"
	"handleOrders"
	"communication"
	"time"
	"driver"

)

var outgoingMsg = make(chan def.Message, 10)
var incomingMsg = make(chan def.Message, 10)
var costMsg = make(chan def.Message, 10)
var orderCompleted = make(chan def.Message, 10)
var quitChan = make(chan int)


func main() {

	elevRun.Elev_init()
	
	go network.Init(outgoingMsg, incomingMsg)

	go elevRun.Run_elev(outgoingMsg)
	go handleOrders.Handle_orders(outgoingMsg, incomingMsg, costMsg, orderCompleted)
	go communication.Handle_msg(incomingMsg, outgoingMsg, costMsg, orderCompleted)

	go quit_program(quitChan)
	<-quitChan
}



func quit_program(quit chan int) {
	for {
		time.Sleep(time.Second)
		if driver.Get_stop_signal() == 1 {
			driver.Set_motor_dir(0)
			quit <- 1
		}
	}
}
