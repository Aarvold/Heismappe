package helpFunc

import(
	"math"
	def "config"
)


func Difference_abs(val1 ,val2 int) int{
	return int(math.Abs(math.Abs(float64(val1))-math.Abs(float64(val2))))
}

func Order_dir(floor,button int)int{
	if button == def.BtnDown{
		return -floor
	}else{
		return floor
	}
}