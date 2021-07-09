# GoRcon
Simple rcon-client written in Golang using [Valve RCON protocol documentation](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol).

## Installation
```text
go get github.com/xStrixU/gorcon
```

## Example code
```go
package main

import (
	"fmt"
	"github.com/xStrixU/gorcon"
	"log"
)

func main() {
	rcon, err := gorcon.Connect("localhost:19132", "password")

	if err != nil {
		log.Fatal(err)
	}

	defer rcon.Close()

	response, err := rcon.SendCommand("say Hello Golang RCON!")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response message: " + response)
}
```

## License
MIT License, see [LICENSE](LICENSE)