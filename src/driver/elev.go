package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lpthread -lcomedi -lm
#include "elev.h"*/
import "C"

const (
	BUTTON_CALL_UP = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)
const N_floors int = int(C.N_FLOORS)

func Elev_init() {
	C.elev_init()
}

func Set_motor_dir(dirn int) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Set_button_lamp(button int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

//value er litt dårlig navn
func Set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func Get_button_signal(button int, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}
