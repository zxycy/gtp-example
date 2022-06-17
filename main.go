package main

import (
	"gtp-example/server"
)

func main() {

	go server.GTPServer()
	select {}
}
