package main

import (
	"go-pocs/manyProdOneCons"
	"runtime"
	"strconv"
)

func main() {
	println("Number of processors: " + strconv.Itoa(runtime.NumCPU()))
	runtime.GOMAXPROCS(runtime.NumCPU())
	manyProdOneCons.ManyProducersOneConsumer()
}
