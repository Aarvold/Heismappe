package asd

import(
	"driver"
	"math"
	"fmt"
)

func Go_to_floor(floor int, c chan int){
	var dir float64
	cur_floor:=<-c
	dir = (float64(floor-cur_floor))/math.Abs(float64(floor-cur_floor))
	driver.Set_motor_dir(int(dir))
	//event som venter på at den kommer til en etashe 
	for{
		fmt.Println("hey")
		cur_floor:=<-c
		fmt.Println(cur_floor)
		if cur_floor == floor{
			driver.Set_motor_dir(0)
			return
		}
	}
}