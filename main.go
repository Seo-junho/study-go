package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan string)
	people := [5]string{"nico", "juno", "dal", "japanguy", "larry"}
	for _, person := range people{
		go isSexy(person, c)
	}
	for i:=0; i <len(people); i++{
		fmt.Print("waiting for", i)
		fmt.Println(<-c)
	}
}

func isSexy(person string, channel chan string) {
	time.Sleep(time.Second * 5)
	channel <- person + " is sexy"
}