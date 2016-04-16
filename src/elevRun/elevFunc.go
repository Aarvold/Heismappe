package elevRun

import (
	def "config"
	"driver"
	//"fmt"
	"helpFunc"
	"math"
	"queue"
	"time"

)


func Elev_init(){
	driver.Elev_init()
	driver.Set_motor_dir(-1)
	defer driver.Set_motor_dir(0)

	queue.Get_backup_from_file()

	for {
		if driver.Get_floor_sensor_signal() != -1 {
			return
		}
	}

}

func Run_elev(outgoingMsg chan def.Message) {
	for {

		floorSensor := driver.Get_floor_sensor_signal()
		update_current_floor(floorSensor)

		if at_top_or_bottom(floorSensor) {
			def.CurDir = -def.CurDir
		}
		if (len(queue.Get_Orders()) > 0) && (floorSensor != -1) {

			set_direction(queue.Get_Orders()[0],def.CurFloor)

			if arrived_at_destination() {
				//fmt.Printf("%sFloat cur floor = %v float orders[0] = %v %s \n", def.ColY, float64(floorSensor), math.Abs(float64(def.Orders[0])), def.ColN)
				driver.Set_motor_dir(0)

				if def.ImConnected{
					send_complete_msg(queue.Get_Orders()[0],outgoingMsg)
				}else{
					turn_off_lights()
				}

				remove_first_element_in_orders_and_save()

				//fmt.Printf("%sOrder to floor %d deleted, current floor = %d\n",def.ColB,queue.Get_Orders()[0],def.CurFloor)

				driver.Set_button_lamp(def.BtnInside, def.CurFloor, 0)

				open_close_door()
				

			}
		}
	}
}

func turn_off_lights(){
	if queue.Get_Orders()[0] < 0 {
		driver.Set_button_lamp(def.BtnDown, def.CurFloor, 0)
	} else {
		driver.Set_button_lamp(def.BtnUp, def.CurFloor, 0)
    }
}

func open_close_door(){
	driver.Set_door_open_lamp(1)
	time.Sleep(2 * time.Second)
	driver.Set_door_open_lamp(0)
}

func remove_first_element_in_orders_and_save(){
	queue.Set_Orders(append(queue.Get_Orders()[:0], queue.Get_Orders()[1:]...))
	queue.Save_backup_to_file()
}

func set_direction(order,curFloor int){
	//var dir float64
	if (math.Abs(float64(order)) - float64(curFloor)) != 0 {
		//difference divided with abs(difference) gives direction = -1 or 1
		dir := (math.Abs(float64(order)) - float64(curFloor)) / (float64(helpFunc.Difference_abs(order, curFloor)))
		driver.Set_motor_dir(int(dir))
		def.CurDir = int(dir)
		//fmt.Printf(def.ColY,"Current dir = %d \n",def.CurDir,def.ColN)
	}
}

func arrived_at_destination()bool{
	return def.CurFloor == int(math.Abs(float64(queue.Get_Orders()[0])))	
}

func at_top_or_bottom(floorSensor int)bool{
	return (floorSensor == 0) || (floorSensor == def.NumFloors-1)
}

func send_complete_msg(order int,outgoingMsg chan def.Message){
	if order < 0 {
		msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnDown, Cost: -1}
		outgoingMsg <- msg
	} else {
		msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnUp, Cost: -1}
		outgoingMsg <- msg
	}
}

func update_current_floor(floorSensorSignal int){
	if floorSensorSignal != -1 {
		def.CurFloor = floorSensorSignal
		driver.Set_floor_indicator(floorSensorSignal)
	}
}


/*--------------------------------------------------------------------------------------*/




