package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)
	people := [2]string{"nico", "juno"}
	for _, person := range people{
		go isSexy(person, c)
	}
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}

func isSexy(person string, channel chan bool) {
	time.Sleep(time.Second * 5)
	fmt.Println(person)
	channel <- true
}