package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf"
)

var build = "develop"

func main() {
	fmt.Println(conf.String("this is the main function"))
	fmt.Println("-----Starting Service : ", build, " -----")
	defer fmt.Println("-----Ended Service-----")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	defer fmt.Println("-----Stopping Service-----")

}
