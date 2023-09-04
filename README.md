[![Build Status](https://github.com/mlavergn/godl/workflows/CI/badge.svg?branch=master)](https://github.com/mlavergn/godl/actions)
[![Go Report](https://goreportcard.com/badge/github.com/mlavergn/godl)](https://goreportcard.com/report/github.com/mlavergn/godl)

# GoDL

Go project implementing a very low-overhead and fast HTTP file server.

NOTE: This project is now "standard" functionality in Go and can be recreated as:

```go
package main

import "net/http"

func main() {
	http.Handle("/", http.FileServer(http.Dir("files")))
}
```

