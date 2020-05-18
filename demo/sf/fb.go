package sf

import "log"

/*
	0 1 2 3 4 5
 	0 1 1 2 3 5 8 13 ....
*/

func main() {

	log.Println(createFb(0))
}

func createFb(n int) int {

	if n <= 1 {
		return n
	}

	first := 0
	second := 1

	for i := 0; i < n-1; i++ {
		sum := first + second
		first = second
		second = sum
	}
	return second
}
