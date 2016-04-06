package asd

import (
	"driver"
	//"fmt"
	def "config"
	"math"
	"time"
)

func Go_to_floor(floor int) {
	var dir float64
	def.CurFloor = driver.Get_floor_sensor_signal()
	dir = (float64(floor - def.CurFloor)) / math.Abs(float64(floor-def.CurFloor))
	driver.Set_motor_dir(int(dir))
	//venter på at den kommer til en etashe
	for {
		def.CurFloor = driver.Get_floor_sensor_signal()
		if def.CurFloor == floor {
			driver.Set_motor_dir(0)
			driver.Set_door_open_lamp(1)
			time.Sleep(2 * time.Second)
			driver.Set_door_open_lamp(0)
			return
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func Update_lights() {
	for {
		def.CurFloor = driver.Get_floor_sensor_signal()
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if driver.Get_button_signal(buttontype, floor) == 1 {
					driver.Set_button_lamp(buttontype, floor, 1)
				}
				if def.CurFloor != -1 {
					driver.Set_floor_indicator(def.CurFloor)
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
