package main

import "fmt"

/*
	两数之和
*/

func main() {

	nums := []int{2, 11, 15, 19, 7}
	target := 9

	fmt.Println(twoSum(nums, target))
}

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(nums); j++ {
			if i != j {
				if nums[i]+nums[j] == target {
					result := []int{i, j}
					return result
				}
			}
		}
	}
	return nil
}