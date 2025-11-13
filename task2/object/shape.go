package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func main() {
	rectangle := &Rectangle{width: 10, height: 5}
	circle := &Circle{radius: 5}
	fmt.Println("长方形的面积:", rectangle.Area(), "长方形的周长:", rectangle.Perimeter())
	fmt.Println("圆形的周长:", circle.Perimeter(), "圆形的面积:", circle.Area())

}
