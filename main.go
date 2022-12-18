package main

import (
	"fmt"
	"os"
	"time"

	"github.com/diyliv/ransomware/pkg/reading"
)

func main() {
	start := time.Now()

	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	reading.Process(f)
	fmt.Printf("The whole operation took: %v\n", time.Since(start))
}
