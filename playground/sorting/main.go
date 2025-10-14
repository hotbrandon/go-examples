package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func mySort(person []Person, less func(i, j int) bool) {
	for i := 0; i < len(person); i++ {
		for j := i + 1; j < len(person); j++ {
			if !less(i, j) {
				person[i], person[j] = person[j], person[i]
			}
		}

	}
}

func main() {

	person := []Person{
		{Name: "John", Age: 30},
		{Name: "Amy", Age: 25},
		{Name: "Bob", Age: 40},
		{Name: "Eve", Age: 35},
	}

	for _, p := range person {
		fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
	}

	sort.Slice(person, func(i, j int) bool {
		return person[i].Age < person[j].Age
	})

	fmt.Println("\nSorted by Age:")
	for _, p := range person {
		fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
	}

	mySort(person, func(i, j int) bool {
		return person[i].Name < person[j].Name
	})

	fmt.Println("\nSorted by name:")
	for _, p := range person {
		fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
	}
}
