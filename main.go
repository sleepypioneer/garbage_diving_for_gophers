package main

import (
	"fmt"
	"garbage_diving/tracer"
	_ "garbage_diving/tracer"
)

func repeatXTimes(x int, functionToRepeat func()) {
	for i := 0; i < x; i++ {
		functionToRepeat()
	}
}

func myFunc() {
	a := make([]int, 1000000)
	fmt.Println(len(a))
}

func main() {
	repeatXTimes(10, myFunc)
	// Uncomment line below to run with tracer
	tracer.WithTrace(repeatXTimes, 10, myFunc)
}
