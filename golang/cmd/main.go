package main

import (
	"fmt"
)

func main() {
	a := []int(nil)
	for i := range a {
		fmt.Println(i)
	}
	fmt.Println(len(a))

}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
