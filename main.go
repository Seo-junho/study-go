package main

import (
	"awesomeProject/mydict"
	"fmt"
)

func main() {
	dictionary := mydict.Dictionary{"first" : "First word"}
	baseWord := "hello"
	dictionary.Add(baseWord, "First")
	dictionary.Search(baseWord)
	err := dictionary.Delete(baseWord)
	word, err2 := dictionary.Search(baseWord)
	if err == nil{
		fmt.Println(err2)
	}
	fmt.Println(word)

}