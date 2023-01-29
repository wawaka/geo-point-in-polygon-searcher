package searcher

import (
	"math"
)

// Line defines a line segment in 2D space
type Line struct {
	P1 Point
	P2 Point
}

// Point defines a point in 2D space
type Point struct {
	X float64
	Y float64
}

// Intersect checks if two lines are intersecting
func Intersect(l1 Line, l2 Line) bool {
	a1 := l1.P2.Y - l1.P1.Y
	b1 := l1.P1.X - l1.P2.X
	c1 := a1*l1.P1.X + b1*l1.P1.Y

	a2 := l2.P2.Y - l2.P1.Y
	b2 := l2.P1.X - l2.P2.X
	c2 := a2*l2.P1.X + b2*l2.P1.Y

	determinant := a1*b2 - a2*b1
	if determinant == 0 {
		return false
	}

	x := (b2*c1 - b1*c2) / determinant
	y := (a1*c2 - a2*c1) / determinant

	if math.Min(l1.P1.X, l1.P2.X) <= x && x <= math.Max(l1.P1.X, l1.P2.X) &&
		math.Min(l1.P1.Y, l1.P2.Y) <= y && y <= math.Max(l1.P1.Y, l1.P2.Y) &&
		math.Min(l2.P1.X, l2.P2.X) <= x && x <= math.Max(l2.P1.X, l2.P2.X) &&
		math.Min(l2.P1.Y, l2.P2.Y) <= y && y <= math.Max(l2.P1.Y, l2.P2.Y) {
		return true
	}

	return false
}
