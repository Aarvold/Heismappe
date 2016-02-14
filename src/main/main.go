package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	"asd"
	"driver"
	//"fmt"
	"time"
)

func main() {
	driver.Elev_init()
	go asd.Update_lights()
	go asd.Go_to_floor(2)
	time.Sleep(time.Second)
	go asd.Go_to_floor(1)
	quit := make(chan int)
	go asd.Quit_program(quit)

	<-quit
}
