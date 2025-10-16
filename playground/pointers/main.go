package main

import (
	"fmt"
)

// Define the Person struct
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// MakePerson returns a Person value
func MakePerson(firstName, lastName string, age int) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// MakePersonPointer returns a pointer to a Person
func MakePersonPointer(firstName, lastName string, age int) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

func main() {
	// Call MakePerson (returns a value)
	personValue := MakePerson("John", "Doe", 30)
	fmt.Println("Person (value):", personValue)

	// Call MakePersonPointer (returns a pointer)
	personPointer := MakePersonPointer("Jane", "Smith", 25)
	fmt.Println("Person (pointer):", *personPointer) // Dereference to print value
}
