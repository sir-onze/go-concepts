package main

import "fmt"

func main() {
	x := 10
	// defer func(x int) {
	// 	fmt.Println("deferred:", x) // x is read WHEN defer runs, not when registered
	// }(x)
	defer fmt.Println(x)
	x = 20
	fmt.Println(x)
}

// output:
// 20
// 20   ← sees the updated x
