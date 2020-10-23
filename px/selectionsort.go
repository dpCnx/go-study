package main

import "fmt"

func main() {

	arr := []int{7, 3, 5, 8, 6, 7, 4, 5}

	for i := len(arr) - 1; i > 0; i-- {
		index := 0
		for j := 1; j <= i; j++ {
			if arr[j] > arr[index] {
				index = j
			}
		}
		arr[i], arr[index] = arr[index], arr[i]
	}

	fmt.Println(arr)
}
