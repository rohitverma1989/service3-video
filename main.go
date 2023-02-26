package main

import (
	"fmt"

	"github.com/ardanlabs/conf"
)

func main() {
	fmt.Println(conf.String("this is the main function"))
}
