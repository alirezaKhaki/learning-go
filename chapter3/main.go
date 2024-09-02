package main

import (
	"fmt"
)

func main() {
	// Call the functions to execute each exercise
	exercise1()
	exercise2()
	exercise3()
}

// Exercise 1: Define a variable named greetings of type slice of strings
// with the following values: "Hello", "Hola", "नमस्कार", "こんにちは", and "Привіт".
// Create a subslice containing the first two values;
// a second subslice with the second, third, and fourth values;
// and a third subslice with the fourth and fifth values.
// Print out all four slices.
func exercise1() {
	greetings := []string{"Hello", "Hola", "नमस्कार", "こんにちは", "Привіт"}

	// Create a subslice with the first two elements
	slice1 := greetings[:2]

	// Create a subslice with elements 1 through 3 (inclusive of 1, exclusive of 4)
	slice2 := greetings[1:4]

	// Create a subslice from the fourth element to the end
	slice3 := greetings[3:]

	// Print the original slice and the three subslices
	fmt.Println("Original slice:", greetings)
	fmt.Println("Subslice 1:", slice1)
	fmt.Println("Subslice 2:", slice2)
	fmt.Println("Subslice 3:", slice3)

	// Explanation:
	// We defined the 'greetings' slice with five international greetings.
	// Three subslices were created:
	// - 'slice1' contains the first two elements.
	// - 'slice2' contains the second to fourth elements.
	// - 'slice3' contains the fourth and fifth elements.
	// All slices are printed to verify their content.
}

// Exercise 2: Define a string variable called message with the value "Hi 😘 and 😊 "
// and print the fourth rune in it as a character, not a number.
func exercise2() {
	message := "Hi 😘 and 😊 "
	// Print the fourth rune (index 3) as a character using %c format specifier
	fmt.Printf("Fourth rune: %c\n", message[3])

	// Explanation:
	// We defined a string 'message' with the value "Hi 😘 and 😊 ".
	// We accessed the fourth rune (index 3) of the string and printed it
	// as a character using the %c format specifier.
}

// Exercise 3: Define a struct called Employee with three fields:
// firstName, lastName, and id. The first two fields are of type string,
// and the last field (id) is of type int. Create three instances of this struct
// using whatever values you’d like. Initialize the first one using the struct literal
// style without names, the second using the struct literal style with names, and
// the third with a var declaration. Use dot notation to populate the fields in the
// third struct. Print out all three structs.
func exercise3() {
	type Employee struct {
		firstName string
		lastName  string
		id        int
	}

	// Initialize the first Employee instance using struct literal without field names
	emp1 := Employee{"John", "Doe", 1}

	// Initialize the second Employee instance using struct literal with field names
	emp2 := Employee{
		firstName: "Jane",
		lastName:  "Smith",
		id:        2,
	}

	// Initialize the third Employee instance using var declaration and dot notation
	var emp3 Employee
	emp3.firstName = "Alice"
	emp3.lastName = "Johnson"
	emp3.id = 3

	// Print all three Employee instances
	fmt.Println("Employee 1:", emp1)
	fmt.Println("Employee 2:", emp2)
	fmt.Println("Employee 3:", emp3)

	// Explanation:
	// We defined the 'Employee' struct with fields 'firstName', 'lastName', and 'id'.
	// Three instances of 'Employee' were created:
	// - 'emp1' using an unnamed struct literal.
	// - 'emp2' using a named struct literal.
	// - 'emp3' using 'var' declaration and dot notation for field assignment.
	// All three instances were printed to verify their values.
}
