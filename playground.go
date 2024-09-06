package main

import "fmt"

type Inner struct {
	A int
}

func (i Inner) IntPrinter(val int) string {
	return fmt.Sprintf("Inner: %d", val)
}

func (i Inner) Double() string {
	return i.IntPrinter(i.A * 2) // Calls Inner's IntPrinter method
}

type Outer struct {
	Inner
}

func (o Outer) IntPrinter(val int) string {
	return fmt.Sprintf("Outer: %d", val)
}

func (o Outer) Double() string {
	return o.IntPrinter(o.A * 2) // Calls Inner's IntPrinter method
}

func main() {
	o := Outer{
		Inner: Inner{A: 10},
	}

	type ITest interface{}
	var test ITest = 5
	change, ok := test.(string)
	fmt.Println(change, ok)
	fmt.Println(o.Double()) // Output: "Inner: 20"
}
