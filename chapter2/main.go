package main

import (
	"fmt"
	"math/cmplx"
)

func main() {
	// 1. Predeclared Types
	// Boolean type
	var isActive bool = true // Explicit declaration
	var isClosed bool        // Zero value: false
	fmt.Println("Boolean:", isActive, isClosed)

	// Integer types
	var smallInt int8 = -128                    // 8-bit signed integer
	var largeUint uint64 = 18446744073709551615 // 64-bit unsigned integer
	fmt.Println("Integers:", smallInt, largeUint)

	// Float types
	var pi float64 = 3.14159 // 64-bit floating-point number
	fmt.Println("Float:", pi)

	// Complex types
	// As was mentioned in the book you do not need to learn this if you not working with it
	var complexNum complex128 = cmplx.Sqrt(-5 + 12i) // Complex number
	fmt.Println("Complex:", complexNum)

	// String and Rune types
	var greeting string = "Hello, Go!" // String
	var char rune = 'G'                // Rune (alias for int32)
	fmt.Println("String and Rune:", greeting, string(char))

	// 2. Zero Value
	// Variables without initialization get their zero value
	var defaultInt int       // Zero value: 0
	var defaultFloat float64 // Zero value: 0.0
	var defaultBool bool     // Zero value: false
	var defaultString string // Zero value: "" (empty string)
	fmt.Println("Zero Values:", defaultInt, defaultFloat, defaultBool, defaultString)

	// 3. Literals
	// Integer literals with different bases and underscores for readability
	var dec int = 123               // Decimal
	var bin int = 0b1010            // Binary
	var oct int = 0o12              // Octal
	var hex int = 0x1A              // Hexadecimal
	var readableInt int = 1_000_000 // Readable integer with underscores
	fmt.Println("Literals:", dec, bin, oct, hex, readableInt)

	// Floating-point and complex literals
	var sci float64 = 1.2e3        // Scientific notation
	var hexFloat float64 = 0x1.2p3 // Hexadecimal floating-point
	fmt.Println("Floating-Point Literals:", sci, hexFloat)

	// String literals: interpreted and raw
	var interpString string = "Hello\nWorld" // Interpreted string
	var rawString string = `Hello\nWorld`    // Raw string
	fmt.Println("Strings:", interpString, rawString)

	// 4. Variable Declarations
	// Using var keyword
	var age int = 30

	// Using short declaration
	name := "Go Developer"

	// Declaring constants
	const Pi = 3.14159            // Untyped constant
	const Greeting = "Hello, Go!" // Typed constant: string
	fmt.Println("Variables and Constants:", age, name, Pi, Greeting)

	// 5. Typed vs. Untyped Constants
	const untyped = 42               // Untyped constant
	var typedFloat float64 = untyped // Used as float64 without explicit conversion
	fmt.Println("Typed vs. Untyped:", untyped, typedFloat)

	// 6. Explicit Type Conversion
	var a int = 10
	var b float64 = float64(a) // Explicit conversion from int to float64
	var c uint = uint(b)       // Explicit conversion from float64 to uint
	fmt.Println("Type Conversions:", a, b, c)

	// 7. Common Pitfalls and Best Practices
	// Unused variables - Uncommenting below lines will cause a compile error due to unused variable
	// var unusedVar int = 50

	// Implicit types - Beware of potential type issues
	const implicitConst = 5                         // Untyped
	var implicitTyped float64 = implicitConst + 0.5 // Works because of compatible context
	fmt.Println("Implicit Constant:", implicitTyped)
}
