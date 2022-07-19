package main

import "gtp-example/gtpv2_test/gtp"

func main() {
	gtp.Run(gtp.Dispatch)
	select {}
}
