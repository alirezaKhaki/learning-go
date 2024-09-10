# Concurrency in Go

**Summary:**

This chapter introduces the concept of concurrency in Go, which involves breaking down a process into independent components and coordinating their data sharing. Unlike traditional approaches using operating system threads and locks, Go employs a CSP-based model using goroutines, channels, and the select keyword.

**Notes:**

- **CSP:** A concurrency model inspired by Tony Hoare's work, focusing on message passing between processes.
- **Goroutines:** Lightweight, stackless threads managed by the Go runtime.
- **Channels:** Typed pipes used for communication between goroutines.
- **select:** A keyword for non-blocking multiplexing between channels.

**Conclusion:**

Go's concurrency model is a powerful and easy-to-understand alternative to traditional thread-based approaches. By leveraging goroutines, channels, and select, developers can create concurrent programs effectively and efficiently.

---

## When to Use Concurrency

### A Deeper Dive into Concurrency in Go

### Understanding Concurrency vs. Parallelism

While the terms are often used interchangeably, concurrency and parallelism are distinct concepts:

- **Concurrency:** This refers to the ability of a system to handle multiple tasks simultaneously, even if they aren't executed at the same time. Go's goroutines are a prime example of concurrency.
- **Parallelism:** This involves the simultaneous execution of multiple tasks on multiple processors or cores. While concurrency enables multiple tasks, parallelism ensures they're executed at the same time.

### Factors Influencing Concurrency Performance

Several factors influence whether concurrency will actually improve performance:

- **Amdahl's Law:** This law states that there's a limit to the speedup that can be achieved through parallelization. If a portion of the work must be done sequentially, it will limit the overall performance gain.
- **Task Granularity:** Smaller, independent tasks are more suitable for concurrency. If tasks are too large or interdependent, the overhead of coordinating them might outweigh the benefits.
- **Communication Overhead:** The cost of communication between concurrent tasks (e.g., through channels) can impact performance. Excessive communication can negate any gains from parallelism.
- **Hardware Capabilities:** The number of cores and their capabilities will determine how much parallelism can be achieved.

### Common Concurrency Patterns in Go

- **Producer-Consumer:** One or more goroutines produce data, and one or more goroutines consume it. Channels are often used for communication.
- **Fan-in/Fan-out:** Multiple goroutines perform tasks on a single input, and their results are combined into a single output.
- **Pipeline:** A series of goroutines, each passing data to the next in a pipeline fashion. This is common in data processing tasks.
- **Worker Pool:** A pool of goroutines is used to execute tasks from a shared queue. This can be efficient for handling many short-lived tasks.

### When to Avoid Concurrency

- **Simple Tasks:** If a task is straightforward and doesn't involve significant computation or I/O, concurrency might add overhead without providing benefits.
- **Sequential Dependencies:** If tasks are heavily dependent on each other, concurrency might introduce complexities and potential deadlocks.
- **Debugging Challenges:** Concurrent programs can be harder to debug due to non-deterministic behavior.

### Best Practices for Concurrency

- **Start Simple:** Begin with small, well-defined concurrent tasks to understand the concepts.
- **Use Channels Effectively:** Channels provide a safe and efficient way to communicate between goroutines.
- **Avoid Deadlocks:** Be mindful of potential deadlocks, especially when using channels and locks.
- **Test Thoroughly:** Test concurrent code under various conditions to ensure it behaves as expected.
- **Consider Alternatives:** If concurrency doesn't provide significant benefits or introduces complexities, explore other approaches like asynchronous programming.

By following these guidelines and carefully considering your specific use case, you can effectively leverage concurrency in Go to improve the performance and scalability of your applications.

---

## Goroutines

### Understanding Goroutines in Go

A **goroutine** is a core concept in Go's concurrency model. Think of it as a lightweight thread managed by Go's runtime. Goroutines make concurrent programming efficient, scalable, and easier to work with compared to traditional threads.

#### Key Concepts

1. **Process**: A process is an instance of a program running on a computer. The operating system gives each process its own resources (like memory) and ensures that other processes can’t interfere with them.

2. **Thread**: A thread is a unit of execution within a process. Multiple threads can share resources within the same process. The operating system schedules these threads across the CPU cores to ensure they all get a chance to run.

3. **Goroutine**: A goroutine is like a lightweight thread, but it's managed by Go’s runtime scheduler, not the operating system. Go's runtime handles the scheduling of goroutines onto threads, and it decides when to pause or resume them. Goroutines are faster and more memory-efficient than traditional threads.

#### Benefits of Goroutines

- **Faster creation**: Creating a goroutine is faster than creating an operating system thread because goroutines are not tied to OS-level resources.
- **Memory efficiency**: Goroutines have small initial stack sizes (as little as a few KB), and their stack can grow as needed. In contrast, threads require larger fixed-size stacks.
- **Faster context switching**: Switching between goroutines is much faster because it doesn't involve the operating system—everything happens in the Go runtime.
- **Optimized scheduling**: Go's scheduler interacts with other Go components, like the garbage collector and network poller, to make smart decisions about which goroutines to run and when.

#### Creating Goroutines

A goroutine is created by placing the `go` keyword before a function call. The function will run concurrently with the rest of the program. You can pass parameters to the function, but the function’s return values are ignored.

```go
go someFunction(param1, param2)  // Launches a goroutine
```

#### Example: Goroutines and Channels

In Go, it's common to use **goroutines** alongside **channels** for concurrent programming. Here’s an example where multiple goroutines process data concurrently and communicate via channels.

```go
func process(val int) int {
    // Simulate some work
    return val * 2
}

func processConcurrently(inVals []int) []int {
    in := make(chan int, len(inVals))  // Input channel
    out := make(chan int, len(inVals)) // Output channel

    // Launch 5 worker goroutines to process data
    for i := 0; i < 5; i++ {
        go func() {
            for val := range in {         // Read from input channel
                out <- process(val)       // Process and write to output channel
            }
        }()
    }

    // Load data into the input channel
    go func() {
        for _, val := range inVals {
            in <- val
        }
        close(in)  // Close the channel when done
    }()

    // Collect results from the output channel
    results := make([]int, len(inVals))
    for i := range results {
        results[i] = <-out
    }
    return results
}

func main() {
    inVals := []int{1, 2, 3, 4, 5}
    results := processConcurrently(inVals)
    fmt.Println(results)  // Output: [2, 4, 6, 8, 10]
}
```

#### Explanation:

- **Goroutines**: The `processConcurrently` function launches 5 worker goroutines, each reading from the `in` channel, processing the data using the `process` function, and sending the result to the `out` channel.
- **Channels**: Channels are used to communicate between the goroutines. The `in` channel is used to pass values to the worker goroutines, and the `out` channel is used to collect the results.
- **Concurrency**: Multiple values are processed concurrently, making the program more efficient.

### Key Features of Goroutines:

- **Non-blocking**: When a goroutine is launched, the main program continues executing without waiting for the goroutine to finish.
- **Closure Use**: It's common in Go to wrap business logic inside goroutines using anonymous functions or closures.

#### The "Function Coloring" Problem

Go avoids the "function coloring" problem (the need to differentiate between synchronous and asynchronous functions with different keywords, as in languages like JavaScript). In Go, any function can be made asynchronous by simply launching it with `go`, regardless of how it's defined.

### Conclusion

Goroutines allow Go to handle concurrency in a simple and efficient way. By combining goroutines with channels, you can build scalable, concurrent systems without worrying about complex thread management or locking mechanisms. Go’s scheduler optimizes the execution of goroutines, making it easy to run thousands of them concurrently without significant overhead.

---

## Understanding Channels in Go

**Channels** in Go are a key feature for communication between goroutines. They allow one goroutine to send data to another goroutine in a safe, synchronized way.

### Creating a Channel

Channels are created using the `make` function, similar to slices and maps. They are reference types, so when you pass a channel to a function, you're passing a reference to it.

```go
ch := make(chan int)  // Create a channel that transports integers
```

By default, the channel is unbuffered, meaning it can hold only one value at a time, and both sending and receiving operations block until the other side is ready.

### Reading from and Writing to a Channel

To **send data** into a channel, place the `<-` operator to the right of the channel variable. To **receive data**, place the `<-` operator to the left of the channel variable.

```go
ch <- 5        // Send the value 5 to the channel
val := <-ch    // Receive the value from the channel and store it in val
```

Each value sent to a channel can be read by only **one goroutine**.

### Channel Direction (Send-Only vs. Receive-Only)

Channels can be restricted to be **send-only** or **receive-only**. This is useful for ensuring that a channel is used properly in specific goroutines.

- **Send-only channel**:

  ```go
  ch chan<- int  // Only for sending data
  ```

- **Receive-only channel**:
  ```go
  ch <-chan int  // Only for receiving data
  ```

This helps the Go compiler enforce how channels are used, ensuring correctness in concurrent communication.

### Unbuffered Channels

By default, channels are **unbuffered**, meaning every **write** to the channel blocks until another goroutine reads from the channel, and every **read** blocks until there is data written to the channel. This forces synchronization between goroutines.

#### Example of an Unbuffered Channel

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int)

    go func() {
        ch <- 5  // Write to the channel
    }()

    val := <-ch  // Read from the channel
    fmt.Println(val)  // Output: 5
}
```

In this case, the goroutine writes `5` to the channel, and the main function reads the value from the channel.

### Buffered Channels

A **buffered channel** allows you to send a fixed number of values to the channel without requiring immediate reads. You specify the buffer size when creating the channel:

```go
ch := make(chan int, 10)  // Create a buffered channel with a capacity of 10
```

- Writing to a buffered channel only blocks if the buffer is full.
- Reading from a buffered channel only blocks if the buffer is empty.

#### Example of a Buffered Channel

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 2)  // Buffered channel with capacity of 2

    ch <- 1  // Doesn't block because the buffer isn't full
    ch <- 2  // Doesn't block

    fmt.Println(<-ch)  // Output: 1
    fmt.Println(<-ch)  // Output: 2
}
```

In this case, values are stored in the buffer, and the reads don't block unless the buffer is empty.

### Functions to Query Channels

- **`len(ch)`**: Returns the number of values currently in the channel buffer.
- **`cap(ch)`**: Returns the capacity of the channel buffer.

For unbuffered channels, both `len` and `cap` return 0 since unbuffered channels have no capacity to store values.

### Blocking and Synchronization

- **Unbuffered channels**: Always block until the other side is ready to send or receive data, enforcing synchronization between goroutines.
- **Buffered channels**: Block only when the buffer is full (for writing) or empty (for reading).

### Key Points:

- **Unbuffered channels** are used for tight synchronization between goroutines.
- **Buffered channels** allow for more flexible, asynchronous communication but block once the buffer is full or empty.
- Channels ensure safe data transfer between goroutines without explicit locks.

Channels are an essential part of Go's concurrency model, allowing goroutines to communicate in a safe and efficient manner, making it easier to build scalable concurrent applications.

---

## Using for-range and Channels

### Using `for-range` with Channels in Go

In Go, you can iterate over values sent through a channel using a `for-range` loop. This is a powerful feature that allows you to read values from a channel until it is closed. It works similarly to a `for-range` loop for slices or maps, but in this case, it only declares one variable, which holds the value received from the channel.

### How `for-range` with Channels Works

- **Loop until the channel is closed**: The loop runs until the channel is closed.
- **Pausing behavior**: If no value is available on the channel, the loop pauses the goroutine until a value is sent to the channel or the channel is closed.

### Syntax

```go
for v := range ch {
    fmt.Println(v)
}
```

- `v` holds the value received from the channel `ch`.
- The loop continues reading values from the channel until it is closed.

### Example of `for-range` with Channels

```go
package main

import (
    "fmt"
)

func sendValues(ch chan int) {
    for i := 1; i <= 5; i++ {
        ch <- i  // Send values to the channel
    }
    close(ch)  // Close the channel when done
}

func main() {
    ch := make(chan int)

    go sendValues(ch)  // Run the sender goroutine

    // Use for-range to receive values from the channel
    for v := range ch {
        fmt.Println(v)  // Output: 1 2 3 4 5
    }
}
```

### Explanation:

- **Goroutine**: The `sendValues` function sends integers from 1 to 5 into the channel `ch` and then closes the channel.
- **for-range**: The main goroutine uses a `for-range` loop to receive and print the values from the channel. The loop automatically stops when the channel is closed.

### Important Points:

1. **Channel closure**: The `for-range` loop only stops when the channel is **closed**. If the channel is not closed, the loop will hang, waiting for more values.
2. **No explicit check needed**: You don’t need to check for the closed channel condition explicitly. The `for-range` loop handles this internally and stops automatically when the channel is closed.
3. **Closing a channel**: Only the sender should close a channel. It is a runtime error to send on a closed channel.

### Manual Alternative to `for-range`

Without `for-range`, you'd typically use a `for` loop and check the second value returned by receiving from the channel, which indicates whether the channel is open or closed:

```go
for {
    v, ok := <-ch
    if !ok {
        break  // Channel is closed
    }
    fmt.Println(v)
}
```

However, `for-range` is simpler and more idiomatic when iterating over channel values.

### Summary

- **`for-range`** is a convenient and idiomatic way to receive values from a channel until it’s closed.
- The loop pauses if no value is available on the channel and resumes when a value is sent.
- The loop automatically stops when the channel is closed, making it an efficient way to read all values from a channel.

---

## Closing a Channel

### Closing a Channel in Go

In Go, when you're done writing to a channel, you **close** it using the `close` function. Closing a channel signals that no more values will be sent on that channel.

```go
close(ch)  // Close the channel
```

### Important Behaviors of a Closed Channel

1. **Cannot Write to a Closed Channel**:

   - Writing to a closed channel will cause a **panic**.
   - You can only **close** a channel once; trying to close it again will also cause a panic.

2. **Reading from a Closed Channel**:
   - Reading from a closed channel **never blocks**.
   - If the channel is **buffered**, any remaining values in the buffer can still be read.
   - Once all buffered values are read (or if the channel is unbuffered), the channel returns the **zero value** for its type.

### Example of Reading from a Closed Channel:

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 3)

    ch <- 1
    ch <- 2
    ch <- 3

    close(ch)  // Close the channel

    // Read all remaining values from the channel
    for i := 0; i < 3; i++ {
        fmt.Println(<-ch)  // Output: 1, 2, 3
    }

    // Reading from the closed channel returns the zero value (0 for int)
    fmt.Println(<-ch)  // Output: 0
}
```

### Detecting if a Channel is Closed Using `comma ok` Idiom

When reading from a channel, if you need to detect whether the channel is closed, you can use the **comma ok** idiom:

```go
v, ok := <-ch
```

- `v` is the value received from the channel.
- `ok` is `true` if the channel is still open, and `false` if the channel is closed.

### Example of Using `comma ok`:

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2

    close(ch)  // Close the channel

    // Read from the channel using comma ok
    for {
        v, ok := <-ch
        if !ok {
            fmt.Println("Channel closed")
            break
        }
        fmt.Println(v)
    }
}
```

**Output:**

```
1
2
Channel closed
```

In this example, after the channel is closed, the loop continues until `ok` is `false`, signaling that no more values are available.

### Key Points:

1. **Closing a Channel**:
   - The responsibility for closing a channel lies with the goroutine that sends values to it.
   - Closing a channel is **not mandatory** unless another goroutine is explicitly waiting for the channel to close (e.g., in a `for-range` loop).
2. **The `comma ok` Idiom**:

   - Use the `comma ok` idiom when reading from a channel to check if the channel is still open.
   - If `ok` is `false`, the channel is closed, and `v` will be the zero value for the channel's type.

3. **Buffered Channels**:

   - If a buffered channel is closed and still has values in the buffer, those values can still be read.
   - Once the buffer is empty, further reads return the zero value.

4. **Concurrency Model**:
   - Channels enforce clear communication between goroutines, making data dependencies explicit. This simplifies reasoning about concurrency.
   - In contrast to other languages, which often rely on shared state and locks, Go channels provide a safer and clearer way to manage concurrent processes.

By closing channels when needed and using the `comma ok` idiom, Go allows for safe, efficient, and clear communication between goroutines.

---

## Understanding How Channels Behave

### **Channel Behaviors:**

| **State**              | **Read**                                                                                                       | **Write**                       | **Close**                                           |
| ---------------------- | -------------------------------------------------------------------------------------------------------------- | ------------------------------- | --------------------------------------------------- |
| **Unbuffered, open**   | Pauses until something is written.                                                                             | Pauses until something is read. | Works (allows closing).                             |
| **Unbuffered, closed** | Returns the zero value for the channel’s type (use `comma ok` to check if closed).                             | Causes a panic (`panic`).       | Causes a panic (`panic`).                           |
| **Buffered, open**     | Pauses if the buffer is empty.                                                                                 | Pauses if the buffer is full.   | Works, and any remaining values stay in the buffer. |
| **Buffered, closed**   | Returns a remaining value in the buffer; if empty, returns the zero value (use `comma ok` to check if closed). | Causes a panic (`panic`).       | Causes a panic (`panic`).                           |
| **Nil**                | Hangs forever.                                                                                                 | Hangs forever.                  | Causes a panic (`panic`).                           |

### **Explanation:**

1. **Unbuffered, open channel**:

   - **Read**: Pauses until something is written.
   - **Write**: Pauses until something is read.
   - **Close**: Works without issues.

2. **Unbuffered, closed channel**:

   - **Read**: Always returns the zero value for the channel type, so use the `comma ok` idiom to check if the channel is closed.
   - **Write**: Causes a **panic** since writing to a closed channel is not allowed.
   - **Close**: Causes a **panic** if you try to close the channel more than once.

3. **Buffered, open channel**:

   - **Read**: Pauses if the buffer is empty until something is written.
   - **Write**: Pauses if the buffer is full until something is read.
   - **Close**: Works fine, leaving any buffered values available for reading.

4. **Buffered, closed channel**:

   - **Read**: Continues returning remaining values in the buffer. Once empty, returns the zero value for the channel type. Use `comma ok` to check if the channel is closed.
   - **Write**: Causes a **panic**.
   - **Close**: Causes a **panic** if the channel is closed again.

5. **Nil channel**:
   - **Read**: Hangs forever, since a nil channel can never be read from.
   - **Write**: Hangs forever, as a nil channel cannot accept writes.
   - **Close**: Causes a **panic** when trying to close a nil channel.

### Key Takeaways:

- Closing a channel is **safe** and **required** only by the sender (the one writing to the channel), and it is generally done when no more data will be sent.
- Attempting to **write** to a closed channel or closing a channel multiple times will cause a **panic**.
- Use the **`comma ok`** idiom when reading from a channel that might be closed to differentiate between actual zero values and those returned due to the channel being closed.

---

## Select

### Understanding the `select` Statement in Go

The **`select`** statement is a powerful feature in Go's concurrency model. It allows a goroutine to wait on multiple communication operations, such as reading from or writing to channels. The primary purpose of `select` is to handle multiple concurrent channel operations elegantly, preventing issues like **starvation** (favoring one operation over others) and **deadlocks**.

### Basic Syntax

The `select` statement looks similar to a `switch` but works with channels:

```go
select {
case v := <-ch1:
    fmt.Println(v)
case v := <-ch2:
    fmt.Println(v)
case ch3 <- x:
    fmt.Println("wrote", x)
case <-ch4:
    fmt.Println("received signal on ch4")
}
```

- **Cases**: Each case involves either reading from or writing to a channel.
- **Random selection**: If multiple cases are ready (i.e., multiple channels can read or write), Go **randomly** chooses one case to proceed, ensuring fairness and preventing starvation.
- **Blocking**: If no cases are ready, the `select` blocks and waits for a channel to become available.

### Preventing Deadlocks with `select`

A **deadlock** occurs when two or more goroutines are blocked, waiting for each other, and no progress can be made. If no channels can proceed in a Go program, the runtime detects this and terminates the program with an error.

#### Example of Deadlock:

```go
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        inGoroutine := 1
        ch1 <- inGoroutine
        fromMain := <-ch2
        fmt.Println("goroutine:", inGoroutine, fromMain)
    }()

    inMain := 2
    ch2 <- inMain
    fromGoroutine := <-ch1
    fmt.Println("main:", inMain, fromGoroutine)
}
```

- **What happens**: Both the main goroutine and the launched goroutine are waiting to read and write from each other’s channels. This leads to a **deadlock**, as neither can proceed.
- **Error**: `"fatal error: all goroutines are asleep - deadlock!"`

#### Avoiding Deadlock with `select`:

By wrapping the read and write operations in a `select`, you can prevent the deadlock:

```go
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        inGoroutine := 1
        ch1 <- inGoroutine
        fromMain := <-ch2
        fmt.Println("goroutine:", inGoroutine, fromMain)
    }()

    inMain := 2
    var fromGoroutine int
    select {
    case ch2 <- inMain:
    case fromGoroutine = <-ch1:
    }

    fmt.Println("main:", inMain, fromGoroutine)
}
```

- **Why it works**: The `select` checks if either channel can proceed. The program avoids deadlock by handling the first available channel operation.
- **Note**: In this version, the `fmt.Println` in the launched goroutine doesn't execute because the main goroutine exits early.

### Common Patterns with `select`

#### 1. **Using `select` in a `for` Loop (For-Select Loop)**

This is a common pattern used for continuously listening to multiple channels:

```go
for {
    select {
    case <-done:
        return  // Exit the loop when done is signaled
    case v := <-ch:
        fmt.Println(v)
    }
}
```

- **`for-select` loop**: This combination is used to continuously listen to channels until some exit condition is met. The loop will block on the `select` statement until a channel operation can proceed.
- **Exiting the loop**: Be sure to have a mechanism to break the loop, such as a signal channel (`done`), otherwise the loop will run indefinitely.

#### 2. **Using `select` with `default` for Non-Blocking Operations**

If you want a **non-blocking** read or write on a channel, use the `default` case:

```go
select {
case v := <-ch:
    fmt.Println("read from ch:", v)
default:
    fmt.Println("no value written to ch")
}
```

- **Non-blocking**: The `default` case is selected when none of the other cases can proceed. This allows the code to move forward without waiting.
- **Caution**: Using `default` inside a `for-select` loop can cause the loop to run continuously, which is inefficient and can waste CPU resources.

### Example of Using `select` to Avoid Deadlock and Handle Multiple Channels

Here’s an example demonstrating the correct use of `select` for handling multiple channels:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    // Goroutine 1: Sends a message after 1 second
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "message from ch1"
    }()

    // Goroutine 2: Sends a message after 2 seconds
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "message from ch2"
    }()

    // Use select to handle both channels
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        }
    }
}
```

- **Multiple channel handling**: The `select` statement allows the program to handle whichever channel receives a message first.
- **Avoids deadlock**: Both channels are listened to concurrently, preventing the goroutines from being blocked indefinitely.

### Summary

- **`select`** is Go's concurrency control structure, used to handle multiple channel operations in a non-blocking way.
- It picks randomly from the cases that are ready to avoid starvation.
- `select` can prevent deadlocks by enabling flexible channel communication patterns, but it's crucial to handle all possible cases correctly.
- The **for-select loop** is a common pattern used to keep reading from or writing to channels in a loop, and you can add a **default** case for non-blocking behavior.

---

## Concurrency Practices and Patterns

### Concurrency Best Practices and Patterns in Go

#### 1. **Keep APIs Concurrency-Free**

When designing APIs in Go, **hide concurrency details** such as channels and mutexes. Exposing these details forces users of your API to manage them, which can lead to errors like deadlocks or improper usage. Your API should be easy to use without requiring users to know about underlying concurrency mechanisms.

- **Avoid exposing channels and mutexes**: Channels should not be part of exported types, functions, or methods.
- **Exception**: If your API is specifically designed to handle concurrency (e.g., concurrency helper libraries), channels may be part of the API.

#### 2. **Goroutines, For Loops, and Capturing Variables**

Before Go 1.22, there was a common issue when launching goroutines inside a `for` loop. In earlier versions, the `for` loop reused the same variable for each iteration, which caused unexpected behavior when capturing variables inside a goroutine.

##### Example of the Issue (Go 1.21 and Earlier):

```go
func main() {
    a := []int{2, 4, 6, 8, 10}
    ch := make(chan int, len(a))

    for _, v := range a {
        go func() {
            ch <- v * 2
        }()
    }

    for i := 0; i < len(a); i++ {
        fmt.Println(<-ch)
    }
}
```

In this code, you expect each goroutine to capture a different value of `v`. However, all the goroutines capture the same final value of `v` (which is `10`), leading to this incorrect output:

```
20
20
20
20
20
```

##### Go 1.22 Behavior:

Starting in Go 1.22, the loop now **creates a new variable for each iteration**, which resolves this issue. Running the same code in Go 1.22 or later produces the correct output:

```
20
8
4
12
16
```

##### Fixing the Issue for Go 1.21 and Earlier:

If you're using Go 1.21 or earlier, you can fix this issue by creating a **copy of the variable** inside the loop.

1. **Variable Shadowing**: Create a new local variable with the same name inside the loop to ensure each goroutine captures a unique value.

   ```go
   for _, v := range a {
       v := v  // Shadowing the loop variable
       go func() {
           ch <- v * 2
       }()
   }
   ```

2. **Passing as a Parameter**: Pass the loop variable as a parameter to the goroutine.

   ```go
   for _, v := range a {
       go func(val int) {
           ch <- val * 2
       }(v)  // Pass v as a parameter to the closure
   }
   ```

#### 3. **Closures and Capturing Variables**

Even though Go 1.22 resolves the issue for loop variables, you still need to be careful with other variables captured by closures. If a closure depends on a variable that changes outside the closure, the closure may not capture the expected value.

- **Best Practice**: Use parameters to pass the current value of a variable into the closure, ensuring each closure gets its own unique copy.

##### Example:

```go
func main() {
    var someVar int

    go func() {
        fmt.Println(someVar)  // This might print an unexpected value if someVar changes
    }()

    someVar = 10
}
```

In the example above, `someVar` may change before the goroutine prints its value. To avoid this, pass `someVar` as a parameter to the closure:

```go
func main() {
    someVar := 5

    go func(val int) {
        fmt.Println(val)  // Now it prints the expected value
    }(someVar)

    someVar = 10
}
```

### Summary of Best Practices:

- **Hide concurrency in APIs**: Channels and mutexes should be internal implementation details and not exposed to API users.
- **For-loop variables in goroutines**: Go 1.22 fixes the issue of capturing variables in loops, but if you're using Go 1.21 or earlier, copy or pass loop variables explicitly.
- **Use parameters for closures**: When closures capture variables, use parameters to ensure they capture the correct value.
- **Concurrency should be an internal detail**: Only expose concurrency tools (like channels) if your library specifically deals with concurrency helpers.

---

## Always Clean Up Your Goroutines

### Preventing Goroutine Leaks in Go

When launching a **goroutine**, you must ensure that it will eventually **exit**. If a goroutine continues running indefinitely or blocks forever, it results in a **goroutine leak**. This can lead to memory that cannot be garbage collected because the runtime is unaware that the goroutine will no longer be used. Unlike variables, goroutines are not automatically freed by the garbage collector once they are done, so it's essential to manage them properly.

#### Example of a Goroutine That May Leak

Here's a simple function that launches a goroutine to generate numbers and sends them over a channel:

```go
func countTo(max int) <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < max; i++ {
            ch <- i
        }
        close(ch)  // Always close the channel after the loop completes
    }()
    return ch
}

func main() {
    for i := range countTo(10) {
        fmt.Println(i)
    }
}
```

In this example:

- The goroutine runs in the background, sending numbers to the channel `ch`.
- The main function reads from the channel until it is closed.

This works as expected when all the values are read. The goroutine sends all values and exits after closing the channel. However, there is a risk if the loop in `main` exits early.

#### The Problem: Goroutine Leak

If the loop in `main` exits before reading all the values from the channel, the goroutine continues trying to send values, but there’s no receiver anymore. This causes the goroutine to **block forever**, leading to a **goroutine leak**.

Example where the loop exits early:

```go
func main() {
    for i := range countTo(10) {
        if i > 5 {
            break  // Exit the loop early
        }
        fmt.Println(i)
    }
}
```

In this case, after the loop breaks, the goroutine continues running, trying to send more values to the channel. Since the `main` function has stopped reading from the channel, the goroutine is blocked forever, causing a memory leak.

### Solution: Use a `done` Channel to Signal Exit

A common way to prevent goroutine leaks is to use a **`done` channel** to signal when a goroutine should exit early.

Here’s how to modify the previous example:

```go
func countTo(max int, done <-chan struct{}) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)  // Ensure the channel is closed when the goroutine exits
        for i := 0; i < max; i++ {
            select {
            case ch <- i:  // Send value to channel if not done
            case <-done:   // Exit early if done signal is received
                return
            }
        }
    }()
    return ch
}

func main() {
    done := make(chan struct{})
    defer close(done)  // Ensure the done channel is closed to avoid leaks

    for i := range countTo(10, done) {
        if i > 5 {
            break  // Exit the loop early
        }
        fmt.Println(i)
    }
}
```

#### Explanation:

1. **`done` channel**: The `done` channel is used to signal the goroutine to stop early.
2. **`select` block**: Inside the goroutine, we use `select` to check if either:
   - A value should be sent to the channel, or
   - The `done` signal is received, which stops the goroutine early.
3. **`defer close(ch)`**: Ensures that the output channel is closed properly when the goroutine exits, whether it completes naturally or exits early.

#### How This Solves the Problem:

- When the loop in `main` exits early, it closes the `done` channel.
- The goroutine receives the signal from the `done` channel and exits, preventing it from blocking and leaking memory.

### Summary of Best Practices:

1. **Always ensure that goroutines exit**: When you launch a goroutine, make sure it can complete its work or be stopped when no longer needed.
2. **Use a `done` channel**: When launching long-running or background goroutines, provide a way to signal them to stop early if they are no longer needed.
3. **Close channels**: Ensure that channels are closed properly to avoid blocking readers and leaking goroutines.

By managing goroutine lifecycles carefully and using `done` channels, you can avoid goroutine leaks and improve the efficiency and reliability of your Go programs.
