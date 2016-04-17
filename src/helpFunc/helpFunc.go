package helpFunc

import(
	"math"
	"fmt"
)

func Difference_abs(val1 ,val2 int) int{
	return int(math.Abs(math.Abs(float64(val1))-math.Abs(float64(val2))))
}

func Get_index(list []int, element int) int {
	i := 0
	for i < len(list){
		if list[i] == element {
			return i
		}
		i++
	}
	fmt.Print("Error in Get_index\n")
	return -1
}

func Append_list(list1,list2 []int)[]int{
	i := 0
	for i < len(list2) {
		list1 = append(list1, list2[i])
		i++
	}
	return list1
}
