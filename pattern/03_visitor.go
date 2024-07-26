package main

import "fmt"

type Shape interface {
	Accept(v Visitor)
}

type Visitor interface {
	VisitCircle(c *Circle)
	VisitRectangle(r *Rectangle)
}

type Circle struct {
	Radius float64
}

func (c *Circle) Accept(v Visitor) {
	v.VisitCircle(c)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r *Rectangle) Accept(v Visitor) {
	v.VisitRectangle(r)
}

type AreaVisitor struct {
	Area float64
}

func (a *AreaVisitor) VisitCircle(c *Circle) {
	a.Area = 3.14 * c.Radius * c.Radius
}

func (a *AreaVisitor) VisitRectangle(r *Rectangle) {
	a.Area = r.Width * r.Height
}

type PerimeterVisitor struct {
	Perimeter float64
}

func (p *PerimeterVisitor) VisitCircle(c *Circle) {
	p.Perimeter = 2 * 3.14 * c.Radius
}

func (p *PerimeterVisitor) VisitRectangle(r *Rectangle) {
	p.Perimeter = 2 * (r.Width + r.Height)
}

func main() {
	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 3, Height: 4},
	}

	areaVisitor := &AreaVisitor{}
	perimeterVisitor := &PerimeterVisitor{}

	for _, shape := range shapes {
		shape.Accept(areaVisitor)
		fmt.Printf("Area: %.2f\n", areaVisitor.Area)

		shape.Accept(perimeterVisitor)
		fmt.Printf("Perimeter: %.2f\n", perimeterVisitor.Perimeter)
	}
}
