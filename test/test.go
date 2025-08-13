package main

import "fmt"

type Person struct {
	Name string
	Age  string
}

func main() {
	person := Person{Name: "Vistar", Age: "23"}
	fmt.Println(person.Name)
	fmt.Println(person.Age)
}
