# Nanux

Nanux is a toolkit for deploying microservice.

*Work in progress. Breaking change might happen. However nanux is already used in production. *

## Getting started

1. Install nanux into your project: `go get github.com/nanux-io/nanux`
2. Choose the transporter you which to use (now transporters that are officially
supported are:
    * [nats](https://nats.io) : <https://github.com/nanux-io/tnats>
    * http : <https://github.com/nanux-io/thttp>

Example of code using nats transporter:

```go
package main

import (
  "github.com/nanux-io/nanux"
  "github.com/nanux-io/tnats"
)

func main() {
  // Instantiate nanux
  t := tnats.New("natsURL", []tnats.Option{})
  n := nanux.New(&t, nil)

  // Add error handler which will be called when an error an action return an error.
  // the error handler must be of type `handler.ManageError`
  errorHandler := func(ctx *interface{}, req nanux.Request) ([]byte, error) {
    // Process the error here.
    return nil, nil
  }
  n.HandleError(errorHandler)
  
  // Add handler which will be called when "my.route" is reached
  handlerForMyRoute := nanux.Handler{
    Fn: func(ctx *interface{}, req nanux.Request) ([]byte, error) {
        return nil, nil
    },
  }

  n.Handle("my.route", handlerForMyRoute)

  // defer the closing of the connection to the transporter
  defer n.Close()

  // Run nanux
  if err := n.Run(); err != nil {
    log.Fatalf("Problem when starting listening incoming requests - %s\n", err)
  }
}

```

## Development

### Installation

1. Launch the docker container : `cd build && docker-compose up`
2. Enter in the container: `docker-compose exec go bash`

### Test

To launch test with generation of html coverage page execut :
`go test -coverprofile=coverage.out &&  go tool cover -html=coverage.out -o coverage.html`

### Editor

The project used the go modules, in order to have tools working in your
editor don't forget to add the container `GOPATH` to the host `GOPATH` when
working on the project.

#### VS Code

For VS Code, add the followin setting for the workspace:

```json
{
  "go.gopath": "${env:GOPATH}:${workspaceFolder}/build/.go"`
}
```

## Contributor

Thanks to [Nicolas Talle](https://github.com/nicolab) for the feedback.
