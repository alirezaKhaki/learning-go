package main

import "fmt"

type Integer interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

func PlusOneThousand[T Integer](in T, out T) T {
	return in + out
}

func main() {
	var test int32 = 100
	a := PlusOneThousand(test, 100000000)
	fmt.Println(a)
}
