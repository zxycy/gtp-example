package main

import "gtp-example/gtpv1_test/gtp"

func main() {
	gtp.Run()
	select {}
}
