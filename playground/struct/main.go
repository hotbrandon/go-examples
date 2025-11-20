package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{
		Name: "Brandon",
	}

	fmt.Println(p)

}
