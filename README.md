# Wavix Go SDK

Wavix Go SDK provides convenient and easy-to-understand functions for interacting with the Wavix API. You can use this package to integrate your Go application or service with Wavix functionality.

The current list of Wavix HTTP REST API can be found here: https://wavix.com/api

## Installation and Usage

### Installation

```sh
go get github.com/wavix/sdk-go
```

### Import

```go
import "github.com/wavix/sdk-go"
```

### Typical Usage Example

```go
package main

import (
    wavix "github.com/wavix/sdk-go"
)

func main() {
    instance := wavix.Init(wavix.ClientOptions{Appid: "<YOUR APPID>"})

    from := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
    to := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

    cdrList, err := instance.Cdr.GetCdrList(wavix.GetCdrListQueryParams{
        Type:   "placed",
        RequiredDateParams: utils.RequiredDateParams{
            From: utils.QueryDateParams(from),
            To:   utils.QueryDateParams(to),
            }, PaginationParams: utils.PaginationParams{Page: 1, PerPage: 5}
    })

    if err != nil {
        panic(err.Message)
    }

    ...
}
```

## License

This SDK is distributed under the MIT License. See [LICENSE](./LICENSE) for more detailed information.
