package main

import (
	"fmt"

	"github.com/mmdriley/travisci"
)

func main() {
	fmt.Printf("%+v\n", travisci.NewClient())
}
