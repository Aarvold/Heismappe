package main
//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
    "driver"
    "time"
    "asd"
    "fmt"
)



func main (){
	driver.Elev_init()
	var cur_floor int

	c := make(chan int)
	go func(){
		fmt.Println("jfsdkgæø")
		for{
			if driver.Get_floor_sensor_signal() != -1{
				cur_floor = driver.Get_floor_sensor_signal()
					c <- cur_floor
			}
			time.Sleep(10*time.Millisecond)
		}
	}()

	asd.Go_to_floor(2, c)
	time.Sleep(time.Second)


}