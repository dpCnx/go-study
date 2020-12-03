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

/*
	https://blog.csdn.net/myz123321/article/details/90727065?utm_medium=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromBaidu-1.not_use_machine_learn_pai&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromBaidu-1.not_use_machine_learn_pai
*/
