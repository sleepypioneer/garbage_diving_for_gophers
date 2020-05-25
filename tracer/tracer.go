package tracer

import (
	"log"
	"os"
	"runtime/trace"
)

// WithTrace runs the program passed with the addition of tracing
func WithTrace(program func(int, func()), paramInt int, paramFunc func()) {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	program(paramInt, paramFunc)
}
