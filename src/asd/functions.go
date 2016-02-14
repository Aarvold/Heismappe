package asd

import (
	"driver"
	//"fmt"
	"math"
	"time"
)

func Go_to_floor(floor int) {
	var dir float64
	var cur_floor int = driver.Get_floor_sensor_signal()
	dir = (float64(floor - cur_floor)) / math.Abs(float64(floor-cur_floor))
	driver.Set_motor_dir(int(dir))
	//event som venter på at den kommer til en etashe
	for {
		cur_floor = driver.Get_floor_sensor_signal()
		if cur_floor == floor {
			driver.Set_motor_dir(0)
			driver.Set_door_open_lamp(1)
			time.Sleep(1 * time.Second)
			driver.Set_door_open_lamp(0)
			return
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func Update_lights() {
	var cur_floor int
	for {
		cur_floor = driver.Get_floor_sensor_signal()
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < driver.N_floors; floor++ {
				if driver.Get_button_signal(buttontype, floor) == 1 {
					driver.Set_button_lamp(buttontype, floor, 1)
				}
				if cur_floor != -1 {
					driver.Set_floor_indicator(cur_floor)
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

func Quit_program(quit chan int) {
	for {
		time.Sleep(time.Second)
		if driver.Get_stop_signal() == 1 {
			quit <- 1
		}
	}
}

func Run_Elev() {

}
