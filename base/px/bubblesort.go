package main

import "fmt"

func main() {

	arr := []int{7, 3, 5, 8, 6, 7, 4, 5}

	for i := len(arr) - 1; i > 0; i-- {
		for j := 1; j <= i; j++ {
			if arr[j] < arr[j-1] {
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
		}
	}

	fmt.Println(arr)
}
