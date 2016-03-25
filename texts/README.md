# texts

Go package for truncating texts.

# Installation

```go
go get github.com/ChristianSiegert/texts
```

# Usage

```go
import "github.com/ChristianSiegert/texts"

texts.Truncate("Hello world", 10, "…", false)
// Result: "Hello …"

texts.Truncate("Hello world", 10, "…", true)
// Result: "Hello wor…"
```

[See documentation](https://godoc.org/github.com/ChristianSiegert/texts).
