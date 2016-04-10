package helpFunc

import(
	"math"
)


func DifferenceAbs(val1 ,val2 int) int{
	return int(math.Abs(math.Abs(float64(val1))-math.Abs(float64(val2))))
}