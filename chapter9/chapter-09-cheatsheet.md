Here's a **Go Error Handling Cheat Sheet** that you can use for interviews. It covers common patterns and important functions for handling errors in Go.

---

# **Go Error Handling Cheat Sheet**

### **Basic Error Handling**

- **Returning Errors from Functions:**
  Functions in Go typically return both the result and an `error`. If the operation fails, the result is often zero-valued, and the error is non-`nil`.

  ```go
  func divide(a, b int) (int, error) {
      if b == 0 {
          return 0, errors.New("division by zero")
      }
      return a / b, nil
  }
  ```

- **Checking and Handling Errors:**
  Always check for errors immediately after calling a function.

  ```go
  result, err := divide(10, 0)
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println("Result:", result)
  }
  ```

### **Creating Custom Errors**

- **Creating Simple Errors:**
  Use `errors.New` to create basic errors.

  ```go
  import "errors"
  var ErrNotFound = errors.New("resource not found")
  ```

- **Using `fmt.Errorf` for Dynamic Errors:**
  Create formatted error messages with `fmt.Errorf`.

  ```go
  err := fmt.Errorf("failed to load file: %s", filename)
  ```

### **Sentinel Errors**

- **What Are Sentinel Errors?**
  Sentinel errors are predefined, reusable errors used to signal specific failure conditions. They're typically defined as package-level variables.

  ```go
  var ErrInvalidID = errors.New("invalid ID")
  ```

- **Checking Sentinel Errors:**
  Use `errors.Is` to check if an error matches a sentinel error.

  ```go
  if errors.Is(err, ErrInvalidID) {
      fmt.Println("Invalid ID provided")
  }
  ```

### **Wrapping and Unwrapping Errors**

- **Wrapping Errors with `fmt.Errorf`:**
  Use `%w` to wrap errors while adding more context.

  ```go
  err := fmt.Errorf("failed to open file: %w", err)
  ```

- **Unwrapping Errors:**
  Use `errors.Unwrap` to retrieve the underlying error.

  ```go
  if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
      fmt.Println("Original error:", unwrappedErr)
  }
  ```

### **Custom Error Types**

- **Defining Custom Error Types:**
  You can define your own error type by implementing the `Error()` method.

  ```go
  type EmptyFieldError struct {
      Field string
  }

  func (e EmptyFieldError) Error() string {
      return fmt.Sprintf("empty field: %s", e.Field)
  }
  ```

- **Using `errors.As` to Handle Custom Errors:**
  `errors.As` lets you check for specific error types.

  ```go
  var emptyFieldErr EmptyFieldError
  if errors.As(err, &emptyFieldErr) {
      fmt.Println("Empty field:", emptyFieldErr.Field)
  }
  ```

### **Combining Multiple Errors**

- **Using `errors.Join` to Combine Errors:**
  If multiple errors occur, you can combine them into one error using `errors.Join`.

  ```go
  var errs []error
  errs = append(errs, errors.New("first error"))
  errs = append(errs, errors.New("second error"))

  combinedErr := errors.Join(errs...)
  fmt.Println(combinedErr)  // Output: [first error second error]
  ```

### **Panic and Recover**

- **Using `panic` for Fatal Errors:**
  Panic stops the program immediately and should only be used for unrecoverable errors.

  ```go
  panic("Something went wrong!")
  ```

- **Recovering from a Panic:**
  Use `recover` within a `defer` statement to catch a panic and prevent the program from crashing.

  ```go
  func riskyOperation() {
      defer func() {
          if r := recover(); r != nil {
              fmt.Println("Recovered from panic:", r)
          }
      }()
      panic("Oops, something went wrong!")
  }
  ```

### **Best Practices for Error Handling**

1. **Always check for errors**: Never ignore an error returned by a function.

   ```go
   result, err := someFunc()
   if err != nil {
       return err
   }
   ```

2. **Wrap errors with context**: Use `fmt.Errorf` and `%w` to add more context to errors while preserving the original error.

   ```go
   return fmt.Errorf("failed to execute task: %w", err)
   ```

3. **Use `errors.Is` to check specific errors**: Use `errors.Is` to compare errors without relying on `==`.

   ```go
   if errors.Is(err, ErrNotFound) {
       // handle specific error
   }
   ```

4. **Use `errors.As` for custom error types**: Use `errors.As` to extract and handle specific custom errors.

   ```go
   var myErr CustomError
   if errors.As(err, &myErr) {
       // handle custom error
   }
   ```

5. **Avoid panic for normal errors**: Reserve `panic` for unrecoverable conditions and use normal error handling for everything else.

---

This cheat sheet covers most of the key error handling topics you might encounter in a Go interview. It highlights common patterns, usage of sentinel errors, custom errors, error wrapping, and `panic/recover` for handling critical issues.
