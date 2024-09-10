# Go Tooling

## **Go's Robust Tooling Ecosystem**

Go's success lies not only in its language features but also in its comprehensive tooling ecosystem. This tooling is designed to streamline development tasks, making it easier for developers to create high-quality software.

**Key Tools and Features:**

- **Built-in Tools:** Go provides essential tools like `go vet`, `go fmt`, `go mod`, `go get`, `go list`, `go work`, `go doc`, and `go build` to handle common development tasks.
- **Testing Support:** The `go test` tool offers extensive testing capabilities, covered in detail in Chapter 15.
- **Third-Party Tools:** Beyond the built-in tools, Go benefits from a vibrant ecosystem of third-party tools that enhance development workflows.

**Conclusion**

Go's strong emphasis on tooling is a major factor in its popularity. By providing a robust set of tools and a thriving ecosystem, Go empowers developers to build efficient and maintainable software.

---

## Using go run to Try Out Small Programs

**Go's `go run` Command: A Quick and Easy Way to Execute Small Programs**

While Go is a compiled language, the `go run` command offers a convenient way to test and execute small programs without the need for explicit compilation.

**How `go run` Works:**

1. **Compilation:** The `go run` command compiles the specified Go source code into a temporary binary file.
2. **Execution:** It then executes the compiled binary.
3. **Cleanup:** After the program finishes, the temporary binary is deleted.

**Key Points:**

- **Rapid Feedback:** `go run` allows for a quick development and testing cycle, similar to interpreted languages.
- **No Permanent Binary:** The compiled binary is created and deleted in a temporary directory, keeping your project directory clean.
- **Ideal for Small Programs:** `go run` is well-suited for testing small scripts or exploring Go's features interactively.

**Conclusion**

The `go run` command provides a valuable tool for Go developers, enabling them to experiment with code and get immediate results. It's a convenient way to leverage Go's power for small-scale tasks and prototyping.

---

## Adding Third-Party Tools with go install

**Managing Third-Party Go Tools with `go install`**

This section explains how to leverage the `go install` command to install third-party Go tools from their source code repositories.

**Installation Process:**

1. **Specify the Package:** Use the `go install` command followed by the path to the main package in the tool's source code repository.
2. **Versioning:** Include the `@version` or `@latest` flag to specify the desired version or get the most recent version, respectively. Omitting this can lead to unexpected behavior.
3. **Download and Compile:** `go install` downloads the tool's source code, its dependencies, compiles everything, and installs the binary.
4. **Installation Location:** By default, binaries are placed in the `go/bin` directory within your home directory. You can modify the `GOBIN` environment variable to change this location.
5. **Executable Path:** Adding the `go/bin` directory (or your custom `GOBIN` location) to your system's path allows you to run the installed tool directly from the command line.

**Example:**

Installing the `hey` load testing tool:

```bash
$ go install github.com/rakyll/hey@latest
```

This downloads `hey`, its dependencies, builds the binary, and installs it in your `go/bin` directory.

**Updating Tools:**

To update an already installed tool to a newer version, simply rerun `go install` with the desired version or `@latest`.

**Key Points:**

- `go install` is convenient for managing third-party Go tools.
- Versioning is crucial to avoid installation issues.
- You can customize the installation location using `GOBIN`.
- Add the installation directory to your path for easy tool execution.
- Installed tools are regular binaries and can be moved or distributed as needed.

**Conclusion:**

The `go install` command simplifies the process of installing and managing third-party Go tools. This approach is widely adopted for distributing developer tools within the Go community.

---

## Improving Import Formatting with goimports

**Keeping Your Imports Tidy with `goimports`**

This section introduces `goimports`, a tool that enhances import formatting in your Go projects.

**What `goimports` Does:**

- **Orders Imports:** Organizes your import statements alphabetically.
- **Removes Unused Imports:** Eliminates imports that are not referenced in your code.
- **Suggests Missing Imports:** Attempts to automatically identify and add necessary imports. (Note that these suggestions may require manual verification.)

**Installation:**

Use `go install` to download `goimports`:

```bash
$ go install golang.org/x/tools/cmd/goimports@latest
```

**Running `goimports`:**

- **Flags:**
  - `-l`: Prints files with incorrect formatting to the console for review.
  - `-w`: Modifies the files directly to incorporate the suggested changes.
- **Target:** The `.` specifies scanning all files in the current directory and its subdirectories.

**Example:**

```bash
$ goimports -l -w .
```

This command runs `goimports` with the `-l` and `-w` flags on all files in your current project. It will list files with formatting issues and modify them in place with the suggested improvements.

**Important Notes:**

- `golang.org/x` packages are part of the Go project but are not part of the core Go standard library. They may have looser compatibility requirements and introduce breaking changes.
- `goimports` suggestions for missing imports might not always be accurate. Manual verification is recommended.

**Conclusion:**

Using `goimports` helps maintain clean and organized import statements in your Go projects, improving code readability and maintainability.

---

## Using Code-Quality Scanners

**Beyond `go vet`: Using Linters for Code Quality**

This section discusses the use of linters to enhance code quality beyond the basic checks provided by `go vet`.

**Linters and Code Style:**

- **Additional Checks:** Linters can identify potential programming errors, style inconsistencies, and non-idiomatic code that `go vet` might miss.
- **Examples:** Common suggestions include proper variable naming, formatted error messages, and comments on public methods and types.

**False Positives and Negatives:**

- **Verification:** While linters are valuable tools, they may occasionally produce inaccurate results (false positives or negatives).
- **Critical Thinking:** It's essential to evaluate linter suggestions and make informed decisions about whether to implement them.

**Ignoring Linter Suggestions:**

- **Comments:** Most linters allow you to add comments to your code to suppress specific warnings.
- **Explanation:** Include a clear explanation in the comment to justify your decision and provide context for future review.

**Conclusion:**

Linters are a valuable asset for maintaining high-quality Go code. By using them in conjunction with `go vet` and carefully evaluating their suggestions, you can write more robust, readable, and idiomatic code.

---

### **Tools**

**Enhancing Code Quality with Static Analysis Tools (including code samples)**

This section explores several static analysis tools that go beyond `go vet` to identify potential issues and enforce code quality standards.

**Key Tools:**

- **staticcheck:**
  - Focus: Wide range of code quality checks (over 150) with minimal false positives.
  - Installation: `go install honnef.co/go/tools/cmd/staticcheck@latest`
  - Example: Identifies unnecessary function calls like `fmt.Sprintf(string)`.
  - Code Sample:

```go
package main

import "fmt"

func main() {
  s := fmt.Sprintf("Hello")
  fmt.Println(s)
}
```

    * Explanation: `staticcheck` flags the use of `fmt.Sprintf` for a single string argument, suggesting a simpler approach using the string directly.

- **revive:**
  - Focus: Style and code quality checks based on the former `golint` tool.
  - Installation: `go install github.com/mgechev/revive@latest`
  - Example: Detects issues like missing comments, inconsistent naming conventions, and misplaced error checks.
  - Code Sample:

```go
package main

import "fmt"

func main() {
  true := false  // Shadowing built-in identifier
  fmt.Println(true)
}
```

    * Explanation: `revive` warns about assigning `false` to a variable named `true`, which shadows the built-in `true` identifier.

- **golangci-lint:**
  - Focus: Comprehensive linting suite running over 50 tools, including `go vet`, `staticcheck`, and `revive`.
  - Installation: Download binary from website (recommended).
  - Benefits:
    - Efficient execution of multiple linters.
    - Customizable configuration (e.g., enabling shadowing checks beyond revive).
  - Example: Flags potential issues like unused variables and identifiers shadowing built-in functions.
  - Code Sample 1:

```go
package main

func main() {
  x := 10
  x = 30  // Unused assignment
}
```

    * Explanation: `golangci-lint` identifies the assignment to `x` as ineffective since the value is never used.

- **Code Sample 2:**

```go
package main

import "fmt"

var b = 20

func main() {
  true := false  // Shadowing built-in identifier
  a := 10
  b := 30
  if true {
    a := 20  // Shadowing variable a
    fmt.Println(a)
  }
  fmt.Println(a, b)
}
```

    * Explanation: `golangci-lint` detects shadowing of both the built-in `true` identifier and the variable `a` within the `if` block.

**General Recommendations:**

- **Start with `go vet`:** Make it mandatory in your build process.
- **Add `staticcheck` next:** Offers valuable checks with low false positives.
- **Consider `revive` for stylistic preferences:** Be cautious of potential false positives/negatives.
- **Explore `golangci-lint` gradually:** Experiment with configuration to best suit your team's needs.
- **Maintain consistent configuration:** Version control the `.golangci.yml` file for consistent linting across developers.

**Conclusion:**

By incorporating these static analysis tools into your workflow, you can significantly improve code quality, maintain consistent style, and catch potential bugs before they cause problems. Remember to evaluate linter suggestions critically and establish configuration that works for your team to avoid unnecessary arguments and ensure effective code reviews.

---

## Using govulncheck to Scan for Vulnerable Dependencies

This section explains how to leverage `govulncheck` to scan your Go projects for vulnerabilities within dependencies.

**Understanding Software Vulnerabilities:**

- Third-party libraries can introduce security risks if they contain unpatched vulnerabilities.
- Developers fix these vulnerabilities in updates, but it's crucial to ensure your project uses the fixed versions.

**govulncheck: A Go Vulnerability Scanner**

- Developed by the Go team to address dependency security concerns.
- Scans project dependencies for known vulnerabilities in both the standard library and third-party modules.
- Relies on a publicly maintained vulnerability database by the Go team.

**Installation:**

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

**Example Usage:**

1. **Code Snippet:**

   ```go
   func main() {
       info := Info{}
       err := yaml.Unmarshal([]byte(data), &info)
       if err != nil {
           fmt.Printf("error: %v\n", err)
           os.Exit(1)
       }
       fmt.Printf("%+v\n", info)
   }
   ```

2. **go.mod File:**

   ```
   module github.com/learning-go-book-2e/vulnerable

   go 1.20
   require gopkg.in/yaml.v2 v2.2.7
   require gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
   ```

3. **Running govulncheck:**

   ```bash
   govulncheck ./...
   ```

   **Expected Output:**

   ```
   Using go1.21 and govulncheck@v1.0.0 with vulnerability data from
   https://vuln.go.dev (last modified 2023-07-27 20:09:46 +0000 UTC).
   Scanning your code and 49 packages across 1 dependent module
   for known vulnerabilities...
   Vulnerability #1: GO-2020-0036
   Excessive resource consumption in YAML parsing in gopkg.in/yaml.v2
   More info: https://pkg.go.dev/vuln/GO-2020-0036
   Module: gopkg.in/yaml.v2
   Found in: gopkg.in/yaml.v2@v2.2.7
   Fixed in: gopkg.in/yaml.v2@v2.2.8
   Example traces found:
   #1: main.go:25:23: vulnerable.main calls yaml.Unmarshal
   Your code is affected by 1 vulnerability from 1 module.
   ```

   - The output identifies a vulnerability (GO-2020-0036) in the `gopkg.in/yaml.v2` library (version 2.2.7).
   - It suggests updating to version 2.2.8 for a fix and highlights the specific line in your code that uses the vulnerable code.

**Resolving the Vulnerability:**

1. Update the dependency:

   ```bash
   go get -u=patch gopkg.in/yaml.v2
   ```

2. Rerun `govulncheck`:

   ```bash
   govulncheck ./...
   ```

   **Expected Output:**

   ```
   Using go1.21and govulncheck@v1.0.0 with vulnerability data from
   https://vuln.go.dev (last modified 2023-07-27 20:09:46 +0000 UTC).
   Scanning your code and 49 packages across 1 dependent module
   for known vulnerabilities...
   No vulnerabilities found.
   ```

   - The update successfully resolves the identified vulnerability.

**Best Practices:**

- Integrate `govulncheck` into your build process for regular vulnerability scanning.
- Strive for minimal dependency changes to minimize potential breakage.
- Update to the latest patch version within a minor release for security fixes.

**Additional Notes:**

- `govulncheck` might eventually be part of the standard Go toolset.
- Refer to the official blog post for more details: (link to be added if available).

---

## Embedding Content into Go Programs

This section explores the `go:embed` directive for embedding files and directories directly into your Go binaries.

**Key Points:**

- **Purpose:** Avoid distributing separate support files with your program.
- **Mechanism:** Uses `go:embed` comments to embed file contents into variables.
- **Variable Types:**
  - `string` or `[]byte` for single files.
  - `embed.FS` for directories or multiple files.
- **Usage:**
  - Import the `embed` package.
  - Place `go:embed` comments before package-level variables.
  - Access embedded content using the variable.

**Example: Embedding a Text File**

```go
package main

import (
    _ "embed"
    "fmt"
    "os"
    "strings"
)

//go:embed passwords.txt
var passwords string

func main() {
    pwds := strings.Split(passwords, "\n")
    if len(os.Args) > 1 {
        for _, v := range pwds {
            if v == os.Args[1] {
                fmt.Println("true")
                os.Exit(0)
            }
        }
        fmt.Println("false")
    }
}
```

**Example: Embedding a Directory**

```go
package main

import (
    "embed"
    "fmt"
    "io/fs"
    "os"
    "strings"
)

//go:embed help
var helpInfo embed.FS

func main() {
    if len(os.Args) == 1 {
        printHelpFiles()
        os.Exit(0)
    }
    data, err := helpInfo.ReadFile("help/" + os.Args[1])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(string(data))
}

func printHelpFiles() {
    fmt.Println("contents:")
    fs.WalkDir(helpInfo, "help",
        func(path string, d fs.DirEntry, err error) error {
            if !d.IsDir() {
                _, fileName, _ := strings.Cut(path, "/")
                fmt.Println(fileName)
            }
            return nil
        })
}
```

**Additional Notes:**

- Use `embed.FS` to represent a virtual filesystem for directories.
- The `WalkDir` function can be used to iterate over embedded files and directories.
- You can embed binary files and use wildcards or ranges to specify file patterns.
- The compiler checks for invalid embedding specifications.

**Conclusion:**

`go:embed` provides a convenient way to include static content within your Go binaries, making them self-contained and easier to distribute. By understanding the concepts and examples provided, you can effectively leverage this feature in your projects.

---

## Embedding Hidden Files in Go Programs

This section explains how to include hidden files within your Go binaries using `go:embed`.

**Default Behavior:**

- Files and directories starting with `.` or `_` (hidden) are excluded by default.

**Overriding the Default:**

1. **Including Hidden Files in Root Directory:**

   - Use `/*` after the directory name in the `go:embed` comment.
     - Example: `//go:embed parent_dir/*`
   - This includes hidden files within the specified directory, but not in subdirectories.

2. **Including All Hidden Files Recursively:**

   - Use `all:` before the directory name in the `go:embed` comment.
     - Example: `//go:embed all:parent_dir`
   - This includes hidden files in all subdirectories of the specified directory.

**Example:**

```go
//go:embed parent_dir
var noHidden embed.FS
//go:embed parent_dir/*
var parentHiddenOnly embed.FS
//go:embed all:parent_dir
var allHidden embed.FS

func main() {
  checkForHidden("noHidden", noHidden)
  checkForHidden("parentHiddenOnly", parentHiddenOnly)
  checkForHidden("allHidden", allHidden)
}

func checkForHidden(name string, dir embed.FS) {
  fmt.Println(name)
  allFileNames := []string{
    "parent_dir/.hidden",
    "parent_dir/child_dir/.hidden",
  }
  for _, v := range allFileNames {
    _, err := dir.Open(v)
    if err == nil {
      fmt.Println(v, "found")
    }
  }
  fmt.Println()
}
```

**Output:**

```
noHidden

parentHiddenOnly
parent_dir/.hidden found

allHidden
parent_dir/.hidden found
parent_dir/child_dir/.hidden found
```

**Explanation:**

- `noHidden` doesn't include hidden files.
- `parentHiddenOnly` includes only the `.hidden` file in the root directory (`parent_dir`).
- `allHidden` includes all hidden files, both in the root directory and subdirectory (`child_dir`).

**Remember:**

- Use these techniques with caution, as hidden files might contain sensitive information.
- Consider the security implications of including hidden files in your binary.

---

## Using go generate for Code Generation

This section explains the `go generate` tool and its role in automating code generation tasks.

**Key Points:**

- `go generate` itself doesn't perform actions; it executes commands based on comments.
- It's commonly used to generate code from specifications or existing code.

**Example: Generating Code from Protocol Buffers**

1. **Protobuf Schema:**

   ```protobuf
   syntax = "proto3";
   message Person {
       string name = 1;
       int32 id = 2;
       string email = 3;
   }
   ```

2. **go:generate Comment:**

   ```go
   //go:generate protoc -I=. --go_out=. --go_opt=module=github.com/learning-go-book-2e/proto_generate --go_opt=Mperson.proto=github.com/learning-go-book-2e/proto_generate/data person.proto
   ```

   - This comment instructs `go generate` to run `protoc` with specific arguments:
     - Generate Go code from `person.proto`.
     - Place the generated code in the current directory.
     - Set the module path and message name for the generated code.

3. **Running go generate:**

   ```bash
   go generate ./...
   ```

   - This creates a new directory (`data`) containing `person.pb.go`.
   - The generated code includes the `Person` struct, methods for marshalling/unmarshalling protobuf data.

**Other Use Cases:**

- Generating string representations for enumeration values using `stringer`.

**Benefits:**

- Automates tedious code generation tasks.
- Improves code consistency and maintainability.

**Additional Notes:**

- Requires installing external tools like `protoc`.
- Consider tool-specific configuration options for customization.

**Conclusion:**

`go generate` is a powerful tool for leveraging code generation tools within your Go projects. By understanding its functionality and common use cases, you can streamline your development process and reduce manual coding effort.

---

## Integrating go generate and Makefiles

This section discusses the interplay between `go generate` and Makefiles, providing guidance on their effective use.

**Separation of Responsibilities:**

- **go generate:** Generate source code based on specifications or existing code.
- **Makefile:** Define build and validation rules.

**Best Practices:**

- **Commit Generated Code:** Include generated code in your version control to ensure transparency and avoid external tool dependencies.
- **Automate go generate:** Add a `go generate` step to your Makefile as a dependency of your build step to ensure consistency.

**Exceptions:**

- **Minor Differences:** If `go generate` produces minor differences (e.g., timestamps) on the same input, consider excluding generated files from version control.
- **Long Build Times:** If `go generate` significantly slows down builds, evaluate whether the benefits outweigh the cost in developer productivity.

**Key Considerations:**

- **Manual vs. Automated:** While manual execution is possible, automation is generally preferred for consistency and to avoid oversight.
- **Version Control:** Committing generated code helps maintain a complete project history.
- **Performance:** Balance the benefits of code generation with potential build time impacts.

**Conclusion:**

By effectively integrating `go generate` and Makefiles, you can streamline your development process, automate code generation tasks, and ensure consistent project builds. Carefully consider the trade-offs between automation and manual intervention based on your specific project requirements and team practices.

---

## Understanding Build Information in Go Binaries

This section explores how Go binaries embed information about their build process for tracking purposes.

**Motivation:**

- Companies need to track software versions deployed in their environments.
- Manual tracking can be error-prone and impractical.

**Benefits of Go's Approach:**

- Go binaries automatically embed build information during compilation.
- This information includes:
  - Go version used.
  - Module versions.
  - Build commands used.
  - Version control system and revision details.

**Viewing Build Information:**

- Use the `go version -m` command to view embedded build information for a binary.

**Example:**

```bash
$ go build vulnerable   # Build the "vulnerable" program
$ go version -m vulnerable
vulnerable: go1.20
 path github.com/learning-go-book-2e/vulnerable (devel)
 dep gopkg.in/yaml.v2 v2.2.7 h1:VUgggvou5XRW9mHwD/yXxIYSMtY0zoKQf/v...
 build -compiler=gc
 build CGO_ENABLED=1
 ... (other build information)
```

**Applications:**

- **Vulnerability Scanning:** Tools like `govulncheck` can scan binaries for known vulnerabilities based on embedded module versions.
- **Improved Deployment Tracking:** Companies can track deployed software versions accurately.

**Additional Notes:**

- Govulncheck provides vulnerability information but may not pinpoint exact code lines.
- Use `go version -m` to identify the exact deployed version and then analyze the source code for fixes.

**Reading Build Information Programmatically:**

- Explore the `debug/buildinfo` package in the Go standard library for programmatic access to embedded build information.

**Conclusion:**

Go's embedded build information feature simplifies software tracking and vulnerability management, empowering companies to maintain efficient and secure deployments. By understanding how to access and utilize this information, you can enhance the reliability and security of your Go applications.

---

## Cross-Compiling Go Binaries

This section explains how to use `go build` to create Go binaries for different operating systems and CPU architectures.

**Cross-Compilation:**

- Go allows you to build binaries for platforms other than your development machine.
- This is achieved by setting the `GOOS` and `GOARCH` environment variables.

**Understanding `GOOS` and `GOARCH`:**

- `GOOS`: Specifies the target operating system (e.g., `linux`, `darwin`, `windows`).
- `GOARCH`: Specifies the target CPU architecture (e.g., `amd64`, `arm64`).

**Cross-Compiling Example:**

```bash
# Build for Linux on x86-64 (Intel)
GOOS=linux GOARCH=amd64 go build
```

**Additional Notes:**

- The valid values for `GOOS` and `GOARCH` can be found in the Go installation documentation.
- Some values may require specific translations (e.g., `darwin` for macOS).

**Conclusion:**

Cross-compilation in Go provides flexibility for deploying your applications to various platforms. By understanding and utilizing the `GOOS` and `GOARCH` environment variables, you can efficiently create binaries tailored to your target environments.

---

## Using Build Tags for Conditional Compilation

This section explains how to use build tags to conditionally include or exclude code based on platform-specific conditions.

**Conditional Compilation:**

- Enables writing code that is specific to certain operating systems, CPU architectures, or Go versions.
- Helps maintain backward compatibility and target specific platforms.

**Methods:**

1. **Filename Conventions:**

   - Prefix filenames with `GOOS_GOARCH` (e.g., `something_windows_arm64.go`).
   - This method is less flexible compared to build tags.

2. **Build Tags:**
   - Use `//go:build` comments before the package declaration.
   - Specify conditions using boolean operators (e.g., `&&`, `||`, `!`).

**Example:**

```go
//go:build !darwin && !linux || (darwin && !go1.12)

package mypackage

// Code specific to platforms other than macOS or Linux, or macOS versions older than 1.12
```

**Built-in Build Tags:**

- `unix`: Matches any Unix-like platform.
- `cgo`: Matches if cgo is supported and enabled.
- `go1.X`: Matches Go versions greater than or equal to 1.X.

**Custom Build Tags:**

- Define your own tags (e.g., `gopher`).
- Control compilation using the `-tags` flag with `go build`, `go run`, or `go test`.

**Example:**

```go
//go:build gopher

package mypackage

// Code that is only included when using the "gopher" build tag
```

**Best Practices:**

- Use build tags for more complex conditions and platform-specific code.
- Consider filename conventions for simple cases.
- Avoid whitespace between `//` and `go:build`.

**Conclusion:**

Build tags provide a powerful mechanism for conditional compilation in Go. By effectively using them, you can tailor your code to specific platforms, maintain compatibility, and manage experimental or unfinished code segments.

---

## Testing Different Go Versions

This section explores methods for verifying compatibility of your Go programs across various Go versions.

**Challenges of Bugs and Compatibility:**

- Despite backward compatibility efforts, bugs may arise in new Go releases.
- Ensuring compatibility with older versions becomes necessary.

**Testing Options:**

1. **Secondary Go Environments:**

   - Install specific Go versions using `golang.org/dl/goX.Y.Z@latest` (replace X, Y, Z with version numbers) and download the corresponding Go tools.
   - Use `goX.Y.Z` command (e.g., `go1.19.2`) to build and test your code with that specific version.
   - Uninstall the secondary environment after testing by removing the directories from `~/sdk` and `~/go/bin`.

**Example:**

```bash
$ go install golang.org/dl/go1.19.2@latest
$ go1.19.2 download
$ go1.19.2 build myproject  # Build your project
```

2. **Containerized Testing Environments:**

   - Consider using containerized environments (like Docker) to manage multiple Go versions without manual installation.

**Learning More with `go help`**

- Explore the `go help` command for comprehensive documentation on Go tools, runtime, and other features.
- Get detailed information on specific topics like import paths, modules, and working with non-public source code.

**Example:**

```bash
$ go help importpath  # Get help on import path syntax
```

**Conclusion:**

Effective testing across Go versions helps ensure compatibility and reduces the risk of regressions when using newer versions. By leveraging secondary environments, containerization, and the `go help` command, you can streamline your testing process and maintain robust Go applications.

---

## Exercises

Here’s a step-by-step solution to the exercises:

### **Exercise 1: Embedding UDHR Text Files**

1. **Step 1: Download Text Files**

   - Go to the [UN’s Universal Declaration of Human Rights (UDHR)](https://www.un.org/en/about-us/universal-declaration-of-human-rights) page.
   - Copy the English version of the text into a file named `english_rights.txt`.
   - Click on the "Other Languages" link and copy a few other language versions into files named according to the language, e.g., `french_rights.txt`, `spanish_rights.txt`, etc.

2. **Step 2: Create the Program**

   - Below is the Go program that reads the embedded files and prints the UDHR in the specified language.

   ```go
   // main.go
   package main

   import (
       "embed"
       "flag"
       "fmt"
       "log"
   )

   //go:embed *.txt
   var files embed.FS

   func main() {
       language := flag.String("language", "english", "The language of the UDHR document")
       flag.Parse()

       fileName := *language + "_rights.txt"
       data, err := files.ReadFile(fileName)
       if err != nil {
           log.Fatalf("Failed to read file: %v", err)
       }

       fmt.Println(string(data))
   }
   ```

3. **Step 3: Run the Program**
   - Save the code above as `main.go`.
   - Run the program using the command:
     ```bash
     go run main.go --language=english
     ```
   - You can replace `english` with any other available language like `french`, `spanish`, etc.

### **Exercise 2: Use `staticcheck`**

1. **Install `staticcheck`**

   - Run the following command to install `staticcheck`:
     ```bash
     go install honnef.co/go/tools/cmd/staticcheck@latest
     ```

2. **Run `staticcheck` on the Program**
   - Run the static analysis tool against your program:
     ```bash
     staticcheck ./...
     ```
   - Fix any issues reported by `staticcheck`. Most issues are usually about unused variables, imports, or common best practices.

### **Exercise 3: Cross-Compile the Program**

1. **Cross-Compile for ARM64 on Windows**

   - Use the following command to cross-compile the program for ARM64 on Windows:
     ```bash
     GOOS=windows GOARCH=arm64 go build -o udhr_windows_arm64.exe main.go
     ```

2. **Cross-Compile for AMD64 on Linux (if on ARM64 Windows)**
   - If you are on an ARM64 Windows computer, cross-compile for AMD64 on Linux using:
     ```bash
     GOOS=linux GOARCH=amd64 go build -o udhr_linux_amd64 main.go
     ```

These steps will help you create the program, ensure it adheres to best practices using `staticcheck`, and cross-compile it for different platforms.
