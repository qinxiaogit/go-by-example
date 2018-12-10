package main

import (
	"fmt"
	"math"
)

type geo interface {
	area() float64
	circ() float64
}

type rect struct {
	width, height float64
}
type circle struct {
	radius float64
}

func (r rect)area() float64{
	return r.width*r.height
}
func (r rect)circ()float64{
	return 2*(r.height+r.width)
}

func (c circle)area() float64{
	return math.Pi*math.Pow(c.radius,2)
}

func (c circle)circ()float64{
	return math.Pi*2*c.radius
}

func measure(g geo){
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.circ())
}

func main() {
	measure(rect{1,2})
	measure(circle{2})
}
