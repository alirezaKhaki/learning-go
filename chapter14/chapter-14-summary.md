# CHAPTER 14 The Context

## What Is the Context?

In Go, the **context** is an instance that implements the `context.Context` interface, which helps manage deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes. Here's a breakdown of how context is used:

### Why Use Context?

- **Context Passing**: Go follows the idiomatic convention of explicitly passing context as the first parameter to functions. This allows you to propagate deadlines, cancelations, and values between functions in your program.
- **Contexts in HTTP Servers**: When dealing with HTTP requests, you can use the context to carry deadlines or other request-specific information through middleware to handlers. This is especially useful when you have to manage request timeouts or cancellation signals.

### Creating Context

- **`context.Background()`**: Use this when you need an empty context as the starting point, typically in command-line applications or at the entry point of your program.

  ```go
  ctx := context.Background()
  ```

- **`context.TODO()`**: Use this during development as a placeholder when you aren't yet sure how context will be handled. This should not be used in production code.

  ```go
  ctx := context.TODO()
  ```

### Adding Context in HTTP Handlers

Since HTTP handlers can't be modified to directly include context, you use two important methods on `http.Request`:

1. **`req.Context()`**: Extracts the context from the request.
2. **`req.WithContext()`**: Returns a new `http.Request` that carries the updated context.

Here's how to use context in middleware:

```go
func Middleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        // Add or modify the context here
        req = req.WithContext(ctx)
        handler.ServeHTTP(rw, req)
    })
}
```

In the handler, the context is extracted, and the context-aware logic is invoked:

```go
func handler(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    result, err := logic(ctx, "example")
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    rw.Write([]byte(result))
}
```

### Using Context in HTTP Clients

When making outgoing HTTP requests, always use `NewRequestWithContext` to pass the context through the request:

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
if err != nil {
    // Handle error
}
resp, err := client.Do(req)
```

### Key Concepts of Context

- **Cancellation**: You can cancel operations if they take too long or aren't needed anymore. For example, an HTTP request could be canceled if the client disconnects.
- **Timeouts**: Contexts can include timeouts to prevent long-running operations from exhausting resources.
- **Values**: Context can store request-scoped values, but it’s not meant for passing global or long-lived data. You’ll see more about how to work with values in contexts next.

---

## Values in Context

In Go, context values can be useful for passing request-scoped information, such as user authentication data, GUIDs, or request metadata, through your program. While you should prefer passing data explicitly, there are cases where context is necessary, especially in middleware and request handlers.

### Using `context.WithValue`

The `context.WithValue` function is used to add values to a context. It creates a child context with the new value and is immutable, meaning the original context remains unchanged. The value can be accessed in deeper layers of the code.

Example of adding and retrieving values from a context:

```go
ctx := context.Background()
ctx = context.WithValue(ctx, "key", "value")

value, ok := ctx.Value("key").(string)
if !ok {
    fmt.Println("key not found")
} else {
    fmt.Println("value:", value)
}
```

### Choosing Keys for Context Values

To avoid conflicts and ensure type safety, don't use simple strings as keys. Instead, define custom types for keys:

```go
type userKey int

const key userKey = iota

ctx := context.WithValue(context.Background(), key, "user123")
value, ok := ctx.Value(key).(string)
if ok {
    fmt.Println("User:", value)
}
```

This method ensures no accidental collisions with other context values, especially when using third-party packages.

### Example: Managing Users in Middleware

Here's an example of how to add user information to a context in an HTTP middleware:

```go
func extractUser(req *http.Request) (string, error) {
    userCookie, err := req.Cookie("identity")
    if err != nil {
        return "", err
    }
    return userCookie.Value, nil
}

func Middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        user, err := extractUser(req)
        if err != nil {
            rw.WriteHeader(http.StatusUnauthorized)
            rw.Write([]byte("Unauthorized"))
            return
        }
        ctx := req.Context()
        ctx = context.WithValue(ctx, key, user)
        req = req.WithContext(ctx)
        h.ServeHTTP(rw, req)
    })
}
```

### Extracting Values from Context in Handlers

In your HTTP handler, you can access the value stored in the context:

```go
func handler(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    user, ok := ctx.Value(key).(string)
    if !ok {
        rw.WriteHeader(http.StatusInternalServerError)
        return
    }
    fmt.Println("User from context:", user)
    rw.Write([]byte("Hello " + user))
}
```

### Keeping Metadata in Context: GUID Tracking Example

In some cases, such as tracking requests with GUIDs, you might prefer to store metadata in the context rather than passing it explicitly. Here’s an example of a middleware that adds a GUID to the context and logs it:

```go
package tracker

import (
    "context"
    "fmt"
    "net/http"
    "github.com/google/uuid"
)

type guidKey struct{}

func contextWithGUID(ctx context.Context, guid string) context.Context {
    return context.WithValue(ctx, guidKey{}, guid)
}

func guidFromContext(ctx context.Context) (string, bool) {
    guid, ok := ctx.Value(guidKey{}).(string)
    return guid, ok
}

func Middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        if guid := req.Header.Get("X-GUID"); guid != "" {
            ctx = contextWithGUID(ctx, guid)
        } else {
            ctx = contextWithGUID(ctx, uuid.New().String())
        }
        req = req.WithContext(ctx)
        h.ServeHTTP(rw, req)
    })
}
```

### Using Context Values in Business Logic

You can pass the context with the GUID through your business logic without requiring the logic itself to be aware of the GUID:

```go
type Logger interface {
    Log(context.Context, string)
}

type LogicImpl struct {
    Logger Logger
}

func (l LogicImpl) Process(ctx context.Context, data string) (string, error) {
    l.Logger.Log(ctx, "Processing request")
    return "Processed: " + data, nil
}

type LoggerImpl struct{}

func (LoggerImpl) Log(ctx context.Context, message string) {
    if guid, ok := guidFromContext(ctx); ok {
        fmt.Printf("GUID: %s - %s\n", guid, message)
    } else {
        fmt.Println(message)
    }
}
```

### Complete Example: Middleware and Context Usage

Here’s how you can wire everything up:

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", handler)

    loggedMux := tracker.Middleware(mux)

    http.ListenAndServe(":8080", loggedMux)
}

func handler(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    user, ok := ctx.Value(key).(string)
    if ok {
        fmt.Fprintf(rw, "Hello %s", user)
    } else {
        fmt.Fprintf(rw, "Hello guest")
    }
}
```

This setup uses context to pass request-scoped data, such as a user identity or tracking GUID, through middleware, making the system flexible and maintainable.

---

## Cancellation

### Timeouts with Context

In addition to manual cancellation, Go’s `context` package provides a mechanism for timeouts. This allows you to automatically cancel a context if an operation takes too long, which is particularly useful when dealing with slow or unreliable external services.

To set up a timeout, use `context.WithTimeout`, which creates a context that is canceled automatically after a specified time. Just like `WithCancel`, it returns a new context and a cancellation function. Here’s a quick example of how to use it:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

select {
case <-time.After(1 * time.Second):
    fmt.Println("operation completed")
case <-ctx.Done():
    fmt.Println("operation timed out:", ctx.Err())
}
```

In this example:

- If the operation finishes within 1 second, you print "operation completed."
- If it takes longer than 2 seconds, the context cancels itself, and you print "operation timed out."

### Example: Using Timeout in HTTP Requests

Timeouts are particularly important in HTTP clients, where requests can hang or take too long to complete. Here’s how you can use `context.WithTimeout` with an HTTP client request:

```go
func fetchData(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }

    // process response body (omitted)
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := fetchData(ctx, "https://example.com")
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Request timed out:", err)
        } else {
            fmt.Println("Request failed:", err)
        }
    }
}
```

### `WithDeadline`

For more control over when the timeout occurs, you can use `context.WithDeadline`, which allows you to set a specific time for the context to expire rather than a duration. For example:

```go
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Your logic here, similar to the timeout example.
```

### Cancellation in Concurrent Operations

When you launch multiple goroutines, such as fetching data from several external services concurrently, you can use context cancellation to stop all operations if one of them fails or takes too long.

Here’s an example where two HTTP services are called concurrently. If one of the services returns an error or takes too long, all operations are canceled:

```go
func makeRequest(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    return http.DefaultClient.Do(req)
}

func fetchServiceData(ctx context.Context, urls []string) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    ch := make(chan error, len(urls))
    var wg sync.WaitGroup
    wg.Add(len(urls))

    for _, url := range urls {
        go func(url string) {
            defer wg.Done()
            resp, err := makeRequest(ctx, url)
            if err != nil || resp.StatusCode != http.StatusOK {
                ch <- fmt.Errorf("failed to fetch from %s: %w", url, err)
                cancel() // cancel all requests
                return
            }
            ch <- nil
        }(url)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    for err := range ch {
        if err != nil {
            return err
        }
    }

    return nil
}

func main() {
    urls := []string{
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/500", // This will cause cancellation
    }

    err := fetchServiceData(context.Background(), urls)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("All services responded successfully")
    }
}
```

### Conclusion

Using `context.WithCancel`, `context.WithTimeout`, and `context.WithDeadline`, you can build responsive and efficient Go programs that handle timeouts and cancellations in a clean and scalable way. By controlling concurrent goroutines and using `Done` channels, you ensure your application remains performant and doesn't leak resources.

---

## Contexts with Deadlines

Contexts with deadlines are an essential tool in Go for managing how long a server request should run. Instead of letting processes run indefinitely or until all resources are used up, you can limit their execution time using contexts with deadlines or timeouts.

### Key Concepts

- **Deadline**: A specific point in time when a context will automatically cancel.
- **Timeout**: A duration after which the context will cancel.

Both are used to enforce limits on how long a request or a function should run, ensuring that system resources are shared fairly among users.

### Creating Deadlines and Timeouts

You can use `context.WithTimeout` to specify how long a context should last before being canceled. Alternatively, you can use `context.WithDeadline` to set a specific time in the future when the context will cancel.

For example:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Simulate some work
select {
case <-time.After(2 * time.Second):
    fmt.Println("Work completed")
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

In this example, the `context.WithTimeout` will cancel the context after 5 seconds. If the simulated work finishes within that time, the program will print "Work completed". Otherwise, it will print the timeout error.

### Deadlines in Nested Contexts

When nesting contexts (such as parent and child contexts), the child's timeout is automatically constrained by the parent's deadline.

```go
parentCtx, parentCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer parentCancel()

childCtx, childCancel := context.WithTimeout(parentCtx, 10*time.Second)
defer childCancel()

select {
case <-childCtx.Done():
    fmt.Println("Child context canceled:", childCtx.Err())
}
```

Even though the child context is given a 10-second timeout, it will still cancel after 5 seconds due to the parent's timeout.

### Handling Timeouts and Cancellations

To understand why a context was canceled, you can use the `Err` method, which returns either `context.Canceled` (manual cancellation) or `context.DeadlineExceeded` (timeout).

For example, when you cancel a context manually:

```go
ctx, cancel := context.WithCancel(context.Background())
cancel()

fmt.Println("Context canceled:", ctx.Err()) // Output: Context canceled: context canceled
```

Or when a context times out:

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

select {
case <-time.After(2 * time.Second):
    fmt.Println("Work completed")
case <-ctx.Done():
    fmt.Println("Timeout error:", ctx.Err()) // Output: Timeout error: context deadline exceeded
}
```

### Practical Example: Making Timed HTTP Requests

Here’s an example of how to use context timeouts with HTTP requests to avoid long-running network calls:

```go
func fetchData(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }

    // Process response (omitted)
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := fetchData(ctx, "https://example.com")
    if err != nil {
        fmt.Println("Error fetching data:", err)
    } else {
        fmt.Println("Data fetched successfully")
    }
}
```

In this example, if the `fetchData` function takes longer than 3 seconds, the request is canceled, and the error is handled.

### Summary

Contexts with deadlines and timeouts are essential for managing long-running tasks, ensuring that your application remains responsive. Whether you're limiting the execution of a single task or coordinating multiple goroutines, Go’s context package provides powerful tools to prevent resource exhaustion and ensure timely cancellations.

This structured approach to timeouts helps ensure that your server scales effectively and can handle a large number of requests without overloading the system.

---

## Key Differences Between Context Tools:

The `context` package in Go provides tools to manage the lifecycle of requests, goroutines, or processes. It allows you to pass metadata, deadlines, timeouts, and cancellation signals through your program, especially in concurrent environments. Let's break down the key tools available in the `context` package and their differences:

### 1. **context.Background()**

- **Purpose**: Returns an empty, non-cancelable context.
- **Usage**: Used as a top-level, root context when no other context is available. It acts as a parent context for the entire flow.
- **When to use**: In entry points such as a `main` function, or when starting a new, independent task.
- **Cancelable?**: No.

```go
ctx := context.Background()
```

### 2. **context.TODO()**

- **Purpose**: Like `context.Background()`, it returns an empty context but indicates that you plan to add a real context later.
- **Usage**: Use it as a placeholder in situations where the context is required but hasn’t been fully designed yet.
- **When to use**: During development, if you're not sure how to handle contexts yet.
- **Cancelable?**: No.

```go
ctx := context.TODO()
```

### 3. **context.WithCancel(parent Context)**

- **Purpose**: Creates a context that can be explicitly canceled. Returns a new context and a `CancelFunc` function.
- **Usage**: When you want to manually signal that a process should be canceled, typically to propagate cancellation to child goroutines.
- **When to use**: Use this when you need to coordinate multiple concurrent operations and cancel them when one fails.
- **Cancelable?**: Yes, via the returned `CancelFunc`.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // Always ensure the cancel function is called
```

### 4. **context.WithDeadline(parent Context, deadline time.Time)**

- **Purpose**: Creates a context that automatically cancels at a specified deadline (`time.Time`). Useful for ensuring that operations don’t exceed a specific time limit.
- **Usage**: When a task must complete by a specific deadline (e.g., ensuring a task doesn't run past a certain time).
- **When to use**: If you need a precise point in time to stop all operations or trigger timeouts.
- **Cancelable?**: Yes, both explicitly and automatically once the deadline passes.

```go
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
defer cancel() // Cancels once the deadline is reached
```

### 5. **context.WithTimeout(parent Context, timeout time.Duration)**

- **Purpose**: Creates a context that automatically cancels after a specified duration.
- **Usage**: Similar to `WithDeadline`, but you specify a relative time (duration) instead of an absolute deadline.
- **When to use**: When you want to set a timeout for an operation (e.g., stopping a network request after 5 seconds).
- **Cancelable?**: Yes, both explicitly and automatically once the timeout expires.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // Cancels after 5 seconds
```

### 6. **context.WithValue(parent Context, key, value any)**

- **Purpose**: Creates a context that carries a key-value pair, allowing you to pass metadata (e.g., user IDs, request IDs) down the call chain.
- **Usage**: When you need to associate metadata or request-specific values with a context and access them in later stages of your code.
- **When to use**: When passing auxiliary data that might be needed across API boundaries, typically in middleware or request handling.
- **Cancelable?**: No. This is only for passing data, not managing cancellation or timeouts.

```go
ctx := context.WithValue(context.Background(), "userID", 1234)
userID := ctx.Value("userID").(int) // Extracting the value
```

### 7. **context.CancelFunc**

- **Purpose**: A function returned by `WithCancel`, `WithTimeout`, or `WithDeadline` that is used to cancel the context explicitly.
- **Usage**: Call the `CancelFunc` to cancel the associated context. This is important to prevent resource leaks.
- **When to use**: Always call this function when the context is no longer needed, either via `defer` or manual calls.

```go
ctx, cancel := context.WithCancel(context.Background())
cancel() // Cancels the context
```

### Key Differences Between Context Tools:

| **Context Function**   | **Purpose**                                | **Auto Cancel?** | **Manual Cancel?** | **Pass Metadata?** |
| ---------------------- | ------------------------------------------ | ---------------- | ------------------ | ------------------ |
| `context.Background`   | Root context, no cancellation              | No               | No                 | No                 |
| `context.TODO`         | Placeholder context                        | No               | No                 | No                 |
| `context.WithCancel`   | Context with manual cancellation           | No               | Yes                | No                 |
| `context.WithDeadline` | Context with a fixed cancellation deadline | Yes              | Yes                | No                 |
| `context.WithTimeout`  | Context with a relative timeout            | Yes              | Yes                | No                 |
| `context.WithValue`    | Context for storing key-value pairs        | No               | No                 | Yes                |

### Summary of Usage:

- **`context.Background`**: For top-level contexts where no cancellation is needed.
- **`context.TODO`**: Temporary context placeholder during development.
- **`context.WithCancel`**: For manual cancellation, useful in complex operations where you may need to stop several processes.
- **`context.WithDeadline`**: To set a hard deadline for a task, useful for requests that must be completed by a certain time.
- **`context.WithTimeout`**: To limit the duration of a task, ensuring it doesn't take too long to execute.
- **`context.WithValue`**: To pass key-value pairs through the context, typically for metadata like request IDs or user info.

By leveraging these tools, Go developers can better manage goroutines, timeouts, cancellation, and metadata, leading to more efficient and predictable concurrent programs.

---

## context.WithTimeout vs context.WithDeadline

`context.WithTimeout` and `context.WithDeadline` are two functions in Go's `context` package that create contexts that automatically cancel after a certain amount of time. While they achieve similar goals, they differ in how they specify when the cancellation should occur:

### 1. `context.WithTimeout`:

- **Purpose**: Creates a context that automatically cancels after a specified **duration**.
- **Usage**: You specify how long (in terms of time duration) the context should live before it gets canceled.
- **Parameters**:
  - `parent context`: The parent context that this new context will be derived from.
  - `timeout duration`: How long to wait before automatically canceling the context.
- **Behavior**: The returned context cancels after the given timeout period has passed.

```go
// Creates a context that cancels after 2 seconds
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

In this example, the context will automatically cancel 2 seconds after the call to `context.WithTimeout`.

### 2. `context.WithDeadline`:

- **Purpose**: Creates a context that automatically cancels at a **specific time** (i.e., a deadline).
- **Usage**: You specify an exact point in time (using `time.Time`) when the context should be canceled.
- **Parameters**:
  - `parent context`: The parent context that this new context will be derived from.
  - `deadline time`: The specific `time.Time` at which the context should be canceled.
- **Behavior**: The returned context cancels when the specified time (deadline) is reached.

```go
// Creates a context that cancels at a specific time (3 seconds from now)
deadline := time.Now().Add(3 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

In this example, the context will automatically cancel at the exact time that is 3 seconds from the current time.

### Key Differences:

1. **Timeout vs. Deadline**:

   - `WithTimeout`: You specify the duration the context should live, e.g., "2 seconds".
   - `WithDeadline`: You specify the exact time when the context should cancel, e.g., "at 12:30 PM on a specific day".

2. **Flexibility**:

   - `WithTimeout` is easier to use when you know how long you want a task to run.
   - `WithDeadline` is useful when you need the task to end at a specific point in time, regardless of when it starts.

3. **Return Values**:
   Both return a new context and a cancellation function (`CancelFunc`). You must call the cancellation function to free resources associated with the context, even if it times out or reaches the deadline.

4. **Use Cases**:
   - Use `context.WithTimeout` when you have a fixed amount of time to perform a task.
   - Use `context.WithDeadline` when you have a known cutoff time at which you must stop, regardless of how much time has passed.

### Example: Difference in Behavior

Let's assume you are dealing with two different use cases:

#### 1. Using `context.WithTimeout`:

```go
func runWithTimeout() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    select {
    case <-time.After(1 * time.Second): // Simulate some work that takes 1 second
        fmt.Println("Task completed before timeout.")
    case <-ctx.Done():
        fmt.Println("Timeout reached:", ctx.Err())
    }
}
```

Here, the task has 2 seconds to complete, but it finishes in 1 second. The context is still valid.

#### 2. Using `context.WithDeadline`:

```go
func runWithDeadline() {
    deadline := time.Now().Add(2 * time.Second) // Set a deadline 2 seconds in the future
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    select {
    case <-time.After(3 * time.Second): // Simulate some work that takes 3 seconds
        fmt.Println("Task completed.")
    case <-ctx.Done():
        fmt.Println("Deadline reached:", ctx.Err())
    }
}
```

In this case, the task takes 3 seconds, but the deadline is only 2 seconds away, so the task will not complete before the deadline, and the context will cancel.

### Summary:

- **`context.WithTimeout`** is used when you have a fixed time limit (duration) for a task.
- **`context.WithDeadline`** is used when you need the task to end at a specific point in time.

---

## Context Cancellation in Your Own Code

In Go, adding context cancellation to your own code can be useful when dealing with long-running tasks or operations that may need to stop when a cancellation is triggered. By periodically checking the context during your computation, you can gracefully handle cancellation and stop processing in response to external signals.

### Key Steps for Handling Context Cancellation in Your Code:

1. **Using `context.Cause(ctx)`**:

   - The `context.Cause(ctx)` function checks whether the context has been canceled or if a timeout has occurred.
   - It returns `nil` if the context is still active, otherwise, it returns an error that you can use to stop your function and possibly return a partial result.

2. **Inserting Periodic Checks**:
   - While performing long-running operations (like loops or heavy computations), insert periodic checks to see if the context has been canceled.
   - If a cancellation is detected, you can either stop immediately or return a partial result, depending on your use case.

### Example 1: Simple Long-Running Task with Context Cancellation

Here's an example of a long-running computation that checks for context cancellation periodically:

```go
func longRunningComputation(ctx context.Context, data string) (string, error) {
    result := ""
    for i := 0; i < 10000; i++ {
        // Perform some processing
        result += data

        // Periodically check if the context has been canceled
        if err := context.Cause(ctx); err != nil {
            // Stop processing and return error
            return result, err
        }
    }
    return result, nil
}
```

In this example, the computation will concatenate a string 10,000 times. If the context is canceled, it will stop and return the partial result accumulated up to that point.

### Example 2: Context Cancellation with a Computational Task (Calculating Pi)

Here's an example using the inefficient Leibniz algorithm to calculate π. The computation will be interrupted when the context is canceled:

```go
func calculatePi(ctx context.Context) (string, error) {
    var sum, d, two big.Float
    two.SetInt64(2)
    d.SetInt64(1)

    i := 0
    for {
        // Check for context cancellation
        if err := context.Cause(ctx); err != nil {
            fmt.Println("cancelled after", i, "iterations")
            return sum.Text('g', 100), err
        }

        // Leibniz algorithm to calculate Pi
        var diff big.Float
        diff.SetInt64(4)
        diff.Quo(&diff, &d)
        if i%2 == 0 {
            sum.Add(&sum, &diff)
        } else {
            sum.Sub(&sum, &diff)
        }
        d.Add(&d, two)
        i++
    }
}
```

In this example:

- The computation iteratively calculates π using the Leibniz formula.
- The `context.Cause(ctx)` check ensures that if the context is canceled (either manually or by a timeout), the loop will break, and the function will return the partial sum of π calculated so far.

### Best Practices:

- **Use `context.Cause` when cancellation needs to propagate an error**: This is especially useful if you want to distinguish between cancellations caused by timeouts and explicit cancellations or to return custom error messages.
- **Insert cancellation checks at logical intervals**: You don't need to check for cancellation at every step in your computation, just at regular intervals where it's feasible to stop processing without compromising performance.

- **Always call the cancellation function**: Whether or not the context was canceled, always ensure that the cancellation function (`CancelFunc`) is called, typically via `defer`.

By periodically checking the context in long-running tasks, you can build responsive and resource-efficient applications that can gracefully handle timeouts and cancellations.

---

## Exercises

Here are the solutions to the exercises:

### 1. Middleware to Create a Context with Timeout

The middleware-generating function creates a context with a timeout based on the provided number of milliseconds. This middleware will wrap the request context and ensure it has a time limit.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Middleware to add a timeout to the request's context
func TimeoutMiddleware(timeoutMillis int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a context with the timeout
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMillis)*time.Millisecond)
			defer cancel() // Ensure we call cancel to free resources
			// Create a new request with the context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// Example handler
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", TimeoutMiddleware(1000)(http.HandlerFunc(HelloHandler)))

	http.ListenAndServe(":8080", mux)
}
```

This middleware applies a timeout of the specified number of milliseconds to every HTTP request that passes through it.

---

### 2. Program to Sum Random Numbers with a Timeout

This program sums randomly generated numbers until either the number `1234` is generated or 2 seconds pass.

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sum := 0
	iterations := 0
	reason := "Timeout"
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				done <- struct{}{}
				return
			default:
				num := rand.Intn(100000000)
				sum += num
				iterations++
				if num == 1234 {
					reason = "Number 1234 found"
					cancel()
				}
			}
		}
	}()

	<-done

	fmt.Printf("Sum: %d\nIterations: %d\nReason: %s\n", sum, iterations, reason)
}
```

This program will generate random numbers until either `1234` is found or the 2-second timeout is reached. It will print out the sum, the number of iterations, and the reason for ending.

---

### 3. Logging with Context and Middleware

In this exercise, you create a logging middleware and extract the log level from the query parameters. You also add functions to store and retrieve the log level in the context.

#### Code:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
)

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
)

// Store the log level in the context
func ContextWithLogLevel(ctx context.Context, level Level) context.Context {
	return context.WithValue(ctx, "log_level", level)
}

// Retrieve the log level from the context
func LogLevelFromContext(ctx context.Context) (Level, bool) {
	level, ok := ctx.Value("log_level").(Level)
	return level, ok
}

// Middleware to extract log_level from the query parameters and store it in the context
func LogLevelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logLevel := r.URL.Query().Get("log_level")
		var level Level

		switch logLevel {
		case string(Debug):
			level = Debug
		case string(Info):
			level = Info
		default:
			// Invalid or missing log_level; no logging level is set
			next.ServeHTTP(w, r)
			return
		}

		// Add log level to the context
		ctx := ContextWithLogLevel(r.Context(), level)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Log function
func Log(ctx context.Context, level Level, message string) {
	var inLevel Level
	var ok bool

	// Get logging level from context
	inLevel, ok = LogLevelFromContext(ctx)
	if !ok {
		// No log level found, return without logging
		return
	}

	if level == Debug && inLevel == Debug {
		fmt.Println(message)
	}
	if level == Info && (inLevel == Debug || inLevel == Info) {
		fmt.Println(message)
	}
}

// Example handler
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	Log(r.Context(), Info, "Processing hello handler")
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", LogLevelMiddleware(http.HandlerFunc(HelloHandler)))

	http.ListenAndServe(":8080", mux)
}
```

#### Explanation:

- **ContextWithLogLevel**: Stores the log level (either `debug` or `info`) in the context.
- **LogLevelFromContext**: Retrieves the log level from the context.
- **LogLevelMiddleware**: Middleware that extracts the `log_level` from the query parameters and adds it to the context.
- **Log function**: Logs a message if the log level from the context matches or exceeds the requested log level (e.g., `debug` messages are only logged if the log level is `debug`).

#### Example Usage:

You can run the server and call it with:

- `/` → no logs
- `/?log_level=info` → logs Info level messages
- `/?log_level=debug` → logs both Debug and Info level messages.
