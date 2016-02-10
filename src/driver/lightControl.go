package driver

func Light_clear_all(){
	ClearBit(LIGHT_UP1)
	ClearBit(LIGHT_UP2)
	ClearBit(LIGHT_UP3)
	ClearBit(LIGHT_DOWN2)
	ClearBit(LIGHT_DOWN3)
	ClearBit(LIGHT_DOWN4)
	ClearBit(LIGHT_COMMAND1)
	ClearBit(LIGHT_COMMAND2)
	ClearBit(LIGHT_COMMAND3)
	ClearBit(LIGHT_COMMAND4)
	ClearBit(LIGHT_STOP)
	ClearBit(LIGHT_FLOOR_IND1)
	ClearBit(LIGHT_FLOOR_IND2)
	// mangler IND3 og IND4 ?
	ClearBit(LIGHT_DOOR_OPEN)
}

func Set_light(){
	Set_bit(LIGHT_STOP)
}