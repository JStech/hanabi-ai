package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 50; i++ {
		fmt.Println(Card(i).Color(), Card(i).Number())
	}
}
