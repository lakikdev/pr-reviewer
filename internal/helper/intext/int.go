package intext

import "fmt"

func New(num int64) *int64 {
	return &num
}

func StringNum(num int) string {
	if num > 0 && num < 10 {
		return fmt.Sprintf("0%d", num)
	}
	return fmt.Sprintf("%d", num)
}

func IntArrayToStringArray(intArray []int64) []string {
	var stringArray []string
	for _, num := range intArray {
		stringArray = append(stringArray, fmt.Sprintf("%d", num))
	}
	return stringArray
}
