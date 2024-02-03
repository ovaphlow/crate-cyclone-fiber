package main

import (
	"ovaphlow/cratecyclone/utilities"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	utilities.InitSlog()

	Serve("localhost:8421")
}
