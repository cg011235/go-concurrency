package main

import (
	"fmt"
	"time"
)

// Run spinner till we calculate 44th number in fibonacci series

func fibo(N uint64) uint64 {
	if N < 2 {
		return N
	}
	return fibo(N - 1) + fibo(N - 2)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func main() {
	go spinner(100 * time.Millisecond)
	const n = 44
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibo(n))
}


//Build with: go build
