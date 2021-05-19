package main

import "fmt"

func main() {
	nums1 := []int{1, 3}
	nums2 := []int{2}
	miDian := findMidianNumSortedArrays(nums1, nums2)
	fmt.Println(miDian)
}

func findMidianNumSortedArrays(nums1 []int, nums2 []int) int {
	len1 := len(nums1)
	len2 := len(nums2)
	t := len1 + len2
	if len1 == 0 { // 其中一个数组为空
		return findMidianNum(nums2)
	}
	if len2 == 0 {
		return findMidianNum(nums1)
	}
	i, j := 0, 0 // 两数组均为非空
	var nums []int
	for count := 0; count < (t); count++ {
		if i == len1 {
			if j < len2 { // 数组nums1被访问完，数组nums2没访问完
				s2 := nums2[j:len2]
				nums = append(nums, s2...)
				break
			}
		}
		if j == len2 {
			if i < len1 { // 数组nums2被访问完，数组nums1没访问完
				s1 := nums1[i:len1]
				nums = append(nums, s1...)
				break
			}
		}
		if nums1[i] < nums2[j] { // 归并排序的思想
			nums = append(nums, nums1[i])
			i++
		} else {
			nums = append(nums, nums2[j])
			j++
		}
	}
	midian := findMidianNum(nums)
	return midian
}

func findMidianNum(num []int) (mid int) {
	if len(num)%2 == 0 {
		mid = (num[len(num)/2] + num[len(num)/2-1]) / 2
	} else {
		mid = num[(len(num)-1)/2]
	}
	return mid
}
