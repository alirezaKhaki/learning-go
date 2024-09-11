# CHAPTER 13 The Standard Library

## io and Friends

### Overview: `io` and Friends in Go

In Go, two of the most important interfaces for handling input and output are **`io.Reader`** and **`io.Writer`**. These interfaces are central to how Go handles reading and writing data in a variety of contexts, from files to network streams.

#### The `io.Reader` Interface

- **Definition**:
  ```go
  type Reader interface {
      Read(p []byte) (n int, err error)
  }
  ```
- **Functionality**: Reads data into a provided byte slice `p`. Returns the number of bytes read (`n`) and an error (`err`), typically `io.EOF` to signal the end of data.

#### The `io.Writer` Interface

- **Definition**:
  ```go
  type Writer interface {
      Write(p []byte) (n int, err error)
  }
  ```
- **Functionality**: Writes the byte slice `p` to the destination. Returns the number of bytes written (`n`) and an error (`err`).

### Key Concepts for `io.Reader` and `io.Writer`

1. **Efficient Memory Management**:

   - `Read` doesn’t return new slices every time, avoiding extra memory allocations.
   - The provided buffer is reused to read chunks of data, which improves performance.

2. **Detecting End of Data**:
   - You know that the reading process is finished when the `Read` method returns the **`io.EOF`** error. This signals the end of the data stream but isn’t an actual failure.

#### Example: Counting Letters from an `io.Reader`

Here’s an example function that reads data from any `io.Reader` and counts the letters:

```go
func countLetters(r io.Reader) (map[string]int, error) {
    buf := make([]byte, 2048)
    out := map[string]int{}
    for {
        n, err := r.Read(buf) // Read into buffer
        for _, b := range buf[:n] { // Process read bytes
            if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
                out[string(b)]++
            }
        }
        if err == io.EOF {
            return out, nil
        }
        if err != nil {
            return nil, err
        }
    }
}
```

#### Using `strings.NewReader`

You can use `strings.NewReader` to create a string-backed `io.Reader`, which makes it easy to test functions like `countLetters`:

```go
s := "The quick brown fox jumped over the lazy dog"
sr := strings.NewReader(s)
counts, err := countLetters(sr)
if err != nil {
    log.Fatal(err)
}
fmt.Println(counts)
```

### Wrapping Readers: Gzip Example

You can wrap different types of `io.Reader` implementations together to process various data formats. For example, to read a gzip-compressed file:

```go
func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
    r, err := os.Open(fileName)
    if err != nil {
        return nil, nil, err
    }
    gr, err := gzip.NewReader(r)
    if err != nil {
        return nil, nil, err
    }
    return gr, func() {
        gr.Close()
        r.Close()
    }, nil
}

r, closer, err := buildGZipReader("my_data.txt.gz")
if err != nil {
    log.Fatal(err)
}
defer closer()
counts, err := countLetters(r)
if err != nil {
    log.Fatal(err)
}
fmt.Println(counts)
```

### Utilities in `io` Package

- **`io.Copy`**: Copies data from an `io.Reader` to an `io.Writer`.

  ```go
  io.Copy(dst io.Writer, src io.Reader)
  ```

- **`io.MultiReader`**: Combines multiple `io.Reader` instances into one.
- **`io.LimitReader`**: Limits the number of bytes read from a reader.

- **`io.MultiWriter`**: Sends data to multiple `io.Writer` instances simultaneously.

### Other Interfaces: `io.Closer` and `io.Seeker`

- **`io.Closer`**: Used to close resources like files.

  ```go
  type Closer interface {
      Close() error
  }
  ```

- **`io.Seeker`**: Used for random access to a data stream, supporting operations like seeking within a file.
  ```go
  type Seeker interface {
      Seek(offset int64, whence int) (int64, error)
  }
  ```

### Helper Functions in `io`

- **`io.ReadAll`**: Reads all the data from an `io.Reader` into a byte slice.
- **`io.NopCloser`**: Adapts an `io.Reader` to an `io.ReadCloser` by adding a `Close` method that does nothing. Useful when a closer is needed but not available.

```go
r := strings.NewReader("Hello, World!")
rc := io.NopCloser(r)
defer rc.Close()
```

### File Operations: `os` Package

For larger file operations, you’ll often work with `*os.File`, which implements both `io.Reader` and `io.Writer`.

- **Reading a File**:

  ```go
  f, err := os.Open("file.txt")
  if err != nil {
      log.Fatal(err)
  }
  defer f.Close()
  ```

- **Writing to a File**:
  ```go
  f, err := os.Create("output.txt")
  if err != nil {
      log.Fatal(err)
  }
  defer f.Close()
  ```

### Conclusion

The **`io` package** provides simple, yet powerful abstractions for handling input and output in Go. By using the **`io.Reader`** and **`io.Writer`** interfaces, you can easily manage various data sources and sinks, and by leveraging utilities like **`io.Copy`** and **`io.MultiWriter`**, you can simplify complex I/O operations.

---

## Time

### Working with Time in Go

The `time` package in Go provides utilities for working with **time** and **durations**. The two main types you’ll use are `time.Duration` and `time.Time`.

### `time.Duration`

A **duration** represents a period of time, such as “2 hours and 30 minutes.” The smallest unit in Go's time package is the **nanosecond**, but you can work with larger units like seconds, minutes, and hours using constants.

#### Example:

```go
d := 2 * time.Hour + 30 * time.Minute
```

This creates a `time.Duration` of 2 hours and 30 minutes.

- You can also parse durations from strings using `time.ParseDuration`, like `"300ms"` or `"1h30m"`.

#### Example:

```go
duration, _ := time.ParseDuration("1h30m")
fmt.Println(duration) // Output: 1h30m0s
```

### Methods on `time.Duration`

- **Formatting**: `String()` returns the duration as a string (e.g., `"1h30m0s"`).
- **Conversion**: Methods like `Hours()`, `Minutes()`, `Seconds()` return the duration in the specified unit.
- **Truncate and Round**: These methods adjust the duration to a multiple of the provided unit.
  ```go
  d := time.Minute + 30 * time.Second
  fmt.Println(d.Truncate(time.Minute)) // Output: 1m0s
  fmt.Println(d.Round(time.Minute))    // Output: 1m0s
  ```

### `time.Time`

**`time.Time`** represents a specific point in time and is associated with a time zone. You get the current time using `time.Now()`.

#### Example:

```go
now := time.Now()
fmt.Println(now) // Output: 2024-09-07 15:04:05.123456 -0700 MST
```

### Comparing Times

- **Equal**: Use the `Equal()` method instead of `==` for comparing times, as it handles time zones.
  ```go
  t1 := time.Now()
  t2 := t1.Add(time.Minute)
  fmt.Println(t1.Equal(t2)) // Output: false
  ```
- **Before and After**: Use `Before()` and `After()` methods to compare times.
  ```go
  t1 := time.Now()
  t2 := t1.Add(time.Minute)
  fmt.Println(t1.Before(t2)) // Output: true
  ```

### Formatting and Parsing Dates

Go uses a unique date and time formatting system based on **January 2, 2006 at 3:04:05 PM (MST)**.

#### Formatting Example:

```go
t := time.Now()
fmt.Println(t.Format("2006-01-02 15:04:05")) // Output: 2024-09-07 15:04:05
```

#### Parsing Example:

```go
t, _ := time.Parse("2006-01-02 15:04:05", "2024-09-07 15:04:05")
fmt.Println(t) // Output: 2024-09-07 15:04:05 +0000 UTC
```

### Methods on `time.Time`

- **Extracting Values**: `Year()`, `Month()`, `Day()`, `Hour()`, `Minute()`, `Second()` allow you to extract specific parts of the time.

  ```go
  now := time.Now()
  fmt.Println(now.Year())   // Output: 2024
  fmt.Println(now.Month())  // Output: September
  fmt.Println(now.Day())    // Output: 7
  ```

- **Add and Subtract Time**: Use `Add()` to add a `time.Duration` to a `time.Time`. Use `Sub()` to find the difference between two times, returning a `time.Duration`.
  ```go
  now := time.Now()
  future := now.Add(2 * time.Hour)
  fmt.Println(future.Sub(now)) // Output: 2h0m0s
  ```

### Monotonic Time

Go uses **monotonic clocks** to track elapsed time, which is crucial because wall clock time can change due to daylight saving time, leap seconds, or NTP updates. Monotonic clocks always increase, making them useful for measuring elapsed time without disruptions.

- Monotonic time is used automatically when you call `time.Now()` and use `Sub()` to calculate time differences.

### Timers and Timeouts

Go’s `time` package provides utilities to handle timeouts and periodic tasks:

- **`time.After()`**: Returns a channel that receives a value after the specified duration.

  ```go
  select {
  case <-time.After(2 * time.Second):
      fmt.Println("Timeout after 2 seconds")
  }
  ```

- **`time.Tick()`**: Creates a channel that delivers values periodically (avoid using it in non-trivial programs, as it cannot be stopped).

  ```go
  for t := range time.Tick(time.Second) {
      fmt.Println(t)
  }
  ```

- **`time.NewTicker()`**: A better alternative to `time.Tick()`, as it provides a `Stop()` method to stop the ticker when it’s no longer needed.

  ```go
  ticker := time.NewTicker(1 * time.Second)
  go func() {
      for t := range ticker.C {
          fmt.Println("Tick at", t)
      }
  }()
  time.Sleep(5 * time.Second)
  ticker.Stop() // Stop the ticker
  ```

- **`time.AfterFunc()`**: Runs a function after a delay.
  ```go
  time.AfterFunc(2*time.Second, func() {
      fmt.Println("Function executed after 2 seconds")
  })
  time.Sleep(3 * time.Second) // Ensure the program doesn't exit before the function runs
  ```

### Conclusion

The `time` package in Go is powerful for working with durations, timeouts, and timestamps. It offers utilities for formatting, comparing, and manipulating time, while providing monotonic time for reliable time measurements.

---

## encoding/json

### JSON in Go: `encoding/json`

Go’s `encoding/json` package provides easy-to-use tools for working with JSON data. You can convert Go data structures to JSON (marshaling) and convert JSON back into Go data structures (unmarshaling).

#### Key Concepts:

- **Marshaling**: Converting Go data types to JSON.
- **Unmarshaling**: Converting JSON to Go data types.

### Struct Tags for JSON

Go uses **struct tags** to specify how fields should be marshaled/unmarshaled.

#### Example JSON:

```json
{
  "id": "12345",
  "date_ordered": "2020-05-01T13:01:02Z",
  "customer_id": "3",
  "items": [
    { "id": "xyz123", "name": "Thing 1" },
    { "id": "abc789", "name": "Thing 2" }
  ]
}
```

#### Go Structs:

```go
type Order struct {
 ID          string    `json:"id"`
 DateOrdered time.Time `json:"date_ordered"`
 CustomerID  string    `json:"customer_id"`
 Items       []Item    `json:"items"`
}

type Item struct {
 ID   string `json:"id"`
 Name string `json:"name"`
}
```

Struct tags like `json:"id"` ensure that the field names in JSON correspond to Go fields. Other struct tag options:

- **Omit empty fields**: `json:"customer_id,omitempty"`
- **Ignore field**: `json:"-"`

### Marshaling and Unmarshaling

- **Unmarshaling**: Convert JSON to Go struct.

```go
var order Order
err := json.Unmarshal([]byte(jsonData), &order)
```

- **Marshaling**: Convert Go struct to JSON.

```go
out, err := json.Marshal(order)
```

### JSON with Readers and Writers

You can use `json.Decoder` and `json.Encoder` for streaming JSON, allowing you to work with large data sets efficiently.

#### Example:

```go
tmpFile, err := os.Create("output.json")
err = json.NewEncoder(tmpFile).Encode(order)
```

### Handling JSON Streams

For multiple JSON objects in a stream:

```go
dec := json.NewDecoder(strings.NewReader(jsonStream))
for {
    err := dec.Decode(&t)
    if err == io.EOF {
        break
    }
    // process t
}
```

### Custom JSON Handling

To customize JSON parsing for fields like dates in non-standard formats, implement the `json.Marshaler` and `json.Unmarshaler` interfaces.

#### Custom Type Example:

```go
type RFC822ZTime struct {
    time.Time
}

func (rt RFC822ZTime) MarshalJSON() ([]byte, error) {
    out := rt.Time.Format(time.RFC822Z)
    return []byte(`"` + out + `"`), nil
}

func (rt *RFC822ZTime) UnmarshalJSON(b []byte) error {
    t, err := time.Parse(`"`+time.RFC822Z+`"`, string(b))
    *rt = RFC822ZTime{t}
    return err
}
```

This allows you to customize how the `DateOrdered` field is marshaled/unmarshaled:

```go
type Order struct {
 DateOrdered RFC822ZTime `json:"date_ordered"`
}
```

### Summary

Go’s `encoding/json` package makes it easy to work with JSON using struct tags for field mapping, decoding JSON from streams, and even customizing the marshaling/unmarshaling process for specific fields like dates.

---

## net/http

### Working with HTTP in Go: `net/http`

Go’s standard library includes a production-quality HTTP client and server. Here’s an overview of how to work with both, using the `net/http` package.

---

### HTTP Client

The `net/http` package provides the `http.Client` type for making HTTP requests.

#### Basic Client Example:

```go
client := &http.Client{
    Timeout: 30 * time.Second, // Set a request timeout
}

req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
if err != nil {
    panic(err)
}

req.Header.Add("X-Client", "MyGoApp") // Add custom headers

res, err := client.Do(req) // Make the request
if err != nil {
    panic(err)
}
defer res.Body.Close()

// Check the response status code
if res.StatusCode != http.StatusOK {
    panic(fmt.Sprintf("Unexpected status code: %v", res.StatusCode))
}

// Process response
var result map[string]interface{}
err = json.NewDecoder(res.Body).Decode(&result)
if err != nil {
    panic(err)
}

fmt.Printf("%+v\n", result) // Print decoded JSON response
```

#### Why Avoid `http.DefaultClient`?

- **No Timeout**: The default client lacks a timeout, making it unsuitable for production.
- **Solution**: Create a custom `http.Client` with a timeout.

---

### HTTP Server

The `http.Server` type handles incoming HTTP requests. You define request handling using the `http.Handler` interface.

#### Basic Server Example:

```go
type HelloHandler struct{}

func (hh HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!\n"))
}

s := &http.Server{
    Addr:         ":8080",
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 90 * time.Second,
    Handler:      HelloHandler{}, // Assign a handler
}

err := s.ListenAndServe()
if err != nil && err != http.ErrServerClosed {
    panic(err)
}
```

#### Using `http.ServeMux` for Routing:

You can define multiple routes with `http.ServeMux`.

```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello!\n"))
})

s := &http.Server{
    Addr:    ":8080",
    Handler: mux, // Use the mux for routing
}

s.ListenAndServe()
```

---

### Middleware in Go

Middleware is used for pre- and post-processing of requests. Middleware wraps a handler and can chain multiple handlers.

#### Example Middleware (Request Timer):

```go
func RequestTimer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r) // Call the next handler
        duration := time.Since(start)
        log.Printf("Request took %v", duration)
    })
}
```

#### Applying Middleware:

```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello!\n"))
})

wrappedMux := RequestTimer(mux)
s := &http.Server{
    Addr:    ":8080",
    Handler: wrappedMux,
}

s.ListenAndServe()
```

---

### Advanced Routing with Go 1.22

Go 1.22 introduced advanced routing features with HTTP verbs and path variables.

```go
mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    fmt.Fprintf(w, "Hello, %s!", name)
})
```

---

### Middleware Libraries

If you want to enhance middleware functionality, libraries like **alice** offer a more elegant way to chain middleware:

```go
chain := alice.New(RequestTimer).ThenFunc(helloHandler)
mux.Handle("/hello", chain)
```

---

### Conclusion

Go’s `net/http` package offers robust HTTP client and server implementations, allowing for concurrency, timeouts, and middleware. It’s suitable for building scalable web applications while providing flexibility with third-party libraries for enhanced functionality.

---

## Structured Logging

### Structured Logging in Go with `log/slog`

Go's standard logging package, `log`, is great for small programs but lacks support for structured logging. In Go 1.21, the `log/slog` package was introduced, providing a solution for creating structured logs. Here's a simplified guide to working with `log/slog`.

---

### Why Structured Logging?

Structured logs are logs with well-defined, consistent formats (like JSON) that make them easier to process programmatically. They’re vital in large-scale applications with millions of users, as they help automate log processing, detect patterns, and find issues efficiently.

---

### Simple Logging with `log/slog`

You can start logging structured messages using the default logger in `log/slog`.

```go
slog.Debug("debug log message")
slog.Info("info log message")
slog.Warn("warning log message")
slog.Error("error log message")
```

Output example:

```
2023/04/20 23:13:31 INFO info log message
2023/04/20 23:13:31 WARN warning log message
2023/04/20 23:13:31 ERROR error log message
```

Each log includes a timestamp, log level, and message.

---

### Adding Custom Fields to Logs

You can easily add key-value pairs to your log messages for structured data.

```go
userID := "fred"
loginCount := 20
slog.Info("user login", "id", userID, "login_count", loginCount)
```

Output example:

```
2023/04/20 23:36:38 INFO user login id=fred login_count=20
```

---

### JSON Log Output

To output logs in JSON format, create a custom logger.

```go
options := &slog.HandlerOptions{Level: slog.LevelDebug}
handler := slog.NewJSONHandler(os.Stderr, options)
mySlog := slog.New(handler)

mySlog.Debug("debug message", "id", "fred", "last_login", time.Now())
```

Output example (JSON):

```json
{
  "time": "2023-04-22T23:30:01.170243-04:00",
  "level": "DEBUG",
  "msg": "debug message",
  "id": "fred",
  "last_login": "2023-01-01T11:50:00Z"
}
```

---

### Optimized Logging with `LogAttrs`

For better performance and fewer memory allocations, you can use `LogAttrs`, which avoids allocating for every key-value pair.

```go
mySlog.LogAttrs(context.TODO(), slog.LevelInfo, "faster logging", slog.String("id", "fred"), slog.Time("last_login", time.Now()))
```

This method uses typed attributes (`slog.String`, `slog.Time`) instead of raw key-value pairs for improved efficiency.

---

### Integrating with `log.Logger`

If you need to bridge the older `log.Logger` with `log/slog`, you can do this easily using `NewLogLogger`.

```go
myLog := slog.NewLogLogger(mySlog.Handler(), slog.LevelDebug)
myLog.Println("using the mySlog Handler")
```

Output example:

```json
{
  "time": "2023-04-22T23:30:01.170269-04:00",
  "level": "DEBUG",
  "msg": "using the mySlog Handler"
}
```

---

### Summary of Key Features

- **Structured Logging**: Add fields to logs for better organization.
- **Log Levels**: Control the verbosity of logs (e.g., Debug, Info, Warn, Error).
- **JSON Logs**: Easily output logs in JSON format for better integration with log management systems.
- **Performance**: Use `LogAttrs` for optimized logging with fewer allocations.

`log/slog` is a powerful and flexible logging tool, making it easier to write structured logs in modern Go programs.

---

## Exercises

Here are simplified steps and code examples for each exercise:

### 1. Web Server that Returns the Current Time in RFC 3339 Format

The task here is to write a simple web server that responds with the current time in RFC 3339 format when a GET request is sent.

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC3339)
    w.Write([]byte(currentTime))
}

func main() {
    http.HandleFunc("/time", timeHandler)
    fmt.Println("Server is running on port 8080...")
    http.ListenAndServe(":8080", nil)
}
```

### 2. Middleware for Logging IP Address of Incoming Requests

This task involves writing a middleware component that logs the IP address of each incoming request using JSON structured logging.

```go
package main

import (
    "net/http"
    "os"
    "time"
    "golang.org/x/exp/slog"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr

        // Log the IP address using slog
        slog.New(slog.NewJSONHandler(os.Stdout, nil)).Info("Incoming request", "IP", ip, "timestamp", time.Now())

        next.ServeHTTP(w, r)
    })
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(time.RFC3339)
    w.Write([]byte(currentTime))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/time", timeHandler)

    wrappedMux := loggingMiddleware(mux)

    http.ListenAndServe(":8080", wrappedMux)
}
```

### 3. Add JSON Response Based on the Accept Header

This task extends the web server to return the current time as JSON or text based on the `Accept` header. By default, it will return plain text, but if the header `Accept: application/json` is sent, it will return a structured JSON response.

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
)

type TimeResponse struct {
    DayOfWeek   string `json:"day_of_week"`
    DayOfMonth  int    `json:"day_of_month"`
    Month       string `json:"month"`
    Year        int    `json:"year"`
    Hour        int    `json:"hour"`
    Minute      int    `json:"minute"`
    Second      int    `json:"second"`
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now()

    acceptHeader := r.Header.Get("Accept")

    if acceptHeader == "application/json" {
        timeResponse := TimeResponse{
            DayOfWeek:   currentTime.Weekday().String(),
            DayOfMonth:  currentTime.Day(),
            Month:       currentTime.Month().String(),
            Year:        currentTime.Year(),
            Hour:        currentTime.Hour(),
            Minute:      currentTime.Minute(),
            Second:      currentTime.Second(),
        }
        jsonResponse, _ := json.Marshal(timeResponse)
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResponse)
    } else {
        w.Write([]byte(currentTime.Format(time.RFC3339)))
    }
}

func main() {
    http.HandleFunc("/time", timeHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Summary

- **Exercise 1**: A simple web server that returns the current time in RFC 3339 format.
- **Exercise 2**: Middleware for logging the IP address using JSON structured logging.
- **Exercise 3**: Adds the ability to return the time as JSON or text, based on the `Accept` header.

You can run each code block to test them and modify them to suit your needs.
