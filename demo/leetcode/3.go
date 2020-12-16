package main

import "fmt"

/*
	无重复字符的最长子串
*/

func main() {

	s := "uukewuukewuukew"

	fmt.Println(lengthOfLongestSubstring(s))

}

func lengthOfLongestSubstring(s string) int {

	val := []byte(s)

	kvMap := make([]int, 128)

	lens := len(s)

	var max, num int
	for i, j := 0, 0; i < lens && j < lens; j++ {
		if kvMap[val[j]] > i {
			i = kvMap[val[j]]
		}
		num = j - i + 1
		if num > max {
			max = num
		}
		kvMap[val[j]] = j + 1
	}
	return max
}
