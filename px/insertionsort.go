package main

import "fmt"

func main() {

	arr := []int{7, 3, 5, 8, 6, 7, 4, 5}

	for i := 1; i < len(arr); i++ {
		insert(i, search(i, arr), arr)
	}

	fmt.Println(arr)
}

func insert(souceIndex int, dest int, arr []int) {

	v := arr[souceIndex]

	for i := souceIndex; i > dest; i-- {
		arr[i] = arr[i-1]
	}

	arr[dest] = v
}

func search(index int, arr []int) int {

	begin := 0
	end := index

	for begin < end {

		mid := (begin + end) >> 1

		if arr[mid] > arr[index] {
			end = mid
		} else {
			begin = mid + 1
		}
	}

	return begin
}
