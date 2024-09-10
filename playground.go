package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 100)
	ch <- 10
	for i := 0; i < 100; i++ {
		defer close(ch)
		go func() {
			ch <- i
		}()
		time.Sleep(time.Second * 1)
	}

	for value := range ch {
		fmt.Println(value)
	}

	fmt.Println("...done")
}
