package main

import "fmt"

// Go is "call by value", but maps and slices are implemented with pointers.
//  Every type in Go is a value type. It’s just that sometimes the value is
// a pointer.
func modifyMap(m map[int]string) {
	m[2] = "six"
	m[4] = "four"
	delete(m, 1)
}

func main() {
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	modifyMap(m)

	fmt.Println(m)
}
