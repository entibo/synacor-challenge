package main

import (
	"fmt"
	prmt "github.com/gitchander/permutation"
)

// 9 2 5 7 3
// blue red shiny concave corroded
func main() {
	nums := []int{2,3,5,7,9}
	p := prmt.New(prmt.IntSlice(nums))
	for p.Next() {
		r := 
			nums[0] + nums[1] * (nums[2]*nums[2]) + (nums[3]*nums[3]*nums[3]) - nums[4]
		if r == 399 {
			fmt.Println(nums)
			break
		}
	}
}