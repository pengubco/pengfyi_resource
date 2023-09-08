package main

import "fmt"

func main() {
	foo()
}

func foo() {
	c := Cat{name: "Tomcat"}
	defer c.Hi()
	defer c.Bye()
	c.name = "Garfield"
}

type Cat struct {
	name string
}

func (c Cat) Hi() {
	fmt.Printf("Hello %s\n", c.name)
}

func (c *Cat) Bye() {
	fmt.Printf("Bye %s\n", c.name)
}
