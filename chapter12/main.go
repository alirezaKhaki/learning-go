package main

import (
	"fmt"
)

func putDataOnChannel(ch *chan int, value int) {
	defer close(*ch)
	*ch <- value
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go putDataOnChannel(&ch1, 1)
	go putDataOnChannel(&ch2, 2)
	go putDataOnChannel(&ch3, 3)

	for {
		select {
		case data := <-ch1:
			fmt.Println(data)
		case data := <-ch2:
			fmt.Println(data)
		case data := <-ch3:
			fmt.Println(data)
		default:
			return
		}

	}

	fmt.Println("here")
}
