package main

import (
	"fmt"
)

func main() {
	bar := &Bar{
		&Foo{},
	}

	bar.Hello()
}

type Foo struct {
}

func (f *Foo) Hello() {
	fmt.Println("Foo -> Hello")
	f.FooBin()
}
func (f *Foo) FooBin() {
	fmt.Println("Foo -> FooBin")
}

type Bar struct {
	*Foo
}

func (f *Bar) Hello() {
	fmt.Println("Bar -> Hello")
	f.FooBin()
}

//func (f *Bar) FooBin() {
//	fmt.Println("Bar -> FooBin")
//}
