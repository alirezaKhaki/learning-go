# CHAPTER 15 Writing Tests

## **Understanding the Basics of Testing in Go**

Go's testing framework consists of two key components:

1. **Standard Library `testing` Package:** Provides types and functions for writing tests.
2. **`go test` Tool:** Executes tests and generates reports.

**Test Structure and Conventions:**

- **Test File Placement:** Tests are typically placed in files with a `_test.go` suffix within the same package as the code being tested.
- **Test Function Naming:** Start with `Test` and include a descriptive name for the test (e.g., `TestAddNumbers`).
- **`testing.T` Parameter:** Test functions take a `*testing.T` parameter for reporting errors and other information.
- **Test Logic:** Use standard Go code to call the tested function and validate results.
- **Error Reporting:** Use `t.Error`, `t.Errorf`, or other methods to report test failures.

**Running Tests:**

- `go test`: Runs tests in the current directory and subdirectories.
- `go test <package_name>`: Runs tests for a specific package.
- `go test -v`: Enables verbose output for more detailed test results.

**Example:**

```go
// adder.go
package adder

func addNumbers(x, y int) int {
    return x + y
}

// adder_test.go
package adder

import "testing"

func TestAddNumbers(t *testing.T) {
    result := addNumbers(2, 3)
    if result != 5 {
        t.Errorf("incorrect result: expected 5, got %d", result)
    }
}
```

**Key Points:**

- Tests are essential for ensuring code quality and correctness.
- Go's testing framework provides a straightforward and efficient way to write and run tests.
- Adhere to naming conventions and best practices for maintainable tests.
- Utilize `go test` command and its flags for effective test execution.

By following these guidelines and leveraging Go's testing tools, you can write comprehensive and reliable tests for your Go applications.

---

## **Reporting Test Failures in Go**

**Key Methods:**

- **`t.Error(format string, args ...any)`:** Marks the test as failed and prints a formatted error message.
- **`t.Errorf(format string, args ...any)`:** Same as `t.Error`, but uses printf-style formatting.
- **`t.Fatal(format string, args ...any)`:** Marks the test as failed and terminates the test function immediately.
- **`t.Fatalf(format string, args ...any)`:** Same as `t.Fatal`, but uses printf-style formatting.

**Choosing the Right Method:**

- **`Error`/`Errorf`:** Use when multiple checks are performed and individual failures should be reported.
- **`Fatal`/`Fatalf`:** Use when a failure indicates a fundamental problem that prevents further testing.

**Example:**

```go
func TestMyFunction(t *testing.T) {
    // Check multiple conditions
    if result != expected {
        t.Errorf("Incorrect result: expected %d, got %d", expected, result)
    }
    if err != nil {
        t.Error("Unexpected error:", err)
    }

    // Check a critical condition
    if criticalError != nil {
        t.Fatalf("Critical error: %v", criticalError)
    }
}
```

**Best Practices:**

- Provide informative error messages to help diagnose issues.
- Use `Fatal` or `Fatalf` for critical failures that prevent further testing.
- Consider using `t.Log` or `t.Logf` for non-fatal debugging messages.

By effectively using these methods, you can write more informative and actionable tests in Go.

---

## Setting Up and Tearing Down in Tests

Sometimes, common state needs to be set up before tests run and cleaned up afterward. This can be handled using a `TestMain` function, which manages the setup and teardown of resources for the tests.

### Example of `TestMain`

```go
var testTime time.Time

func TestMain(m *testing.M) {
    fmt.Println("Set up stuff for tests here")
    testTime = time.Now()
    exitVal := m.Run() // Run all tests
    fmt.Println("Clean up stuff after tests here")
    os.Exit(exitVal)    // Exit with test result code
}

func TestFirst(t *testing.T) {
    fmt.Println("TestFirst uses stuff set up in TestMain", testTime)
}

func TestSecond(t *testing.T) {
    fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}
```

- **TestMain** is invoked once per package and sets up any necessary state.
- After setup, `m.Run()` is called to run all tests, and `os.Exit(exitVal)` ensures the process exits with the correct status.
- Both `TestFirst` and `TestSecond` can access the initialized package-level variable `testTime`.

### Output Example

```bash
$ go test
Set up stuff for tests here
TestFirst uses stuff set up in TestMain 2020-09-01 21:42:36.231508
TestSecond also uses stuff set up in TestMain 2020-09-01 21:42:36.231508
PASS
Clean up stuff after tests here
```

### Use Cases for `TestMain`

1. Setting up data in external repositories (e.g., databases).
2. Initializing package-level variables (although generally discouraged).

### Cleanup for Individual Tests

For cleaning up resources after each test, Go provides the `Cleanup` method on `*testing.T`. This is useful for managing temporary files or other resources created during tests.

### Example of `t.Cleanup`

```go
func createFile(t *testing.T) (_ string, err error) {
    f, err := os.Create("tempFile")
    if err != nil {
        return "", err
    }
    defer func() { err = errors.Join(err, f.Close()) }()

    t.Cleanup(func() {
        os.Remove(f.Name()) // Remove temp file after test
    })
    return f.Name(), nil
}

func TestFileProcessing(t *testing.T) {
    fName, err := createFile(t)
    if err != nil {
        t.Fatal(err)
    }
    // Test logic, no need to worry about cleanup
}
```

### Using `t.TempDir` for Temporary Files

You can avoid manual cleanup by using `t.TempDir()`, which automatically creates a temporary directory for the test and removes it after completion.

```go
func createFile(tempDir string) (_ string, err error) {
    f, err := os.CreateTemp(tempDir, "tempFile")
    if err != nil {
        return "", err
    }
    defer func() { err = errors.Join(err, f.Close()) }()

    return f.Name(), nil
}

func TestFileProcessing(t *testing.T) {
    tempDir := t.TempDir()    // Create a temp directory
    fName, err := createFile(tempDir)
    if err != nil {
        t.Fatal(err)
    }
    // Perform test actions
}
```

---

## Testing with Environment Variables in Go

When writing tests for code that relies on environment variables, Go provides a helpful method called `t.Setenv`. This method allows you to set environment variables within your test and ensures they are reset when the test finishes.

### Example of `t.Setenv`

```go
// Assume ProcessEnvVars is a function that reads environment variables
// and returns a configuration struct with an OutputFormat field.

func TestEnvVarProcess(t *testing.T) {
    t.Setenv("OUTPUT_FORMAT", "JSON") // Set environment variable for the test
    cfg := ProcessEnvVars() // Process the environment variables

    if cfg.OutputFormat != "JSON" {
        t.Error("OutputFormat not set correctly")
    }
    // OUTPUT_FORMAT is automatically reset after the test
}
```

### Key Points

- `t.Setenv` sets an environment variable only for the duration of the test.
- After the test exits, the environment variable is reset to its previous state.
- This ensures that tests don't interfere with each other or the system's environment settings.

### Best Practices

- **Minimize Direct Use of Environment Variables**: It’s a good idea to limit how much of your code is aware of environment variables. Instead, you should read the environment variables early in your program (like in `main`) and pass them into the rest of your application via configuration structs.
- **Use Configuration Libraries**: Libraries like **Viper** or **envconfig** can help manage configuration, including environment variables. You can also use **GoDotEnv** to handle `.env` files for development and CI environments.

---

## Storing Sample Test Data for Go Tests

In Go, you can store sample data for tests in a directory named **`testdata`**. This directory is reserved by Go for testing purposes, making it an ideal place to store files your tests need to read or process.

### How to Use `testdata`:

1. **Create the `testdata` Directory**: Inside your package directory, create a folder named `testdata`. This folder will hold your sample data files.
2. **Access Files with Relative Paths**: When accessing files from the `testdata` directory in your tests, use relative paths. Go automatically sets the current working directory to your package’s directory when running tests, so you can reliably reference files inside `testdata`.

### Example:

```go
func TestFileReading(t *testing.T) {
    data, err := os.ReadFile("testdata/sample.txt")
    if err != nil {
        t.Fatal(err)
    }
    // Process data from sample.txt in your test
    t.Log(string(data))
}
```

### Key Points:

- **Directory Name**: Always name the directory `testdata`—this is a convention Go reserves for test-related files.
- **Relative Paths**: Reference the files using relative paths to the `testdata` folder.
- **Current Working Directory**: When running `go test`, the working directory is set to the package, ensuring that the relative file paths work consistently.

By following this approach, you can keep your test data organized and make it easy for your tests to access the necessary files.

### **Directory schema**

Here's a sample directory schema for organizing Go tests using the `testdata` folder:

```
project-root/
│
├── mypackage/
│   ├── main.go
│   ├── main_test.go
│   ├── testdata/
│   │   ├── sample.txt
│   │   ├── input.json
│   │   └── config.yaml
│   └── go.mod
└── go.sum
```

### Directory Explanation:

- **`project-root/`**: This is the root directory of your project.
  - **`mypackage/`**: The package you're testing, which contains your Go source files and tests.
    - **`main.go`**: Your main source code file.
    - **`main_test.go`**: Your test file containing the test cases.
    - **`testdata/`**: A subdirectory reserved for storing test files.
      - **`sample.txt`**, **`input.json`**, **`config.yaml`**: Example test data files that your tests will use.
    - **`go.mod`**: The module file for your Go project.
    - **`go.sum`**: Dependency file generated by Go.

### How to Access Files in Tests:

In your test file (`main_test.go`), you can access the files in `testdata` using relative paths like this:

```go
func TestFileReading(t *testing.T) {
    data, err := os.ReadFile("testdata/sample.txt")
    if err != nil {
        t.Fatal(err)
    }
    t.Log(string(data))
}
```

This schema keeps your test data organized and easily accessible within your Go package.

---

## Caching Test Results

Go's test caching system works similarly to its package caching. If the source code or test data hasn't changed, Go caches the test results and skips re-running them. This can save time when running multiple tests.

### Key Points:

- **Test Caching**: Go automatically caches test results when tests pass, and no source or testdata files have changed.
- **Recompilation and Re-running**: If you change any file in the package (including files in the `testdata` directory), Go will recompile and rerun the tests.
- **Forcing Tests to Run**: You can bypass the cache and force tests to always run by using the `-count=1` flag.

### Example:

```bash
go test -count=1 ./...
```

This command forces Go to re-run all tests, even if the cache indicates that they haven't changed.

---

## Testing Your Public API

In Go, you can test both exported and unexported functions by placing your test files within the same package or using the `_test` package suffix to test only the public API.

### Key Concepts:

- **Same Package Testing**: This allows access to both exported and unexported functions.
- **Public API Testing (`_test` suffix)**: By using a different package name (e.g., `packagename_test`), you only test exported functions, forcing you to interact with your package as a consumer would.

### Example:

#### Production Code (`adder.go`):

```go
package pubadder

// AddNumbers adds two integers.
func AddNumbers(x, y int) int {
    return x + y
}
```

#### Test Code (`adder_public_test.go`):

```go
package pubadder_test

import (
    "github.com/yourmodule/pubadder"
    "testing"
)

func TestAddNumbers(t *testing.T) {
    result := pubadder.AddNumbers(2, 3)
    if result != 5 {
        t.Errorf("incorrect result: expected 5, got %d", result)
    }
}
```

### Key Points:

- **Test as Black Box**: Using `pubadder_test` as the package name ensures you only interact with exported features like a client using your package.
- **Imports**: Even though the tests and code are in the same directory, you need to import the package to test its public API.

This practice ensures you can test your package's functionality as it would be used externally.

---

## Using go-cmp to Compare Test Results

The `go-cmp` package from Google simplifies comparing complex data types in Go tests by providing a detailed description of differences between two instances. This can be particularly useful when comparing structs, maps, or slices, as it helps identify discrepancies without writing verbose code.

### Example with `go-cmp`

#### Defining a Struct and Factory Function

Here's a simple struct `Person` and a factory function `CreatePerson` that returns a `Person` with the current timestamp:

```go
type Person struct {
    Name      string
    Age       int
    DateAdded time.Time
}

func CreatePerson(name string, age int) Person {
    return Person{
        Name:      name,
        Age:       age,
        DateAdded: time.Now(),
    }
}
```

#### Testing with `go-cmp`

In the test file, import the `cmp` package and use its `Diff` function to compare expected and actual outputs:

```go
import (
    "testing"
    "github.com/google/go-cmp/cmp"
)

func TestCreatePerson(t *testing.T) {
    expected := Person{
        Name: "Dennis",
        Age:  37,
    }
    result := CreatePerson("Dennis", 37)

    if diff := cmp.Diff(expected, result); diff != "" {
        t.Error(diff)  // If there's a difference, print it
    }
}
```

In this case, `cmp.Diff` compares the two `Person` instances. If there is any difference, it returns a string describing what doesn’t match.

#### Output of a Failed Test

If the test fails, `go-cmp` will show detailed output. For instance, when comparing `DateAdded`:

```
--- FAIL: TestCreatePerson (0.00s)
    --- Person{
        Name: "Dennis",
        Age: 37,
    -   DateAdded: s"0001-01-01 00:00:00 +0000 UTC",
    +   DateAdded: s"2024-09-17 12:53:58.087229 +0000 UTC",
    }
FAIL
```

The `-` and `+` indicate the differences in `DateAdded`, which cannot be controlled because `time.Now()` generates the current time.

#### Ignoring Fields with a Custom Comparator

To ignore fields like `DateAdded`, you can create a custom comparator function:

```go
comparer := cmp.Comparer(func(x, y Person) bool {
    return x.Name == y.Name && x.Age == y.Age
})
```

This comparator checks only the `Name` and `Age` fields, ignoring `DateAdded`.

Modify the test to include the custom comparator:

```go
if diff := cmp.Diff(expected, result, comparer); diff != "" {
    t.Error(diff)
}
```

### Benefits of `go-cmp`

- **Detailed Differences**: `go-cmp` provides clear output showing exactly what fields differ.
- **Custom Comparators**: You can create custom comparison functions to handle specific cases, such as ignoring certain fields.
- **Readable Tests**: It reduces verbosity in test code, making comparisons cleaner and more maintainable.

For further customization and features, check the [go-cmp documentation](https://pkg.go.dev/github.com/google/go-cmp/cmp).

---

## Running Table Tests

Table-driven tests in Go are a powerful way to manage and organize multiple test cases for a function without repeating code. This approach allows you to define your test cases in a structured way and reduces redundancy.

### Example: Table-Driven Tests

Let's see how to implement table-driven tests with a function `DoMath`:

#### Function to Test

Here’s the `DoMath` function that performs various arithmetic operations:

```go
func DoMath(num1, num2 int, op string) (int, error) {
    switch op {
    case "+":
        return num1 + num2, nil
    case "-":
        return num1 - num2, nil
    case "*":
        return num1 * num2, nil
    case "/":
        if num2 == 0 {
            return 0, errors.New("division by zero")
        }
        return num1 / num2, nil
    default:
        return 0, fmt.Errorf("unknown operator %s", op)
    }
}
```

#### Table-Driven Test

Instead of writing repetitive test code, you use a table-driven approach to organize your test cases:

```go
func TestDoMath(t *testing.T) {
    data := []struct {
        name     string
        num1     int
        num2     int
        op       string
        expected int
        errMsg   string
    }{
        {"addition", 2, 2, "+", 4, ""},
        {"subtraction", 2, 2, "-", 0, ""},
        {"multiplication", 2, 2, "*", 4, ""},
        {"division", 2, 2, "/", 1, ""},
        {"bad_division", 2, 0, "/", 0, "division by zero"},
        {"unknown_op", 2, 2, "%", 0, "unknown operator %"},
    }

    for _, d := range data {
        t.Run(d.name, func(t *testing.T) {
            result, err := DoMath(d.num1, d.num2, d.op)
            if result != d.expected {
                t.Errorf("Expected %d, got %d", d.expected, result)
            }
            var errMsg string
            if err != nil {
                errMsg = err.Error()
            }
            if errMsg != d.errMsg {
                t.Errorf("Expected error message `%s`, got `%s`", d.errMsg, errMsg)
            }
        })
    }
}
```

### Explanation

1. **Define Test Cases**:

   - Create a slice of structs, where each struct represents a test case.
   - Fields include the name of the test case, input parameters, expected result, and expected error message.

2. **Run Tests**:

   - Loop over the test cases and use `t.Run` to run each test as a subtest.
   - For each subtest, call the function with parameters from the current struct.
   - Compare the actual output and error against the expected results.

3. **Error Comparison**:
   - If an error occurs, check its message against the expected error message.
   - For custom error types or named sentinel errors, use `errors.Is` or `errors.As` instead of comparing messages directly.

### Benefits of Table-Driven Tests

- **Reduced Redundancy**: Avoids repetitive code by defining test cases in a structured way.
- **Clarity**: Makes it clear what inputs and expected results are being tested.
- **Maintainability**: Easier to add or modify test cases without changing the core test logic.

This approach helps keep your test code concise and easy to manage, especially when dealing with functions that need extensive testing across various scenarios.

---

## Running Tests Concurrently

Running tests concurrently can speed up the testing process, but it's important to handle concurrency issues carefully. Here's a simplified guide to running tests concurrently in Go and addressing common pitfalls:

### Running Tests Concurrently

By default, Go runs tests sequentially. To run tests concurrently, use the `t.Parallel()` method. Here’s a basic example:

```go
func TestMyCode(t *testing.T) {
    t.Parallel()
    // Test code goes here
}
```

### Key Points

1. **Use `t.Parallel()`**: Call `t.Parallel()` as the first line inside your test function to mark it as a parallel test.

2. **Shared State**: Avoid parallel tests if they rely on shared mutable state. Tests should be independent to avoid conflicts.

3. **Environment Variables**: Parallel tests using `t.Setenv()` can cause issues. Avoid using `t.Setenv()` in parallel tests.

### Example: Parallel Table Tests

Consider a function to test and its corresponding table test:

```go
func toTest(input int) int {
    // Example function
    return input * 2
}

func TestParallelTable(t *testing.T) {
    data := []struct {
        name  string
        input int
        output int
    }{
        {"test1", 10, 20},
        {"test2", 30, 60},
        {"test3", 50, 100},
    }

    for _, d := range data {
        d := d // Shadow variable to avoid concurrency issues
        t.Run(d.name, func(t *testing.T) {
            t.Parallel() // Run this test in parallel
            result := toTest(d.input)
            if result != d.output {
                t.Errorf("Expected %d, got %d", d.output, result)
            }
        })
    }
}
```

### Explanation

1. **Shadowing**: By declaring `d := d` inside the loop, you create a local copy of `d` for each iteration. This prevents all parallel tests from seeing the same variable value.

2. **Using `t.Parallel()`**: This enables the test function to run in parallel with other tests that are also marked as parallel.

3. **Avoiding Shared State**: Ensure that each test is self-contained and does not rely on shared state to prevent conflicts and flaky tests.

### Troubleshooting

- **Loop Variable Issue**: In Go versions before 1.22, the loop variable `d` was shared across all parallel tests. Shadowing the variable, as shown above, solves this problem.

- **Environment Variables**: Be cautious about using `t.Setenv()` with parallel tests, as concurrent modifications can lead to unpredictable results.

By following these guidelines, you can effectively run tests concurrently in Go and ensure that they are both fast and reliable.

---

## Checking Your Code Coverage

### Checking Code Coverage in Go

Code coverage helps you identify which parts of your code are tested and which aren't. However, even with 100% code coverage, bugs can still exist. Here's a simplified explanation of how to check and interpret code coverage using Go.

### Step-by-Step Guide

1. **Running Tests with Code Coverage**
   To calculate code coverage while running your tests, use the `-cover` flag:

   ```bash
   go test -v -cover
   ```

   This command runs the tests and displays a summary of the code coverage (e.g., 87.5%).

2. **Saving Coverage Data to a File**
   You can save the coverage data to a file using the `-coverprofile` flag:

   ```bash
   go test -v -cover -coverprofile=c.out
   ```

   This creates a file (`c.out`) with the coverage data.

3. **Viewing Code Coverage in HTML**
   To visualize which lines are covered, use the `cover` tool to generate an HTML report:

   ```bash
   go tool cover -html=c.out
   ```

   This opens a web page in your browser that highlights the code:

   - **Green**: Covered by tests.
   - **Red**: Not covered by tests.
   - **Gray**: Not testable (e.g., comments or package declarations).

4. **Improving Coverage**
   After seeing uncovered lines, you can add additional test cases. For example, if a bad operator isn't covered, you can add this case:

   ```go
   {"bad_op", 2, 2, "?", 0, `unknown operator ?`},
   ```

   After rerunning the tests and viewing the updated coverage, you'll notice improved coverage.

5. **Identifying Bugs Despite 100% Coverage**
   Even with full coverage, bugs can still exist. For instance, a multiplication operation might be wrong:
   ```go
   {"another_mult", 2, 3, "*", 6, ""},
   ```
   If you get an unexpected result, like:
   ```
   Expected 6, got 5
   ```
   It means there's a logic bug (e.g., mistakenly adding instead of multiplying).

### Key Takeaways

- **Code coverage is a useful metric**, but it **does not guarantee bug-free code**. It's possible to have 100% coverage and still miss edge cases or logical errors.
- Make sure to review the logic and functionality in addition to ensuring high coverage.

By using coverage tools in Go, you can ensure that most of your code is tested while being mindful that thorough tests and logic validation are also essential.

---

## Fuzzing