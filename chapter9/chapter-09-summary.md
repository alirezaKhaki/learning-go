# CHAPTER 9 Errors

## Error Handling in Go

This document explores error handling in Go, a fundamental concept for developers transitioning from languages with exceptions.

### Go's Approach to Errors

Go employs a distinct approach to error handling compared to languages with exceptions. Here's the core principle:

- Functions **return an error value** of type `error` as the last return value to signal an issue.
- The calling function checks this error by comparing it to `nil`. If `nil`, there's no error. Otherwise, an error has occurred.

This method fosters clear and explicit error handling, promoting robust code.

### Error Handling Example

Let's dissect a basic function demonstrating error handling:

```go
func calcRemainderAndMod(numerator, denominator int) (int, int, error) {
  if denominator == 0 {
    return 0, 0, errors.New("denominator is 0")
  }
  return numerator / denominator, numerator % denominator, nil
}
```

- This function calculates the remainder and modulo of two integers.
- It checks for division by zero, a potential error.
- If `denominator` is zero, an error object is created using `errors.New("denominator is 0")`.
- On successful execution, `nil` is returned for the error.

### Handling Errors in the Calling Function

The calling function retrieves the error value and acts accordingly:

```go
func main() {
  numerator := 20
  denominator := 3
  remainder, mod, err := calcRemainderAndMod(numerator, denominator)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  fmt.Println(remainder, mod)
}
```

- This `main` function calls `calcRemainderAndMod`.
- It retrieves the returned values, including the error (`err`).
- An `if` statement checks if `err` is not `nil` (meaning an error occurred).
- If there's an error, it's printed, and the program exits with an error code (`os.Exit(1)`).
- If no error occurs, the remainder and modulo are printed.

**Key Points:**

- Go's error handling promotes explicit checks using `if` statements.
- The error handling code is indented for better readability, separating it from the main logic.

### The `error` Interface

The `error` type in Go is an interface that defines a single method:

```go
type error interface {
  Error() string
}
```

Any type that implements this method can be considered an error. Returning `nil` from a function signifies successful execution (no error).

### Why Go Avoids Exceptions

Go steers clear of exceptions for two primary reasons:

1. **Clarity and Control Flow:** Exceptions introduce hidden code paths, making it difficult to track error handling. This can lead to unexpected crashes or incorrect program behavior.
2. **Enforced Error Handling:** Go's approach compels developers to explicitly handle errors or acknowledge their ignorance using an underscore (`_`) for the error variable. This promotes code that's more maintainable and less prone to errors.

By prioritizing clear and explicit error handling, Go fosters robust and well-structured code.

---

## Using Strings for Simple Errors in Go

This section explores creating errors from strings in Go, offering two methods from the standard library:

### 1. `errors.New` Function

- Creates an error object from a provided string.
- The string becomes the error message accessible through the `Error` method.
- `fmt.Println` automatically calls `Error` when printing an error object.

**Example:**

```go
func doubleEven(i int) (int, error) {
  if i % 2 != 0 {
    return 0, errors.New("only even numbers are processed")
  }
  return i * 2, nil
}

func main() {
  result, err := doubleEven(1)
  if err != nil {
    fmt.Println(err) // prints "only even numbers are processed"
  }
  fmt.Println(result)
}
```

### 2. `fmt.Errorf` Function

- Creates an error object with a formatted string message.
- Uses `fmt.Printf` verbs to dynamically include data in the error message.
- Similar to `errors.New`, the formatted string is accessible through `Error`.

**Example:**

```go
func doubleEven(i int) (int, error) {
  if i % 2 != 0 {
    return 0, fmt.Errorf("%d isn't an even number", i)
  }
  return i * 2, nil
}
```

**Choosing the Right Method:**

- For simple, pre-defined error messages, `errors.New` is sufficient.
- When you need dynamic information (like a specific value) in the error message, use `fmt.Errorf` for a more informative message.

**Benefits of String-Based Errors:**

- Simplicity and ease of use.
- Suitable for conveying clear and concise error messages.

**Limitations:**

- Lack context for complex errors (e.g., location of the error in the code).

---

## Sentinel Errors

Sentinel errors signal that processing cannot continue due to a problem with the current state. The term "sentinel errors" was coined by Dave Cheney, a well-known developer in the Go community, in his blog post "Don't Just Check Errors, Handle Them Gracefully." The concept stems from the practice of using specific values to signify that no further processing is possible.

In Go, sentinel errors are typically package-level variables and are conventionally named starting with `Err` (with exceptions like `io.EOF`). These errors should be treated as read-only. Although the Go compiler cannot enforce this, altering their value is considered a programming error.

Sentinel errors are commonly used to indicate that processing cannot proceed. For example, the `archive/zip` package defines several sentinel errors, such as `ErrFormat`, which occurs when data that does not represent a ZIP file is passed. Here's a sample code demonstrating this:

```go
package main

import (
    "bytes"
    "fmt"
    "archive/zip"
)

func main() {
    data := []byte("This is not a zip file")
    notAZipFile := bytes.NewReader(data)
    _, err := zip.NewReader(notAZipFile, int64(len(data)))
    if err == zip.ErrFormat {
        fmt.Println("Told you so")
    }
}
```

In this example, the error `zip.ErrFormat` is checked using `==` to determine if the passed data is not a valid ZIP file. This pattern is commonly used to handle sentinel errors in Go.

#### Example of Another Sentinel Error

Another sentinel error is `rsa.ErrMessageTooLong`, defined in the `crypto/rsa` package. It indicates that a message cannot be encrypted because it exceeds the allowable length for the provided public key. Similarly, `context.Canceled`, a sentinel error in the `context` package, is covered in Chapter 14.

#### Defining Sentinel Errors

Before defining a sentinel error, ensure that it is necessary. Once defined, it becomes part of the public API, meaning that future backward-compatible releases must include it. Often, it's better to reuse existing sentinel errors from the standard library or to define a custom error type that provides more context about the condition causing the error.

Sentinel errors are best used when they indicate a specific state where further processing is impossible, and no additional information is required to explain the error. For testing sentinel errors, use the `==` operator when the function's documentation explicitly mentions that it returns a sentinel error.

#### Using Constants for Sentinel Errors

Dave Cheney suggested that constants could be useful for sentinel errors. Here’s an example:

```go
package consterr

type Sentinel string

func (s Sentinel) Error() string {
    return string(s)
}
```

This code creates a type that implements the `error` interface. You could then define sentinel errors like this:

```go
package mypkg

const (
    ErrFoo = consterr.Sentinel("foo error")
    ErrBar = consterr.Sentinel("bar error")
)
```

While this approach might seem like a good solution at first, it isn't considered idiomatic in Go. If two constant errors from different packages had the same error string, they would be considered equal, which can lead to unintended behavior. Errors created with `errors.New` are only equal to themselves or variables explicitly assigned their value, which prevents this issue.

The use of sentinel errors in Go aligns with Go’s design philosophy: simplicity and trust in the developers. Sentinel errors should be rare and handled by convention, rather than enforcing strict language rules. Though these errors are mutable, accidental reassignment is unlikely due to Go’s package-level variable handling.

#### Note:

Sentinel errors should not be used frequently. If the error state requires additional context, it's better to use custom error types. This will allow more information to be passed about the error condition, making error handling and debugging easier. This will be discussed further in the next section.

---

## Errors Are Values

Since `error` in Go is an interface, it allows developers to define custom errors that can include additional information for logging or error handling. For instance, adding a status code to an error can help in classifying the error without relying on string comparisons, which are prone to change. Below is an example of how to achieve this.

#### Step 1: Defining Status Codes

First, define an enumeration to represent different status codes. The `iota` keyword is used to create unique constants for different error statuses.

```go
type Status int

const (
    InvalidLogin Status = iota + 1
    NotFound
)
```

#### Step 2: Creating a Custom Error Type

Define a struct `StatusErr` to hold the status code and an error message. This custom error type implements the `error` interface by defining the `Error()` method, which returns the error message.

```go
type StatusErr struct {
    Status  Status
    Message string
}

func (se StatusErr) Error() string {
    return se.Message
}
```

#### Step 3: Using the Custom Error Type

Now, use the `StatusErr` in a function to provide more details about what went wrong. The example below shows how to use `StatusErr` when an invalid login occurs or when a file is not found:

```go
func LoginAndGetData(uid, pwd, file string) ([]byte, error) {
    token, err := login(uid, pwd)
    if err != nil {
        return nil, StatusErr{
            Status:  InvalidLogin,
            Message: fmt.Sprintf("invalid credentials for user %s", uid),
        }
    }
    data, err := getData(token, file)
    if err != nil {
        return nil, StatusErr{
            Status:  NotFound,
            Message: fmt.Sprintf("file %s not found", file),
        }
    }
    return data, nil
}
```

In this case, the `StatusErr` provides detailed information about the error, such as whether it was due to an invalid login or a missing file.

#### Important Note

Even when defining custom error types, always return `error` as the function's error return type. This allows different error types to be returned from a function and enables the caller to handle errors without relying on specific types.

#### Potential Pitfall: Uninitialized Custom Error

When using custom error types, avoid returning an uninitialized instance. Consider the following broken example:

```go
func GenerateErrorBroken(flag bool) error {
    var genErr StatusErr
    if flag {
        genErr = StatusErr{
            Status: NotFound,
        }
    }
    return genErr
}

func main() {
    err := GenerateErrorBroken(true)
    fmt.Println("GenerateErrorBroken(true) returns non-nil error:", err != nil)
    err = GenerateErrorBroken(false)
    fmt.Println("GenerateErrorBroken(false) returns non-nil error:", err != nil)
}
```

Output:

```
true
true
```

Even though the error should be `nil` in the second case, it isn’t. This happens because `error` is an interface. As mentioned in the book’s section on "Interfaces and nil," both the underlying type and value must be `nil` for an interface to be considered `nil`. In this case, the underlying type is not `nil`.

#### Fixing the Issue

There are two ways to fix this:

1. **Return `nil` explicitly when no error occurs:**
   This is the most common approach, as it ensures that an actual `nil` value is returned when no error occurs.

   ```go
   func GenerateErrorOKReturnNil(flag bool) error {
       if flag {
           return StatusErr{
               Status: NotFound,
           }
       }
       return nil
   }
   ```

2. **Use a variable of type `error` to hold the error:**
   By defining the error variable as `error`, it ensures the correct handling of the error.

   ```go
   func GenerateErrorUseErrorVar(flag bool) error {
       var genErr error
       if flag {
           genErr = StatusErr{
               Status: NotFound,
           }
       }
       return genErr
   }
   ```

#### Best Practice

When using custom errors, it’s important not to define local variables of the custom error type. Instead, use the `error` type for error variables or return `nil` explicitly when no error occurs. Additionally, avoid using type assertions or switches to access the fields of a custom error. Instead, use `errors.As`, which will be discussed later in the book. This is the idiomatic way to check for and handle custom errors in Go.

---

## Wrapping Errors

In Go, when an error is propagated through your code, it's often useful to add additional context, such as the function where the error occurred or details about the failed operation. This practice is called **wrapping** an error. When multiple errors are wrapped together, this structure is referred to as an **error tree**.

#### Wrapping Errors with `fmt.Errorf`

The Go standard library provides a built-in way to wrap errors using `fmt.Errorf` and the `%w` verb. This verb preserves the original error while appending a new message, allowing the original error to be "unwrapped" later. Here’s an example:

```go
func fileChecker(name string) error {
    f, err := os.Open(name)
    if err != nil {
        return fmt.Errorf("in fileChecker: %w", err) // Wrapping the error
    }
    f.Close()
    return nil
}

func main() {
    err := fileChecker("not_here.txt")
    if err != nil {
        fmt.Println(err) // Prints the wrapped error
        if wrappedErr := errors.Unwrap(err); wrappedErr != nil {
            fmt.Println(wrappedErr) // Prints the original error
        }
    }
}
```

**Output:**

```
in fileChecker: open not_here.txt: no such file or directory
open not_here.txt: no such file or directory
```

In this example, the `fmt.Errorf` function wraps the error returned from `os.Open` using `%w`, appending the message `"in fileChecker"`. The `errors.Unwrap` function is used to retrieve the original error (i.e., the error returned by `os.Open`).

#### Unwrapping Errors

The `errors.Unwrap` function returns the original (wrapped) error. If there’s no wrapped error, it returns `nil`. However, directly calling `errors.Unwrap` is rare; instead, it’s more common to use `errors.Is` and `errors.As` to find specific wrapped errors, which will be discussed in the next section.

#### Custom Wrapping with Your Own Error Type

To wrap errors using a custom error type, the error type must implement the `Unwrap` method, which returns the underlying error. Here’s an updated version of the `StatusErr` custom error type:

```go
type StatusErr struct {
    Status  Status
    Message string
    Err     error // Holds the original error
}

func (se StatusErr) Error() string {
    return se.Message
}

func (se StatusErr) Unwrap() error {
    return se.Err
}
```

With this updated structure, `StatusErr` can now wrap another error while also holding a status code. Here’s an example of how it can be used:

```go
func LoginAndGetData(uid, pwd, file string) ([]byte, error) {
    token, err := login(uid, pwd)
    if err != nil {
        return nil, StatusErr{
            Status:  InvalidLogin,
            Message: fmt.Sprintf("invalid credentials for user %s", uid),
            Err:     err, // Wrapping the error
        }
    }

    data, err := getData(token, file)
    if err != nil {
        return nil, StatusErr{
            Status:  NotFound,
            Message: fmt.Sprintf("file %s not found", file),
            Err:     err, // Wrapping the error
        }
    }

    return data, nil
}
```

In this example, the `StatusErr` type wraps the errors from the `login` and `getData` functions. The error message includes details about the failure, while the underlying error is preserved for further handling or debugging.

#### When to Wrap Errors

Not all errors need to be wrapped. If a library returns an error with unnecessary implementation details that aren't relevant to the user, it's fine to create a new error instead of wrapping the original one. The decision to wrap or replace an error depends on the context of the application.

#### Creating a New Error Without Wrapping

If you need to create a new error while including the message from an existing error without wrapping it, use `fmt.Errorf` with the `%v` verb instead of `%w`. Here’s an example:

```go
err := internalFunction()
if err != nil {
    return fmt.Errorf("internal failure: %v", err) // Creating a new error without wrapping
}
```

This method allows you to retain the original error's message in the new error's string, but the original error is not preserved for unwrapping. This is useful when the original error’s details are not needed beyond its message.

#### Conclusion

Wrapping errors is a powerful feature in Go that enables developers to add context while preserving the original error for later retrieval. However, it’s essential to understand when wrapping is necessary and when it’s better to return a new error. By following Go's conventions around error handling, code becomes more readable, maintainable, and easier to debug.

---

## Wrapping Multiple Errors

In Go, there are situations where a function might generate multiple errors, such as when validating the fields in a struct. Instead of returning just one error or an array of errors (`[]error`), Go allows multiple errors to be merged into a single error. This can be achieved using the `errors.Join` function or multiple `%w` verbs in `fmt.Errorf`.

#### Using `errors.Join` to Merge Multiple Errors

The `errors.Join` function merges multiple errors into a single error, making it easy to return and handle multiple error conditions. Here's an example using a `Person` struct where several fields might be invalid:

```go
type Person struct {
    FirstName string
    LastName  string
    Age       int
}

func ValidatePerson(p Person) error {
    var errs []error

    if len(p.FirstName) == 0 {
        errs = append(errs, errors.New("field FirstName cannot be empty"))
    }

    if len(p.LastName) == 0 {
        errs = append(errs, errors.New("field LastName cannot be empty"))
    }

    if p.Age < 0 {
        errs = append(errs, errors.New("field Age cannot be negative"))
    }

    if len(errs) > 0 {
        return errors.Join(errs...) // Merge errors into a single error
    }

    return nil
}
```

In this example, if one or more fields are invalid, `errors.Join` will combine the errors into a single error that can be returned. This allows the caller to handle all errors in one place.

#### Merging Errors Using Multiple `%w` Verbs

Another way to merge multiple errors is by using multiple `%w` verbs in `fmt.Errorf`:

```go
err1 := errors.New("first error")
err2 := errors.New("second error")
err3 := errors.New("third error")

err := fmt.Errorf("first: %w, second: %w, third: %w", err1, err2, err3)
```

In this case, each error is wrapped, and all the errors are combined into a single formatted string. However, this approach doesn't support the unwrapping of each individual error.

#### Custom Error Type Supporting Multiple Wrapped Errors

You can implement a custom error type that holds multiple wrapped errors by defining an `Unwrap` method that returns a slice of errors (`[]error`). This is useful if you need to handle multiple wrapped errors explicitly.

Here’s how to define such a custom error type:

```go
type MyError struct {
    Code   int
    Errors []error
}

func (m MyError) Error() string {
    return errors.Join(m.Errors...).Error()
}

func (m MyError) Unwrap() []error {
    return m.Errors
}
```

With this implementation, `MyError` can wrap multiple errors. The `Unwrap` method returns the slice of errors, allowing each error to be processed individually.

#### Handling Multiple Wrapped Errors

When dealing with errors that may wrap zero, one, or multiple errors, use a type switch to determine how to handle the wrapped errors. Here's an example:

```go
var err error
err = funcThatReturnsAnError()

switch err := err.(type) {
case interface{ Unwrap() error }: // Handle a single wrapped error
    innerErr := err.Unwrap()
    // process innerErr
case interface{ Unwrap() []error }: // Handle multiple wrapped errors
    innerErrs := err.Unwrap()
    for _, innerErr := range innerErrs {
        // process each innerErr
    }
default: // Handle no wrapped error
    // process err directly
}
```

Since the Go standard library does not define specific interfaces for errors with either `Unwrap` variant (i.e., single or multiple errors), this code uses anonymous interfaces in a type switch to match the methods dynamically.

#### Important Considerations

- **Use `errors.Join`**: The `errors.Join` function is the most straightforward way to handle multiple errors when they need to be returned together.
- **Type Switching**: If you need to work with custom error types that wrap multiple errors, use a type switch to differentiate between single and multiple wrapped errors.
- **Avoid Directly Calling `errors.Unwrap`**: Instead, rely on higher-level functions like `errors.Is` and `errors.As` to inspect error trees when possible.

By following these approaches, error handling in Go becomes more flexible and structured, especially when dealing with multiple error conditions.

---

## Is and As

### `errors.Is` and `errors.As` in Go

When wrapping errors, checking for specific errors becomes more complex because direct comparison using `==` no longer works. Go provides two functions, `errors.Is` and `errors.As`, to help handle these cases.

#### 1. `errors.Is`

The `errors.Is` function checks whether an error (or any error it wraps) matches a specific sentinel error, even if it's deep within the error tree. It allows you to check if an error contains a specific error, regardless of how many times it's been wrapped.

**Example using `errors.Is`:**

```go
func fileChecker(name string) error {
    f, err := os.Open(name)
    if err != nil {
        return fmt.Errorf("in fileChecker: %w", err) // Wrapping error
    }
    f.Close()
    return nil
}

func main() {
    err := fileChecker("not_here.txt")
    if err != nil {
        if errors.Is(err, os.ErrNotExist) { // Checking for the sentinel error
            fmt.Println("That file doesn't exist")
        }
    }
}
```

**Output:**

```
That file doesn't exist
```

In this example, even though the error from `os.Open` is wrapped inside another error, `errors.Is` successfully detects the `os.ErrNotExist` error.

#### Customizing `errors.Is`

If your custom error type needs special comparison logic (e.g., comparing complex fields), you can implement the `Is` method for that error type. Here's an example:

```go
type ResourceErr struct {
    Resource string
    Code     int
}

func (re ResourceErr) Error() string {
    return fmt.Sprintf("%s: %d", re.Resource, re.Code)
}

// Custom Is method to match on partial fields
func (re ResourceErr) Is(target error) bool {
    if other, ok := target.(ResourceErr); ok {
        ignoreResource := other.Resource == ""
        ignoreCode := other.Code == 0
        matchResource := other.Resource == re.Resource
        matchCode := other.Code == re.Code
        return (matchResource && matchCode) || (matchResource && ignoreCode) || (ignoreResource && matchCode)
    }
    return false
}
```

This custom `Is` method allows you to check for errors that match partially, such as any errors related to a database, no matter the code:

```go
if errors.Is(err, ResourceErr{Resource: "Database"}) {
    fmt.Println("The database is broken:", err)
}
```

#### 2. `errors.As`

While `errors.Is` checks for specific instances of errors, `errors.As` is used to check if an error (or any error it wraps) matches a specific **type**. This is helpful when you need to extract custom error details.

**Example using `errors.As`:**

```go
type MyErr struct {
    Codes []int
}

func (me MyErr) Error() string {
    return fmt.Sprintf("codes: %v", me.Codes)
}

func main() {
    err := AFunctionThatReturnsAnError()
    var myErr MyErr
    if errors.As(err, &myErr) { // Check if err is of type MyErr
        fmt.Println(myErr.Codes)
    }
}
```

In this example, `errors.As` checks if `err` is of type `MyErr`. If it is, it assigns the error to the `myErr` variable and allows you to access its fields.

#### Checking for Interface Types with `errors.As`

You can also use `errors.As` to check if an error matches an interface type. This lets you find errors that implement certain methods.

```go
err := AFunctionThatReturnsAnError()
var coder interface {
    CodeVals() []int
}
if errors.As(err, &coder) { // Check if err implements the interface
    fmt.Println(coder.CodeVals())
}
```

This example checks if `err` implements a method `CodeVals` that returns a slice of integers.

#### Summary

- **Use `errors.Is`** to check for a specific instance of an error or a known sentinel error.
- **Use `errors.As`** when you need to check for a specific type of error or extract information from custom errors.

These functions make it easier to handle errors that may be wrapped deep within other errors, improving error handling and debugging in Go.

---

## Wrapping Errors with defer

### Wrapping Errors with `defer`

When wrapping multiple errors with the same message in a function, the code can become repetitive. Instead of wrapping each error individually with `fmt.Errorf`, Go allows you to simplify the process using a `defer` statement.

#### Original Approach (Repetitive Error Wrapping)

In this example, the error is wrapped at each stage of the process:

```go
func DoSomeThings(val1 int, val2 string) (string, error) {
    val3, err := doThing1(val1)
    if err != nil {
        return "", fmt.Errorf("in DoSomeThings: %w", err)
    }

    val4, err := doThing2(val2)
    if err != nil {
        return "", fmt.Errorf("in DoSomeThings: %w", err)
    }

    result, err := doThing3(val3, val4)
    if err != nil {
        return "", fmt.Errorf("in DoSomeThings: %w", err)
    }

    return result, nil
}
```

Here, every time an error occurs, it is wrapped with the same message: `"in DoSomeThings: %w"`. This leads to repetitive wrapping at every return point.

#### Simplified Approach Using `defer`

Using a `defer` statement, you can handle error wrapping in a single place at the end of the function, reducing the repetition:

```go
func DoSomeThings(val1 int, val2 string) (_ string, err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("in DoSomeThings: %w", err) // Wrap the error at the end
        }
    }()

    val3, err := doThing1(val1)
    if err != nil {
        return "", err // No need to wrap here
    }

    val4, err := doThing2(val2)
    if err != nil {
        return "", err // No need to wrap here
    }

    return doThing3(val3, val4)
}
```

#### Key Points:

1. **Named Return Values**: The function’s return values are named (`_ string` and `err error`). This is required so the `defer` statement can access `err` and modify it before returning.
2. **Defer Function**: The `defer` function checks if `err` is not `nil`. If there’s an error, it wraps the original error with a message.

3. **Simplified Error Handling**: Errors from `doThing1`, `doThing2`, and `doThing3` are returned directly without wrapping. The wrapping only happens once in the `defer` function, making the code cleaner.

#### When to Use This Pattern

This pattern works best when you want to wrap every error with the same message, such as indicating which function the error occurred in. If you need to add different messages for each error, you’ll still need to wrap each error individually with `fmt.Errorf`.

#### Example with Customized Messages

If you want more detailed messages at each error point, you can still use `fmt.Errorf` for each error, but the `defer` pattern won't add much value:

```go
func DoSomeThings(val1 int, val2 string) (_ string, err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("in DoSomeThings: %w", err)
        }
    }()

    val3, err := doThing1(val1)
    if err != nil {
        return "", fmt.Errorf("doThing1 failed: %w", err)
    }

    val4, err := doThing2(val2)
    if err != nil {
        return "", fmt.Errorf("doThing2 failed: %w", err)
    }

    return doThing3(val3, val4)
}
```

In this case, each error can have a specific message, and the `defer` function can still add a general message if needed.

### Summary:

- The `defer` approach reduces repetitive error wrapping by wrapping the error only once at the end of the function.
- This method is most useful when every error needs the same message.
- For more detailed error messages, it's still necessary to wrap each error individually with `fmt.Errorf`.

---

## panic and recover

### Understanding Panics and Recovers in Go

#### What is a Panic?

In Go, a **panic** is similar to an error in languages like Java or Python, but it occurs when the Go runtime encounters a situation it cannot handle. This is typically due to programming mistakes, such as:

- Attempting to read beyond the end of a slice.
- Passing invalid values to built-in functions (e.g., passing a negative size to `make`).

When a panic happens:

- The function where the panic occurs stops immediately.
- Any deferred functions in that function are executed in reverse order.
- The panic propagates up the call stack, running deferred functions along the way.
- If the panic reaches the `main` function without being recovered, the program exits with a stack trace.

#### Example of a Panic

Here’s an example that triggers a panic using the `panic` function:

```go
func doPanic(msg string) {
    panic(msg) // Triggers a panic
}

func main() {
    doPanic("This is a panic")
}
```

**Output:**

```
panic: This is a panic
goroutine 1 [running]:
main.doPanic(...)
    /tmp/sandbox567884271/prog.go:6
main.main()
    /tmp/sandbox567884271/prog.go:10 +0x5f
```

When this code runs, the `panic` function halts execution, and the message "This is a panic" is printed, followed by the stack trace.

#### How to Recover from a Panic

You can use `recover` to gracefully handle a panic. The `recover` function must be called from within a `defer` statement. When a panic happens, `recover` captures the panic value and stops the program from crashing.

Here’s how `recover` is used:

```go
func div60(i int) {
    defer func() {
        if v := recover(); v != nil { // Recover from panic
            fmt.Println("Recovered from panic:", v)
        }
    }()
    fmt.Println(60 / i) // Panic if i is 0 (division by zero)
}

func main() {
    for _, val := range []int{1, 2, 0, 6} {
        div60(val) // 0 will cause a panic
    }
}
```

**Output:**

```
60
30
Recovered from panic: runtime error: integer divide by zero
10
```

Here’s what happens:

- The program attempts to divide 60 by values in the slice.
- When it tries to divide by 0, it triggers a **panic**.
- The deferred function runs, and `recover` captures the panic message `"runtime error: integer divide by zero"`.
- The program doesn’t crash; instead, it prints the error and continues.

#### Key Points About `panic` and `recover`:

1. **Panic Propagation**:
   When a panic occurs, Go executes all deferred functions in the current function and continues running deferred functions in the calling functions until it reaches `main`, where the program exits.

2. **Using `recover`**:

   - `recover` is used to stop the panic from propagating and crashing the program.
   - It must be used within a `defer` statement because only deferred functions are executed during a panic.

3. **When to Use `panic`**:

   - **Panics should be reserved for fatal situations** where the program cannot continue, such as serious programming errors.
   - It's usually better to handle errors with explicit checks (like returning an error) rather than relying on panics.
   - Use `recover` to log panics or handle graceful shutdowns, but avoid continuing execution after a panic unless you are confident the program can recover safely.

4. **Avoid Panic in Libraries**:
   If you're building a library for others to use, **never let panics escape your library**. Convert panics into errors using `recover`, and return an error instead of crashing the user's program.

5. **Panics in Goroutines**:
   - If a panic occurs in a goroutine, it behaves the same as in the main goroutine—deferred functions are run in the current function, but the panic will not propagate outside the goroutine.
   - If the panic is not recovered within the goroutine, the program will exit.

#### Example of Panic Recovery in a Library:

Here’s an example where a library function recovers from a panic and converts it into an error:

```go
func safeDivide(a, b int) (int, error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic occurred: %v", r) // Convert panic into error
        }
    }()

    if b == 0 {
        panic("division by zero") // Trigger panic
    }

    return a / b, nil
}

func main() {
    result, err := safeDivide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }
}
```

**Output:**

```
Error: panic occurred: division by zero
```

The `safeDivide` function recovers from the panic caused by dividing by zero and returns an error instead of crashing the program.

### Summary:

- **Panic** is used for fatal errors, but it’s not meant for regular error handling.
- **Recover** stops a panic and allows the program to continue running or handle errors more gracefully.
- **Use panics carefully**, and only when necessary—Go encourages explicit error handling over panics.

### Simplify

Let’s simplify the concepts of **panic** and **recover** in Go.

### What is a Panic?

A **panic** in Go happens when something goes seriously wrong, like trying to divide by zero or reading outside the bounds of a list. It's Go's way of saying, "I don't know how to handle this, something is broken!" When a panic occurs:

- The current function stops running.
- Go runs any deferred functions (things you told Go to run before the function exits).
- Then Go continues this process up the call stack until it reaches `main`.
- If no one "recovers" from the panic, the program crashes, and Go prints an error message and a stack trace.

### Example of Panic

Here’s a simple function that causes a panic:

```go
func doPanic() {
    panic("Something went wrong!") // Triggers a panic
}

func main() {
    doPanic() // Calls the function, which panics
}
```

When this runs, the program immediately stops with a message: "panic: Something went wrong!" and shows a stack trace to help debug where the panic happened.

### What is Recover?

`recover` is Go’s way of catching a panic so that the program doesn’t crash. It must be used inside a `defer` statement, which means it will run **after** the main function logic but before the function finishes.

### Example of Recover

```go
func divide(a, b int) {
    defer func() {
        if r := recover(); r != nil { // Catch the panic
            fmt.Println("Recovered from panic:", r)
        }
    }()

    result := a / b
    fmt.Println("Result:", result)
}

func main() {
    divide(10, 2)  // This works
    divide(10, 0)  // This causes a panic (division by zero)
}
```

**What happens:**

- When `divide(10, 0)` tries to divide by zero, Go triggers a **panic**.
- But because of the `defer` and `recover`, the panic is caught, and the program doesn’t crash. Instead, you see the message: "Recovered from panic: runtime error: integer divide by zero."

### Key Points:

1. **Panic**: Use this when something goes really wrong, and you want the program to stop immediately.
2. **Recover**: This lets you catch the panic and handle it gracefully, like logging the error instead of letting the program crash.
3. **Use `defer` with `recover`**: You must use `recover` inside a deferred function to catch the panic before it crashes the program.

### When to Use Panic and Recover?

- **Panic** should only be used for unexpected or fatal issues. It's not meant for regular errors.
- **Recover** can be used to gracefully handle panics, such as logging the error or cleaning up before the program exits.

### Summary:

- **Panic** is for serious problems where Go doesn't know how to continue.
- **Recover** is used to catch the panic and handle it so the program doesn’t crash.
- Use panics **carefully** and only for critical issues. Most of the time, you should handle errors in a more explicit way (like returning an error).

---

## Getting a Stack Trace from an Error

### Why Use a Stack Trace?

A **stack trace** shows all the function calls that happened before an error. This is useful for debugging because it helps you see where things went wrong. In Go, **panics** automatically print a stack trace, but regular errors don't.

### Adding Stack Traces to Errors

By default, Go doesn't include stack traces in errors. If you want stack traces, you have to add them yourself or use third-party libraries that do this for you.

### Using Third-Party Libraries

Some libraries, like **CockroachDB’s error library**, automatically add stack traces to errors. This saves you from doing it manually.

### How to See the Stack Trace

If you're using a library that includes stack traces, you can print the error and the stack trace like this:

```go
package main

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

// A function that returns an error
func doSomething() error {
	return errors.New("something went wrong")
}

func main() {
	// Call the function and capture the error
	err := doSomething()

	if err != nil {
		// Print the error with the stack trace
		fmt.Printf("%+v\n", err)
	}
}  // Prints the error with a stack trace
```

### Hiding Full Paths in Stack Traces

Stack traces often show the **full file path** on your computer. If you don’t want to expose that (e.g., for security reasons), you can use the **`-trimpath`** option when building your program. This removes the full path and just shows the package name.

```bash
go build -trimpath -o myprogram
```

### Key Points:

- Go doesn't show stack traces for regular errors by default, only for panics.
- You can use libraries like **CockroachDB’s error library** to automatically add stack traces to errors.
- Use `fmt.Printf("%+v", err)` to print the error with the stack trace.
- Use the `-trimpath` flag when building to hide the full file paths from the stack trace.

This way, you get better error information and control what details are shown in the stack trace!

---

## Exercises

Here’s how you can approach each of the exercises based on the code improvements for error handling.

### 1. Create a Sentinel Error for Invalid ID

A **sentinel error** is a constant error that represents a specific condition. In this exercise, we’ll create a sentinel error for an **invalid ID** and use `errors.Is` to check for it in `main`.

#### Steps:

- Define a sentinel error `ErrInvalidID`.
- Modify the `getEmployeeByID` function to return `ErrInvalidID` when the ID is invalid.
- In `main`, use `errors.Is` to check for this error and print a message.

**Code:**

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel error for invalid ID
var ErrInvalidID = errors.New("invalid employee ID")

// A function that returns an error if the ID is invalid
func getEmployeeByID(id int) (string, error) {
    if id <= 0 {
        return "", ErrInvalidID // Return sentinel error for invalid ID
    }
    return "Employee Name", nil
}

func main() {
    id := -1 // Invalid ID

    // Call the function and handle the error
    _, err := getEmployeeByID(id)
    if err != nil {
        if errors.Is(err, ErrInvalidID) { // Check for the sentinel error
            fmt.Println("Error: Invalid employee ID provided")
        }
    }
}
```

**Output:**

```
Error: Invalid employee ID provided
```

### 2. Define a Custom Error Type for an Empty Field

Here, we’ll define a **custom error type** for an empty field in an `Employee` struct. The error will include the field name that is empty, and we’ll use `errors.As` in `main` to check for this custom error and print the field name.

#### Steps:

- Define a custom error type `EmptyFieldError`.
- Modify the function to return this error if any employee field is empty.
- In `main`, use `errors.As` to check for this error and print the field name.

**Code:**

```go
package main

import (
    "errors"
    "fmt"
)

// Custom error type for empty field
type EmptyFieldError struct {
    Field string
}

func (e EmptyFieldError) Error() string {
    return fmt.Sprintf("empty field: %s", e.Field)
}

// Function that returns an error for empty fields
func checkEmployeeFields(name, position string) error {
    if name == "" {
        return EmptyFieldError{Field: "Name"}
    }
    if position == "" {
        return EmptyFieldError{Field: "Position"}
    }
    return nil
}

func main() {
    name := ""
    position := "Developer"

    // Call the function and handle the error
    err := checkEmployeeFields(name, position)
    if err != nil {
        var emptyFieldErr EmptyFieldError
        if errors.As(err, &emptyFieldErr) { // Check for the custom error
            fmt.Printf("Error: %s field is empty\n", emptyFieldErr.Field)
        }
    }
}
```

**Output:**

```
Error: Name field is empty
```

### 3. Return a Single Error Containing All Errors

In this exercise, instead of returning just the first error found, we will accumulate all errors and return them as a single error using `errors.Join`.

#### Steps:

- Modify the function to collect all errors.
- Use `errors.Join` to combine the errors and return a single error.

**Code:**

```go
package main

import (
    "errors"
    "fmt"
)

// Custom error type for empty fields, re-used from previous example
type EmptyFieldError struct {
    Field string
}

func (e EmptyFieldError) Error() string {
    return fmt.Sprintf("empty field: %s", e.Field)
}

// Function that returns a single error containing all discovered errors
func checkAllEmployeeFields(name, position string) error {
    var errs []error

    if name == "" {
        errs = append(errs, EmptyFieldError{Field: "Name"})
    }
    if position == "" {
        errs = append(errs, EmptyFieldError{Field: "Position"})
    }

    if len(errs) > 0 {
        return errors.Join(errs...) // Combine errors into a single error
    }

    return nil
}

func main() {
    name := ""
    position := ""

    // Call the function and handle the error
    err := checkAllEmployeeFields(name, position)
    if err != nil {
        fmt.Printf("Errors occurred:\n%v\n", err)
    }
}
```

**Output:**

```
Errors occurred:
[empty field: Name empty field: Position]
```

### Summary:

- **Exercise 1**: A **sentinel error** (`ErrInvalidID`) was created for invalid employee IDs, and `errors.Is` was used to check for this error.
- **Exercise 2**: A **custom error type** (`EmptyFieldError`) was defined to represent empty fields, and `errors.As` was used to check and handle this error.
- **Exercise 3**: Multiple errors were collected and returned as a single error using `errors.Join`. This allows you to handle all issues in one go instead of stopping at the first error.
