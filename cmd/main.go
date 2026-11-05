package main

import (
	"fmt"
	"portscan"
	"time"
)

func main() {
	scanner, err := portscan.NewScanner(100, 500*time.Millisecond)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(scanner.GetSettings())

	fmt.Println(scanner.SetConcurrency(0))
	fmt.Println(scanner.SetTimeout(0))


	scanner.SetConcurrency(200)
	scanner.SetTimeout(600*time.Millisecond)
	fmt.Println(scanner.GetSettings())
}