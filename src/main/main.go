package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	"asd"
	"driver"
	//"fmt"
	//"time"
	def "config"
	"network"
)

func main() {
	driver.Elev_init()
	go asd.Update_lights()

	outgoingMsg := make(chan def.Message)
	incomingMsg := make(chan def.Message)
	go network.Init(outgoingMsg, incomingMsg)

	quit := make(chan int)
	go asd.Quit_program(quit)

	<-quit
}
