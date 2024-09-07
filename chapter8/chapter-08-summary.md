## **CHAPTER 8: Generics**

"Don’t repeat yourself" is common software engineering advice. It’s better to reuse a data structure or a function than to re-create it because it’s hard to keep code changes in sync between duplicated code. In a strongly typed language like Go, the type of every function parameter and every struct field must be known at compile time. This strictness enables the compiler to help validate that your code is correct, but sometimes you’ll want to reuse the logic in a function or the fields in a struct with different types. Go provides this functionality via **type parameters**, which are colloquially referred to as **generics**.

## **Generics Reduce Repetitive Code and Increase Type Safety**

Go is a statically typed language, which means that the types of variables and parameters are checked when the code is compiled. Built-in types (maps, slices, channels) and functions (such as `len`, `cap`, or `make`) are able to accept and return values of different concrete types, but until Go 1.18, user-defined Go types and functions could not.

If you are accustomed to dynamically typed languages, where types are not evaluated until the code runs, you might not understand what the fuss is about generics, and you might be a bit unclear on what they are. It helps to think of them as **type parameters**.

#### **Example without Generics:**

```go
func divAndRemainder(num, denom int) (int, int, error) {
    if denom == 0 {
        return 0, 0, errors.New("cannot divide by zero")
    }
    return num / denom, num % denom, nil
}
```

Similarly, you create structs by specifying the type for the fields when the struct is declared:

```go
type Node struct {
    val int
    next *Node
}
```

In some situations, however, it’s useful to write functions or structs that leave the specific type of a parameter or field unspecified until it is used. This is where **generics** come into play.

### **Non-Generic Data Structure: Orderable**

Consider a case where you want to implement a binary tree for different types like `int` or `float64`. Without generics, you'd have to write separate implementations for each type or use an interface like this:

```go
type Orderable interface {
    // Order returns:
    // a value < 0 when the Orderable is less than the supplied value,
    // a value > 0 when the Orderable is greater than the supplied value,
    // and 0 when the two values are equal.
    Order(any) int
}
```

With the `Orderable` interface, we can modify our tree implementation:

```go
type Tree struct {
    val Orderable
    left, right *Tree
}

func (t *Tree) Insert(val Orderable) *Tree {
    if t == nil {
        return &Tree{val: val}
    }
    switch comp := val.Order(t.val); {
    case comp < 0:
        t.left = t.left.Insert(val)
    case comp > 0:
        t.right = t.right.Insert(val)
    }
    return t
}
```

We also need an `OrderableInt` type:

```go
type OrderableInt int

func (oi OrderableInt) Order(val any) int {
    return int(oi - val.(OrderableInt))
}
```

In the `main` function:

```go
func main() {
    var it *Tree
    it = it.Insert(OrderableInt(5))
    it = it.Insert(OrderableInt(3))
    // etc...
}
```

#### **Issue: Type Safety**

While this code works, it doesn’t allow the compiler to validate that the values inserted into your data structure are all of the same type. If you had an `OrderableString` type, you could mix `OrderableInt` and `OrderableString` in the same tree, which would cause a runtime panic.

```go
type OrderableString string

func (os OrderableString) Order(val any) int {
    return strings.Compare(string(os), val.(string))
}
```

```go
var it *Tree
it = it.Insert(OrderableInt(5))
it = it.Insert(OrderableString("nope"))  // Compiles, but panics at runtime!
```

### **Generics: A Better Solution**

With generics, you can define a data structure once and reuse it for multiple types, with the added benefit of compile-time type safety. You can create a binary tree that works with any comparable type and catch type mismatches before running the program.

Generics help you avoid runtime errors like the one above, as they allow the compiler to verify that the types being used are correct during compilation.

---

## Introducing Generics in Go

Generics have been a much-requested feature since Go's inception, and after years of research, the Go team proposed a practical solution via the Type Parameters Proposal. Generics allow for writing flexible, type-safe code without compromising on Go’s core principles of fast compilation, readable code, and efficient execution.

### Defining a Generic Stack

Generics are implemented using type parameters, which are specified in brackets. For instance, a stack can be defined generically as follows:

```go
type Stack[T any] struct {
    vals []T
}
```

In this definition:

- `T` is a type parameter, enclosed in brackets. Using capital letters for type parameters is customary.
- `any` is a type constraint using a Go interface to indicate that `T` can be any type.

### Implementing Methods Using Generics

Methods on generic types are defined similarly to methods on non-generic types but include the type parameter:

```go
func (s *Stack[T]) Push(val T) {
    s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.vals) == 0 {
        var zero T
        return zero, false
    }
    top := s.vals[len(s.vals)-1]
    s.vals = s.vals[:len(s.vals)-1]
    return top, true
}
```

### Handling Zero Values

Handling zero values in generics is slightly different because not all types can use `nil` as a zero value:

- The `var` declaration initializes the variable to the zero value of type `T`, ensuring compatibility across types.

### Using the Generic Stack

Here is how you might use the `Stack` with type `int`:

```go
func main() {
    var intStack Stack[int]
    intStack.Push(10)
    intStack.Push(20)
    intStack.Push(30)
    v, ok := intStack.Pop()
    fmt.Println(v, ok)  // Output: 30 true
}
```

Attempting to push a different type (like a string) would result in a compiler error, ensuring type safety.

### Extending Functionality with Comparable

If you need to check for equality within the stack, you must use the `comparable` interface instead of `any`:

```go
type Stack[T comparable] struct {
    vals []T
}

func (s Stack[T]) Contains(val T) bool {
    for _, v := range s.vals {
        if v == val {
            return true
        }
    }
    return false
}
```

This adjustment allows the `Contains` method to use `==` for comparison, which is not possible with the `any` type constraint due to lack of guaranteed comparability.

### Practical Example with Comparable Stack

```go
func main() {
    var s Stack[int]
    s.Push(10)
    s.Push(20)
    s.Push(30)
    fmt.Println(s.Contains(10))  // Output: true
    fmt.Println(s.Contains(5))   // Output: false
}
```

This approach demonstrates the practical use of generics to maintain type safety and flexibility, adhering to Go's design principles.

---

## Generic Functions Abstract Algorithms

Generic functions in Go provide a powerful way to abstract algorithms so they can work with different types, enhancing code reusability and maintainability. The generic implementations of `Map`, `Reduce`, and `Filter` functions are perfect examples of this capability, as outlined in the Type Parameters Proposal.

### Generic Map Function

The `Map` function transforms a slice of one type into a slice of another type using a specified function:

```go
func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
    r := make([]T2, len(s))
    for i, v := range s {
        r[i] = f(v)
    }
    return r
}
```

**Explanation**:

- `Map` has two type parameters, `T1` and `T2`, allowing it to work with slices of any types.
- `f` is a function that transforms an element of type `T1` to type `T2`.
- This function iterates over a slice of `T1`, applies `f` to each element, and stores the result in a new slice of `T2`.

### Generic Reduce Function

The `Reduce` function combines all elements of a slice into a single value using a reduction function:

```go
func Reduce[T1, T2 any](s []T1, initializer T2, f func(T2, T1) T2) T2 {
    r := initializer
    for _, v := range s {
        r = f(r, v)
    }
    return r
}
```

**Explanation**:

- `Reduce` also uses two type parameters, `T1` for the slice element type and `T2` for the result and initial value type.
- It starts with an `initializer` and iteratively applies the reduction function `f` to combine each element of the slice into this initial value.

### Generic Filter Function

The `Filter` function creates a new slice containing only those elements that match a specified condition:

```go
func Filter[T any](s []T, f func(T) bool) []T {
    var r []T
    for _, v := range s {
        if f(v) {
            r = append(r, v)
        }
    }
    return r
}
```

**Explanation**:

- `Filter` has one type parameter, `T`, and works on slices of that type.
- `f` is a predicate function that returns `true` for elements that should be included in the resulting slice.

### Example Usage

Here's how you might use these functions with a slice of strings:

```go
words := []string{"One", "Potato", "Two", "Potato"}
filtered := Filter(words, func(s string) bool {
    return s != "Potato"
})
fmt.Println(filtered)  // Output: [One Two]

lengths := Map(filtered, func(s string) int {
    return len(s)
})
fmt.Println(lengths)  // Output: [3 3]

sum := Reduce(lengths, 0, func(acc int, val int) int {
    return acc + val
})
fmt.Println(sum)  // Output: 6
```

These examples demonstrate the flexibility and type safety provided by generics, allowing you to write more abstract and reusable code. You can try these functions in the Go Playground or use them in practical projects to see how they can simplify operations on various data types.

---

## Generics and Interfaces

You can use any interface as a type constraint, not just `any` and `comparable`. For example, say you wanted to make a type that holds any two values of the same type, as long as the type implements `fmt.Stringer`. Generics make it possible to enforce this at compile time:

```go
type Pair[T fmt.Stringer] struct {
    Val1 T
    Val2 T
}
```

You can also create interfaces that have type parameters. For example, here’s an interface with a method that compares against a value of the specified type and returns a `float64`. It also embeds `fmt.Stringer`:

```go
type Differ[T any] interface {
    fmt.Stringer
    Diff(T) float64
}
```

You’ll use these two types to create a comparison function. The function takes in two `Pair` instances that have fields of type `Differ`, and returns the `Pair` with the closer values:

```go
func FindCloser[T Differ[T]](pair1, pair2 Pair[T]) Pair[T] {
    d1 := pair1.Val1.Diff(pair1.Val2)
    d2 := pair2.Val1.Diff(pair2.Val2)
    if d1 < d2 {
        return pair1
    }
    return pair2
}
```

Note that `FindCloser` takes in `Pair` instances that have fields that meet the `Differ` interface. `Pair` requires that its fields are both of the same type and that the type meets the `fmt.Stringer` interface; this function is more selective. If the fields in a `Pair` instance don’t meet `Differ`, the compiler will prevent you from using that `Pair` instance with `FindCloser`.

Now define a couple of types that meet the `Differ` interface:

```go
type Point2D struct {
    X, Y int
}

func (p2 Point2D) String() string {
    return fmt.Sprintf("{%d,%d}", p2.X, p2.Y)
}

func (p2 Point2D) Diff(from Point2D) float64 {
    x := p2.X - from.X
    y := p2.Y - from.Y
    return math.Sqrt(float64(x*x) + float64(y*y))
}

type Point3D struct {
    X, Y, Z int
}

func (p3 Point3D) String() string {
    return fmt.Sprintf("{%d,%d,%d}", p3.X, p3.Y, p3.Z)
}

func (p3 Point3D) Diff(from Point3D) float64 {
    x := p3.X - from.X
    y := p3.Y - from.Y
    z := p3.Z - from.Z
    return math.Sqrt(float64(x*x) + float64(y*y) + float64(z*z))
}
```

And here’s what it looks like to use this code:

```go
func main() {
    pair2Da := Pair[Point2D]{Point2D{1, 1}, Point2D{5, 5}}
    pair2Db := Pair[Point2D]{Point2D{10, 10}, Point2D{15, 5}}
    closer := FindCloser(pair2Da, pair2Db)
    fmt.Println(closer)

    pair3Da := Pair[Point3D]{Point3D{1, 1, 10}, Point3D{5, 5, 0}}
    pair3Db := Pair[Point3D]{Point3D{10, 10, 10}, Point3D{11, 5, 0}}
    closer2 := FindCloser(pair3Da, pair3Db)
    fmt.Println(closer2)
}
```

---

## Use Type Terms to Specify Operators

### Explanation of Generics and Operators in Go with Code Samples and Examples

This section from the book "Learning Golang" covers how to use Go's generics with operators like `/` and `%`. Here's a breakdown with additional examples and explanations to make things clearer.

#### Defining a Generic Type Constraint for Integer Types

In Go, generics allow you to write functions that work with different types, but sometimes, you need the types to support specific operators (like division and modulus). To do this, Go allows you to define **type constraints** that specify which types are allowed.

The following `Integer` interface restricts types to those that are integer types, as the modulus operator is only valid for integers.

```go
type Integer interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 | uintptr
}
```

This defines a list of all supported integer types, allowing you to use operators like `%` and `/` within a generic function.

#### Writing a Generic Function Using `Integer` Constraint

Now, you can write a generic version of the `divAndRemainder` function that works with any of the integer types defined in the `Integer` interface:

```go
func divAndRemainder[T Integer](num, denom T) (T, T, error) {
    if denom == 0 {
        return 0, 0, errors.New("cannot divide by zero")
    }
    return num / denom, num % denom, nil
}

func main() {
    var a uint = 18_446_744_073_709_551_615
    var b uint = 9_223_372_036_854_775_808
    fmt.Println(divAndRemainder(a, b))
}
```

This function divides two numbers and returns both the quotient and remainder. You can now use this function with any integer type, including `uint`, `int`, and others listed in the `Integer` interface.

### Extending the Example to Handle Custom Types

The problem arises when you try to use the `divAndRemainder` function with a custom type, even if the underlying type is an integer. Consider the following example:

```go
type MyInt int

var myA MyInt = 10
var myB MyInt = 20

fmt.Println(divAndRemainder(myA, myB))
```

This will produce a compile-time error because `MyInt` does not match exactly with the types in the `Integer` constraint. The error message suggests using `~` to allow any type whose underlying type is one of the listed types.

#### Fixing the Issue with Custom Types Using `~`

By adding a `~` before the type term, you allow types that have the specified type as their **underlying type**:

```go
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
```

Now, the function will work with both `int` and user-defined types like `MyInt`.

### More Advanced: Ordered Types and Comparison Functions

Go also supports defining generic functions that work with ordered types (types that support comparison operators like `<`, `>`, etc.). For example:

```go
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64 |
    ~string
}
```

This interface allows you to create functions that work with any type that can be compared. In Go 1.21, the `cmp` package even defines this `Ordered` interface, along with functions like `Compare` and `Less` for comparing values.

### Practical Example with a Custom Struct

You can also combine type elements and method elements in an interface. For example, suppose you want a type to have both an underlying type of `int` and a `String()` method:

```go
type PrintableInt interface {
    ~int
    String() string
}
```

This would allow you to write functions that work with any type based on `int` that also implements the `String()` method.

### Example of a Compilation Error

It's important to note that Go will let you define impossible interfaces or type constraints. For example, this would not work because the `int` type does not have a `String()` method:

```go
type ImpossiblePrintableInt interface {
    int
    String() string
}

type ImpossibleStruct[T ImpossiblePrintableInt] struct {
    val T
}
```

If you try to instantiate `ImpossibleStruct` with an `int`, you'll get a compiler error because `int` does not implement the `String()` method. Go’s type system helps you avoid such mistakes at compile time.

### Simplified Example with Slices and Maps

Finally, in addition to primitive types, Go's type constraints can also include more complex types like slices, maps, arrays, channels, and structs.

```go
type SliceOrMap interface {
    []int | map[string]int
}
```

This constraint would allow you to write a function that works with either slices of `int` or maps from `string` to `int`.

### Summary

- Use **type constraints** like `Integer` to specify the types that a generic function can work with.
- Use the `~` operator in type terms to allow user-defined types with specific underlying types.
- Go 1.21 introduced the `Ordered` interface, which makes it easier to write functions for types that support comparison.
- Be careful with impossible type constraints—Go will catch these errors at compile time.

### Step-by-Step Example

This explanation should make it easier to understand how generics and type constraints work with operators in Go. If anything is unclear or you'd like more examples, feel free to ask!

Let's work through an example to demonstrate Go's generics, type constraints, and how to work with operators like `/` and `%` in a more hands-on way. I’ll explain the steps as we go.

### Step-by-Step Example: `divAndRemainder` with Generics

#### Problem: Division and Modulus with Different Integer Types

We want to write a function that divides two numbers and also gives us the remainder (using `%`), but we want it to work with **any** integer type: `int`, `int32`, `uint`, etc. Go doesn't let you just pass any type without some way to restrict the types that work with division (`/`) and modulus (`%`).

Here's how to solve this using **generics**.

### Step 1: Define a Type Constraint for Integer Types

First, we define a **type constraint** called `Integer`, which is just an interface that lists all integer types we want to support.

```go
type Integer interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64
}
```

This says: "The `Integer` interface can be any of the listed integer types."

### Step 2: Write a Generic Function Using This Type Constraint

Now, we can write a function `divAndRemainder` that works with any type that fits the `Integer` constraint. The function will take two arguments (`num` and `denom`), perform division and modulus, and return the quotient, remainder, and an error if the denominator is zero.

Here’s how the generic function looks:

```go
func divAndRemainder[T Integer](num, denom T) (T, T, error) {
    if denom == 0 {
        return 0, 0, fmt.Errorf("cannot divide by zero")
    }
    return num / denom, num % denom, nil
}
```

#### Key Parts:

- `[T Integer]` specifies that `T` can be **any type** listed in the `Integer` interface.
- Inside the function, we use `num / denom` to get the quotient and `num % denom` to get the remainder.

### Step 3: Use the Function with Different Types

Let’s now use this function in the `main` function with different integer types.

#### Using `int`

```go
func main() {
    a := 20  // type int
    b := 3   // type int
    quotient, remainder, err := divAndRemainder(a, b)

    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Quotient: %d, Remainder: %d\n", quotient, remainder)
    }
}
```

Output:

```
Quotient: 6, Remainder: 2
```

#### Using `uint`

You can also use the same function with `uint` (unsigned integers):

```go
func main() {
    var a uint = 100  // type uint
    var b uint = 7    // type uint
    quotient, remainder, err := divAndRemainder(a, b)

    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Quotient: %d, Remainder: %d\n", quotient, remainder)
    }
}
```

Output:

```
Quotient: 14, Remainder: 2
```

### Explanation of Generics

1. **Generics** allow the function to work with multiple types without repeating code for each one.
2. The **type constraint** (`Integer`) restricts the types that the generic can work with, ensuring that the types support division (`/`) and modulus (`%`).
3. The **T** in `[T Integer]` is a placeholder for any type that satisfies the `Integer` interface. So, the function works with `int`, `uint`, `int32`, etc.

### Step 4: Handling Custom Types

Suppose you have a custom type like `MyInt` (which is just an alias for `int`), but you want to use it in the `divAndRemainder` function. This won’t work by default because Go's type system treats `MyInt` as a separate type from `int`. Here’s how we can handle that:

```go
type MyInt int

var myA MyInt = 10
var myB MyInt = 3

fmt.Println(divAndRemainder(myA, myB))  // Compile-time error
```

This will throw an error because `MyInt` is not directly listed in the `Integer` interface. To solve this, we can use the `~` symbol, which allows types with the same underlying type (like `MyInt` and `int`) to satisfy the constraint.

```go
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
```

Now, `divAndRemainder` will work with `MyInt` as well:

```go
func main() {
    var myA MyInt = 10
    var myB MyInt = 3

    quotient, remainder, err := divAndRemainder(myA, myB)

    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Quotient: %d, Remainder: %d\n", quotient, remainder)
    }
}
```

Output:

```
Quotient: 3, Remainder: 1
```

### Summary of Key Concepts

1. **Generics**: Allow you to write a single function that works with many types.
2. **Type Constraints**: Define which types are valid in a generic function (like `Integer` for types that support `%` and `/`).
3. **The `~` Operator**: Allows custom types with the same underlying type to satisfy the constraint.

By using generics, you're able to handle a wide range of types with minimal code duplication, while also enforcing the use of valid types for specific operations like division or modulus. This should help you understand how type constraints work in Go, especially in cases where operators are involved.

---

## Type Inference and Generics

Let's break down the explanation and example to make it clearer, without using personal references.

---

### Type Inference and Generics in Go

In Go, type inference is a feature that allows the compiler to determine the type automatically in many cases. For example, when declaring variables using `:=`, Go infers the type from the value:

```go
a := 10  // Go infers that 'a' is of type 'int'
```

This also applies to generic functions, where Go can often infer the types of the parameters and return values. However, type inference does not always work. Specifically, when a type parameter is used only as a return value and not as an argument, Go cannot infer the return type, so it must be explicitly specified.

### Example: Convert Function with Generics

Consider the following code that demonstrates a situation where type inference does not work:

```go
type Integer interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64
}

func Convert[T1, T2 Integer](in T1) T2 {
    return T2(in)
}

func main() {
    var a int = 10
    b := Convert[int, int64](a)  // type arguments must be specified
    fmt.Println(b)
}
```

#### Explanation:

1. The `Integer` interface restricts the types that can be used in the generic function `Convert`. It lists various integer types such as `int`, `uint`, `int64`, etc.
2. The `Convert` function is a generic function with two type parameters `T1` and `T2`, both constrained by the `Integer` interface. It takes an input of type `T1` and returns a value of type `T2`.
3. The function works by converting the input value `in` of type `T1` into the output type `T2` using a simple type conversion: `T2(in)`.

#### Why Type Inference Fails

In this case, Go cannot infer the return type (`T2`) because there is no input parameter that matches `T2`. Therefore, the types must be specified explicitly when calling the function. The following call:

```go
b := Convert[int, int64](a)
```

specifies that the input type `T1` is `int` and the return type `T2` is `int64`. Without explicitly specifying these types, the compiler would not know how to handle the conversion.

### Summary

- **Type inference** works in situations where the compiler can deduce types from function arguments.
- In cases where a type is used only in the return value (like `T2` in this example), type inference does not work, and the types must be explicitly specified.

---

## Type Elements Limit Constants

In Go, when using generics with type constraints, the compiler ensures that the operations (like adding a constant) are valid for **all** types listed in the constraint. It effectively checks against the "least capacity" type in the list, meaning that the constant must be compatible with the smallest possible type in the constraint.

### Example with the `Integer` Interface

Consider the following `Integer` interface:

```go
type Integer interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64
}
```

- **`int8`** can hold values from -128 to 127.
- **`uint8`** can hold values from 0 to 255.

Even though larger types like `int64` can hold much larger values, if a constant exceeds the capacity of the smallest types (like `int8` or `uint8`), the code will not compile.

### Invalid Example

```go
func PlusOneThousand[T Integer](in T) T {
    return in + 1_000  // INVALID
}
```

This code will **not compile** because `1_000` exceeds the maximum value of `int8` (127) and `uint8` (255). Go prevents this to avoid overflow.

### Valid Example

```go
func PlusOneHundred[T Integer](in T) T {
    return in + 100  // VALID
}
```

This will compile and work fine because `100` is within the valid range of all types, including `int8` and `uint8`.

### Summary

Go ensures that any constant used with generics must fit within the **smallest** capacity of the types listed in the type constraint. If the constant can't fit into the smallest type (like `int8` or `uint8`), the code won't compile, even if it would work for larger types.

---

## Combining Generic Functions with Generic Data Structures

Let's break this down step by step to make it easier to understand how to combine generic functions with generic data structures, like a binary tree.

### Goal

The goal is to build a **binary tree** that can work with any type of data. To do that, the tree needs a way to compare values (like numbers or struct fields) to decide where to place them in the tree or whether they exist in the tree. This is done using a **generic function** that compares two values of the same type.

### Step 1: Define a Comparison Function Type

The first step is to define a function type that can compare two values of the same type. This is done using **generics** with the `OrderableFunc` type:

```go
type OrderableFunc[T any] func(t1, t2 T) int
```

- `T any`: `T` is a placeholder for any type (this is what makes it generic).
- The function `OrderableFunc[T]` compares two values of type `T` (`t1` and `t2`) and returns an `int`.
  - It returns a negative number if `t1 < t2`.
  - It returns zero if `t1 == t2`.
  - It returns a positive number if `t1 > t2`.

### Step 2: Define the Tree Structure

Next, define a **generic binary tree**. The `Tree` struct holds two things:

1. The root node of the tree.
2. The comparison function (`OrderableFunc`).

Each **node** in the tree stores the actual value (`val`), and pointers to its left and right children (`left`, `right`).

```go
type Tree[T any] struct {
    f    OrderableFunc[T]
    root *Node[T]
}

type Node[T any] struct {
    val         T
    left, right *Node[T]
}
```

### Step 3: Create the Tree with a Comparison Function

To create a new tree, a constructor function is used. This function takes a comparison function (`OrderableFunc`) as input and returns an empty tree.

```go
func NewTree[T any](f OrderableFunc[T]) *Tree[T] {
    return &Tree[T]{f: f}
}
```

### Step 4: Adding Values to the Tree

The `Tree` needs a way to add values. When a value is added, it’s placed in the correct position in the tree using the comparison function.

- The `Add` method in `Tree` calls `Add` on the root `Node`, and the comparison function (`f`) is passed to determine where to place the value.

```go
func (t *Tree[T]) Add(v T) {
    t.root = t.root.Add(t.f, v)
}
```

### Step 5: Adding Values to a Node

The `Add` method for `Node` uses the comparison function to place values in the correct position.

- If the new value is less than the current node’s value, it goes to the left.
- If the new value is greater than the current node’s value, it goes to the right.

```go
func (n *Node[T]) Add(f OrderableFunc[T], v T) *Node[T] {
    if n == nil {
        return &Node[T]{val: v}
    }
    switch r := f(v, n.val); {
    case r <= -1:
        n.left = n.left.Add(f, v)
    case r >= 1:
        n.right = n.right.Add(f, v)
    }
    return n
}
```

### Step 6: Checking if a Value Exists in the Tree

The `Contains` method checks if a value exists in the tree, again using the comparison function.

```go
func (t *Tree[T]) Contains(v T) bool {
    return t.root.Contains(t.f, v)
}

func (n *Node[T]) Contains(f OrderableFunc[T], v T) bool {
    if n == nil {
        return false
    }
    switch r := f(v, n.val); {
    case r <= -1:
        return n.left.Contains(f, v)
    case r >= 1:
        return n.right.Contains(f, v)
    }
    return true
}
```

### Step 7: Using the Tree with Integers

Now, the tree can be used with any type, like integers. The `cmp.Compare` function is used to compare integers.

```go
t1 := NewTree(cmp.Compare[int])  // Create a tree for integers
t1.Add(10)  // Add values to the tree
t1.Add(30)
t1.Add(15)
fmt.Println(t1.Contains(15))  // true
fmt.Println(t1.Contains(40))  // false
```

### Step 8: Using the Tree with Structs

For more complex data types, like structs, a custom comparison function can be written. Here’s an example with a `Person` struct:

```go
type Person struct {
    Name string
    Age  int
}

func OrderPeople(p1, p2 Person) int {
    out := cmp.Compare(p1.Name, p2.Name)  // Compare by name first
    if out == 0 {
        out = cmp.Compare(p1.Age, p2.Age)  // If names are the same, compare by age
    }
    return out
}
```

The tree can be created with this custom comparison function:

```go
t2 := NewTree(OrderPeople)
t2.Add(Person{"Bob", 30})
t2.Add(Person{"Maria", 35})
t2.Add(Person{"Bob", 50})
fmt.Println(t2.Contains(Person{"Bob", 30}))  // true
fmt.Println(t2.Contains(Person{"Fred", 25})) // false
```

### Summary

- A **generic binary tree** can be created using **generic functions** for comparison.
- The tree works with any data type, as long as a comparison function (`OrderableFunc`) is provided.
- For integers, the `cmp.Compare` function can be used.
- For structs, a custom comparison function is needed.

This setup allows for a flexible binary tree that can handle different types of data efficiently.

---

## More on comparable

### More on `comparable` in Go

In Go, certain types are **comparable**, meaning they can be compared using the `==` and `!=` operators. However, using these operators with types that are not comparable will cause a **runtime panic**. This issue persists when using the `comparable` interface with **generics**.

### Comparable Interface

Go provides the `comparable` interface to restrict generic functions to only work with types that can be safely compared using `==` and `!=`. However, not all types are comparable, and it's important to understand when a type can and cannot be used with this interface.

### Example: Defining `Thinger` Interface and Implementations

Consider an interface `Thinger` and two implementations: `ThingerInt` and `ThingerSlice`.

```go
type Thinger interface {
    Thing()
}

type ThingerInt int

func (t ThingerInt) Thing() {
    fmt.Println("ThingInt:", t)
}

type ThingerSlice []int

func (t ThingerSlice) Thing() {
    fmt.Println("ThingSlice:", t)
}
```

- `ThingerInt`: An implementation of `Thinger` where the underlying type is `int`.
- `ThingerSlice`: Another implementation where the underlying type is `[]int`.

### Example: Generic Function with `comparable`

A generic function is defined that only accepts values that implement the `comparable` interface:

```go
func Comparer[T comparable](t1, t2 T) {
    if t1 == t2 {
        fmt.Println("equal!")
    }
}
```

This function checks if two values of type `T` are equal. The `T` type is constrained by `comparable`, so the function can only accept types that support `==` and `!=`.

### Legal Comparisons

The following comparisons work because both `int` and `ThingerInt` are **comparable** types:

```go
var a int = 10
var b int = 10
Comparer(a, b)  // prints "equal!"

var a2 ThingerInt = 20
var b2 ThingerInt = 20
Comparer(a2, b2)  // prints "equal!"
```

In this case:

- `int` is a basic type that supports `==` and `!=`.
- `ThingerInt` (which is based on `int`) is also comparable because its underlying type is `int`.

### Illegal Comparisons with Non-Comparable Types

However, the following comparison will **not compile** because `ThingerSlice` (a slice of integers) is not a comparable type:

```go
var a3 ThingerSlice = []int{1, 2, 3}
var b3 ThingerSlice = []int{1, 2, 3}
Comparer(a3, b3)  // compile fails: "ThingerSlice does not satisfy comparable"
```

Slices, maps, and functions are **not comparable** in Go, so the code fails to compile.

### Pitfall: Interfaces and Comparability

Things get tricky when using interface types like `Thinger`. The following example works because `ThingerInt` is a comparable type:

```go
var a4 Thinger = a2
var b4 Thinger = b2
Comparer(a4, b4)  // prints "equal!"
```

However, if `ThingerSlice` is assigned to variables of type `Thinger`, the code compiles, but it will **panic at runtime**:

```go
a4 = a3
b4 = b3
Comparer(a4, b4)  // compiles, panics at runtime
```

The reason for the panic is that the underlying type of `a4` and `b4` is `ThingerSlice`, which is not comparable. Although the code compiles, Go will throw a **runtime panic** when it tries to compare `ThingerSlice`.

### Why the Panic Occurs

Go allows assigning any type that implements an interface to a variable of that interface type (in this case, `Thinger`). But not all types that implement an interface are comparable. If the underlying type of the interface is non-comparable (like `ThingerSlice`), the program will panic when comparison operators are used.

Let's explain why the panic occurs when comparing values of a non-comparable type, like `ThingerSlice`, through a clearer step-by-step process.

### What’s Happening with Interfaces and Comparability

In Go, **interfaces** define a set of methods that a type must implement. Any type that implements those methods can be assigned to a variable of the interface type.

#### Step 1: The Interface `Thinger`

Consider the `Thinger` interface, which has a single method `Thing()`:

```go
type Thinger interface {
    Thing()
}
```

Any type that implements the `Thing()` method can be assigned to a variable of type `Thinger`.

#### Step 2: Implementing the `Thinger` Interface

Now, there are two types that implement this interface:

```go
type ThingerInt int

func (t ThingerInt) Thing() {
    fmt.Println("ThingInt:", t)
}

type ThingerSlice []int

func (t ThingerSlice) Thing() {
    fmt.Println("ThingSlice:", t)
}
```

- `ThingerInt` implements `Thinger`, and its underlying type is `int` (which is **comparable**).
- `ThingerSlice` also implements `Thinger`, but its underlying type is `[]int` (which is **not comparable**).

#### Step 3: Assigning to `Thinger` Variables

Because both `ThingerInt` and `ThingerSlice` implement the `Thinger` interface, they can be assigned to variables of type `Thinger`:

```go
var a4 Thinger = ThingerInt(10)  // OK, since ThingerInt is comparable
var a5 Thinger = ThingerSlice{1, 2, 3}  // OK, since ThingerSlice also implements Thinger
```

Even though `a4` and `a5` are both of type `Thinger`, their underlying types are different (`ThingerInt` and `ThingerSlice`).

#### Step 4: Using `comparable` with Interfaces

Now, when a function like `Comparer[T comparable]` is used, it expects the values of `T` to be **comparable** (i.e., the `==` and `!=` operators can be used on them).

```go
func Comparer[T comparable](t1, t2 T) {
    if t1 == t2 {
        fmt.Println("equal!")
    }
}
```

When this function is called with `a4` (of type `ThingerInt`), it works fine because `ThingerInt` is based on `int`, which is comparable.

```go
var a4 Thinger = ThingerInt(10)
var b4 Thinger = ThingerInt(10)
Comparer(a4, b4)  // Works, prints "equal!"
```

#### Step 5: The Problem with `ThingerSlice`

The issue arises when `ThingerSlice` is assigned to `a4` and `b4`. Even though `ThingerSlice` implements `Thinger`, it is **not comparable** because slices in Go cannot be compared using `==`.

```go
var a4 Thinger = ThingerSlice{1, 2, 3}
var b4 Thinger = ThingerSlice{1, 2, 3}
Comparer(a4, b4)  // This causes a runtime panic
```

- At compile time, the Go compiler allows this because `a4` and `b4` are both of type `Thinger`, and the `Thinger` interface can technically be used with `Comparer`.
- However, at **runtime**, Go realizes that `a4` and `b4` are actually `ThingerSlice` (whose underlying type is `[]int`), which is **not comparable**.

When Go tries to compare the `ThingerSlice` values using `==`, it results in a **runtime panic** with the message:

```
panic: runtime error: comparing uncomparable type main.ThingerSlice
```

### Why the Panic Occurs

- **Interfaces in Go** allow assigning any type that implements the required methods to a variable of the interface type.
- **Not all types that implement an interface are comparable**. For example, slices (`[]int`) are not comparable using `==`.
- When using generics constrained by `comparable`, the actual underlying type of the interface variable is checked at **runtime**.
- If the underlying type is non-comparable (like a slice), Go will allow the code to compile, but it will **panic at runtime** when it tries to compare the values using `==`.

### Summary

- The panic happens because `ThingerSlice`, though it implements the `Thinger` interface, has an underlying type (`[]int`) that is **not comparable**.
- The **compiler** doesn’t check the underlying type when using an interface with generics, but at **runtime**, Go sees that slices can’t be compared, leading to a panic.

### Conclusion

- **Comparable types**: Only certain types in Go are comparable (e.g., basic types like `int`, `float64`, and custom types with comparable underlying types).
- **Non-comparable types**: Slices, maps, and functions are not comparable. Using them in a generic function constrained by `comparable` will result in a compile-time or runtime error.
- **Runtime panic**: When using interfaces, care must be taken because the underlying type may not be comparable, leading to runtime panics if comparison operators are used.

For a deeper dive into why this design choice was made in Go, the blog post ["All Your Comparable Types"](https://blog.golang.org/comparable) by Robert Griesemer provides more technical details on how comparable types and generics interact.

---

## Things That Are Left Out

### Things That Are Left Out of Go Generics

Go remains a small and focused language, and its implementation of generics deliberately leaves out some features commonly found in other languages. This section outlines several features that are not part of Go’s generics, and why they have been excluded.

#### Operator Overloading

In languages like Python, Ruby, and C++, **operator overloading** allows user-defined types to define custom implementations for operators (such as `+`, `-`, `[]`, etc.). Go does not have and will not add operator overloading.

This means:

- User-defined container types cannot use **range** for iteration or **[]** for indexing.
- **Operator overloading** can lead to less readable code, as developers may assign unexpected or inconsistent meanings to operators. For example, in C++, `<<` can mean "shift bits left" for integers or "write the value on the right to the left" for streams.

By avoiding operator overloading, Go aims to keep code **more predictable and readable**.

#### No Method Type Parameters

Go generics do not support additional type parameters on methods. Consider this example, where `Map` and `Reduce` functions are implemented as methods of a generic `functionalSlice` type:

```go
type functionalSlice[T any] []T

// THIS DOES NOT WORK
func (fs functionalSlice[T]) Map[E any](f func(T) E) functionalSlice[E] {
    out := make(functionalSlice[E], len(fs))
    for i, v := range fs {
        out[i] = f(v)
    }
    return out
}

// THIS DOES NOT WORK
func (fs functionalSlice[T]) Reduce[E any](start E, f func(E, T) E) E {
    out := start
    for _, v := range fs {
        out = f(out, v)
    }
    return out
}
```

The intent here is to chain methods together, but this does not work in Go. Instead, calls must be **nested** or **invoked separately** with intermediate variables. While chaining method calls like this might appeal to fans of functional programming, Go does not support **parameterized methods**. The Go generics proposal discusses the reasoning for excluding this feature, primarily to keep things simple and to avoid chaining complexity.

#### No Variadic Type Parameters

Go generics do not support **variadic type parameters**. In regular Go functions, variadic input parameters can be used to accept a varying number of arguments. However, for generics, there is no way to define a type pattern for those variadic parameters. For example, you cannot specify that the parameters alternate between `string` and `int`.

Variadic parameters must all match a **single declared type**, which can be generic or non-generic, but there is no flexibility for varying type patterns.

#### Other Features Left Out

Some additional features, commonly seen in other languages, have been omitted from Go’s generics. These include:

- **Specialization**: This feature allows a function or method to have type-specific versions in addition to the generic version. Since Go does not support function or method overloading, this feature is not being considered.
- **Currying**: This allows a function or type to be partially instantiated based on another generic function or type by specifying some of the type parameters. Go does not have currying.

- **Metaprogramming**: Metaprogramming allows code to be generated and run at compile-time to produce code that runs at runtime. Go avoids this type of complexity to keep the language more straightforward.

### Conclusion

Go’s implementation of generics is intentionally minimal, leaving out several features that are present in other languages. Features like **operator overloading**, **method type parameters**, **variadic type patterns**, and more **esoteric features** such as **specialization**, **currying**, and **metaprogramming** are excluded to keep Go simple, predictable, and readable. The focus of Go’s generics is on providing essential functionality without introducing complexity or reducing code readability.

### Simple explanation

### No Method Type Parameters (Simplified)

In Go, **type parameters** can be used for types and functions, but **they cannot be added to individual methods**. This means you can't define extra type parameters just for a method of a type. Let's break this down with a simpler explanation and example.

### What Are Type Parameters on Methods?

In some programming languages, it's possible to add extra type parameters to methods within a generic type. This allows methods to operate on different types, even if the type itself is working with a specific type.

Go does **not** allow this feature. You can have generic types and generic functions, but you cannot add extra type parameters to methods of those types.

#### Example: Why This Doesn't Work in Go

Let's say we have a **slice** type that stores a list of values, and we want to define a method on it that converts each item in the slice to another type using a function. This method would have its own extra type parameter for the new type.

```go
type MySlice[T any] []T

// THIS DOES NOT WORK IN GO
func (ms MySlice[T]) ConvertTo[E any](f func(T) E) MySlice[E] {
    result := make(MySlice[E], len(ms))
    for i, v := range ms {
        result[i] = f(v)
    }
    return result
}
```

- `T`: The type of the elements in the `MySlice`.
- `E`: A **new type** that we want to convert the slice elements into, using the function `f`.

This code would try to convert a slice of one type (`T`) to another type (`E`) using the method `ConvertTo`. But this **doesn't work** in Go because Go does not allow adding extra type parameters (`E`) to methods.

### A Simpler Example

Let's use a simpler example with integers and strings. Suppose we want a method that converts a slice of integers into a slice of strings, but we want to keep the type of the original slice generic.

```go
type MySlice[T any] []T

// THIS DOES NOT WORK IN GO
func (ms MySlice[int]) ToStrings() MySlice[string] {
    result := make(MySlice[string], len(ms))
    for i, v := range ms {
        result[i] = fmt.Sprintf("%d", v) // Convert int to string
    }
    return result
}
```

Here:

- `MySlice[int]` is a slice of integers.
- `ToStrings` is supposed to convert each `int` into a `string`.

This would require **type parameters** for the method (`ToStrings`), but Go doesn't allow this. So, Go prevents this from compiling.

### How to Work Around This in Go

Since Go doesn't allow extra type parameters on methods, you need to use **regular functions** instead. You can still use generics with functions, but they can't be chained as methods on a type.

Here’s how the same thing would work with a function:

```go
// A generic function to convert a slice of one type to another
func ConvertTo[T, E any](ms []T, f func(T) E) []E {
    result := make([]E, len(ms))
    for i, v := range ms {
        result[i] = f(v)
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3}
    strings := ConvertTo(numbers, func(n int) string {
        return fmt.Sprintf("%d", n)
    })
    fmt.Println(strings) // Output: ["1", "2", "3"]
}
```

In this function:

- `ConvertTo` is a generic function that works with any slice of type `T` and converts it to a slice of type `E` using a conversion function `f`.
- We convert integers to strings in this example.

### Summary

Go **does not support method-specific type parameters**, meaning you can't define extra generic types just for methods. Instead, you can use generic functions that take type parameters and work with different types outside the context of the method. This keeps Go's generics simpler and avoids adding extra complexity to methods.

---

## Idiomatic Go and Generics

### Impact of Generics on Idiomatic Go

With the introduction of **generics**, some best practices for writing Go have changed. However, it is important to note that older Go code without generics still works perfectly fine. This section discusses how generics affect idiomatic Go and provides guidance on when to use them.

### Key Changes with Generics

1. **No More `float64` for All Numbers**:
   Before generics, it was common to use `float64` to represent any numeric type in Go. Now, generics allow for more flexibility by letting functions and data structures handle multiple numeric types without relying on `float64`.

2. **Use `any` Instead of `interface{}`**:
   The new keyword `any` replaces the older `interface{}` when representing an unspecified type. For example, use `any` for functions or data structures that can work with any type.

   Old way:

   ```go
   func printValue(v interface{}) {
       fmt.Println(v)
   }
   ```

   New way:

   ```go
   func printValue(v any) {
       fmt.Println(v)
   }
   ```

3. **Handle Multiple Slice Types with One Function**:
   With generics, it’s now possible to write a single function that handles slices of various types, removing the need to write separate functions for different types of slices.

### Performance Impact of Generics

The performance impact of generics is still evolving. As of Go 1.20, there’s no significant compilation-time penalty for using generics, and the runtime performance has been improving. However, there are cases where generics might be **slower** than traditional Go code.

#### Example: Generics vs. Interface in a Simple Function

Consider the following two functions:

1. **Using an Interface**:

   ```go
   type Ager interface {
       age() int
   }

   func doubleAge(a Ager) int {
       return a.age() * 2
   }
   ```

2. **Using a Generic Type**:
   ```go
   func doubleAgeGeneric[T Ager](a T) int {
       return a.age() * 2
   }
   ```

In Go 1.20, the **generic version** of this function is about **30% slower** than the interface version for simple cases. This is due to how the Go compiler currently generates code for generics:

- The Go compiler creates unique functions for different **underlying types**, but pointer types share a single generated function.
- This means that for shared generic functions, the compiler adds **runtime lookups** to handle different types, which slows down performance.

However, for **nontrivial functions**, there’s generally no significant performance difference between the two approaches. The key takeaway is that switching to generics **solely for performance reasons** may not always result in faster code.

### Comparison with Other Languages

In languages like **C++**, generics often improve performance by generating a **unique function for each concrete type** at compile time. This approach avoids runtime lookups and makes the resulting binary faster, but it also increases binary size.

Go’s current generics implementation, however, shares some functions between different types, particularly pointer types, and this introduces **runtime overhead**. Over time, Go's generics implementation is expected to become more efficient, and the performance impact may reduce.

### Best Practices for Using Generics

1. **Use Generics for Flexibility, Not Performance**:
   Don’t rush to convert every function to use generics. The primary benefit of generics is making code more flexible and reusable, not necessarily faster.

2. **Test and Benchmark**:
   Use Go’s **benchmarking** and **profiling** tools to measure the impact of generics on your code. Write maintainable code that’s fast enough for your use case, rather than assuming generics will always result in performance improvements.

3. **Expect Future Improvements**:
   As Go’s generics implementation matures, both the **compilation speed** and **runtime performance** of generic code are likely to improve in future Go versions.

### Summary

- Generics bring flexibility to Go by allowing code to handle multiple types without using `interface{}` or `float64` for generic cases.
- While generics can sometimes be slower than traditional interface-based functions, this is mainly due to how the Go compiler handles shared functions for different types.
- Don’t feel pressured to convert all existing Go code to use generics immediately. Write code that is clear, maintainable, and meets performance needs.

---

## Adding Generics to the Standard Library

### Adding Generics to the Standard Library

When Go first introduced generics in **Go 1.18**, the release was quite conservative. While the new types `any` and `comparable` were added, no changes were made to the standard library to support generics at that time. However, generics are now more widely adopted in the Go community, and starting with **Go 1.21**, the standard library includes functions that make use of generics to implement common algorithms for working with slices, maps, and concurrency.

### Changes in Go 1.21

Here are some of the key additions to the standard library that make use of generics:

1. **`any` Replaces `interface{}`**:

   - In Go 1.21, almost all instances of `interface{}` in the standard library were replaced with `any`. This makes code easier to read and aligns with the new generic syntax.

2. **New Generic Functions for Slices and Maps**:

   - The **slices** and **maps** packages now include generic functions that simplify working with these data structures.
   - For example, `Equal` and `EqualFunc` were added to make comparing slices easier.

3. **Slice Manipulation**:

   - Functions like `Insert`, `Delete`, and `DeleteFunc` in the **slices** package allow developers to avoid manually writing complex slice-handling code. These functions use generics to handle slices of any type.

4. **Map Cloning**:

   - The **maps.Clone** function was introduced to create shallow copies of maps efficiently. It leverages Go’s runtime features to provide a faster, safer way to clone maps, using generics to handle maps with any key-value type.

5. **Concurrency Enhancements with `sync.OnceValue`**:
   - The **sync** package now includes **`sync.OnceValue`** and **`sync.OnceValues`**, which use generics to create functions that are invoked only once and return one or two values. This is useful for scenarios where a value should only be computed or fetched once in a thread-safe manner.

### Future of Generics in the Standard Library

As the Go community becomes more comfortable with generics, it is expected that more parts of the standard library will adopt them. The goal is to provide reusable, efficient, and easy-to-read generic functions for common tasks so that developers don't have to write their own implementations for things like slice and map manipulation.

In future releases of Go, expect even more generic functions and types in the standard library to simplify common programming tasks.

---

## Exercises

Here are the exercises along with explanations and starter code for solving the problems using generics in Go:

### Exercise 1: Double the Value of Integers and Floats

Write a generic function that doubles the value of any integer or float that’s passed in to it.

#### Solution Outline:

- Define an interface that includes both integer and floating-point types.
- Write a function that accepts a value of the generic type and doubles it.

```go
// Define a generic interface that includes integer and floating-point types.
type Number interface {
    int | int8 | int16 | int32 | int64 | float32 | float64
}

// Generic function that doubles the value.
func Double[T Number](value T) T {
    return value * 2
}

func main() {
    // Example usage:
    fmt.Println(Double(10))      // int
    fmt.Println(Double(10.5))    // float64
    fmt.Println(Double(int64(8))) // int64
}
```

### Exercise 2: Printable Interface

Define a generic interface called `Printable` that matches a type implementing `fmt.Stringer` and having an underlying type of `int` or `float64`. Then write a function that takes in a `Printable` and prints its value using `fmt.Println`.

#### Solution Outline:

- Use the `fmt.Stringer` interface and restrict the underlying type to `int` and `float64`.
- Write types that implement the `fmt.Stringer` interface and fulfill the `Printable` constraint.

```go
// Define Printable interface that implements fmt.Stringer and has underlying types of int or float64.
type Printable interface {
    fmt.Stringer
    ~int | ~float64
}

// Define a type that implements fmt.Stringer and has an underlying int type.
type MyInt int

func (m MyInt) String() string {
    return fmt.Sprintf("MyInt value: %d", m)
}

// Define a type that implements fmt.Stringer and has an underlying float64 type.
type MyFloat float64

func (m MyFloat) String() string {
    return fmt.Sprintf("MyFloat value: %.2f", m)
}

// Function that takes in a Printable and prints its value using fmt.Println.
func PrintValue[T Printable](value T) {
    fmt.Println(value)
}

func main() {
    // Example usage:
    var myInt MyInt = 42
    var myFloat MyFloat = 42.42
    PrintValue(myInt)
    PrintValue(myFloat)
}
```

### Exercise 3: Singly Linked List Data Type

Write a generic singly linked list data type. Each element holds a comparable value and has a pointer to the next element in the list.

#### Solution Outline:

- Use a generic type for the linked list nodes and the list itself.
- Implement methods to add elements, insert elements at a specific position, and return the index of a value.

```go
// Node represents a single element in the linked list.
type Node[T comparable] struct {
    value T
    next  *Node[T]
}

// LinkedList represents the linked list itself.
type LinkedList[T comparable] struct {
    head *Node[T]
}

// Add method adds a new element to the end of the linked list.
func (l *LinkedList[T]) Add(value T) {
    newNode := &Node[T]{value: value}
    if l.head == nil {
        l.head = newNode
        return
    }
    current := l.head
    for current.next != nil {
        current = current.next
    }
    current.next = newNode
}

// Insert method adds an element at the specified position in the linked list.
func (l *LinkedList[T]) Insert(value T, position int) {
    newNode := &Node[T]{value: value}
    if position == 0 {
        newNode.next = l.head
        l.head = newNode
        return
    }
    current := l.head
    for i := 0; current != nil && i < position-1; i++ {
        current = current.next
    }
    if current == nil {
        return
    }
    newNode.next = current.next
    current.next = newNode
}

// Index method returns the position of the supplied value, -1 if it's not present.
func (l *LinkedList[T]) Index(value T) int {
    current := l.head
    position := 0
    for current != nil {
        if current.value == value {
            return position
        }
        current = current.next
        position++
    }
    return -1
}

func main() {
    list := LinkedList[int]{}

    // Add elements to the list.
    list.Add(10)
    list.Add(20)
    list.Add(30)

    // Insert element at position 1.
    list.Insert(15, 1)

    // Find index of a value.
    fmt.Println(list.Index(20)) // Output: 2
    fmt.Println(list.Index(40)) // Output: -1
}
```

### Summary of Exercises

1. **Double Function**: Write a function that doubles the value of an integer or float using a generic interface for numeric types.
2. **Printable Interface**: Define a `Printable` interface that implements `fmt.Stringer` and has an underlying type of `int` or `float64`, then print its value.
3. **Singly Linked List**: Implement a generic singly linked list that can add elements, insert them at specific positions, and return the index of a given value.
