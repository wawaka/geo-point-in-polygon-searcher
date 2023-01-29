package searcher

import (
	"math"
)

func MultiPolygonContains(multipolygon [][][][]float64, point []float64) bool {
	for _, polygon_with_holes := range multipolygon {
		if PolygonWithHolesContains(polygon_with_holes, point) {
			return true
		}
	}
	return false
}

func PolygonWithHolesContains(polygon [][][]float64, point []float64) bool {
	if !PolygonContains(polygon[0], point) {
		return false
	}
	for _, hole := range polygon[1:] {
		if PolygonContains(hole, point) {
			return false
		}
	}
	return true
}

func PolygonContains(polygon [][]float64, p []float64) bool {
	contains := false
	for i := 0; i < len(polygon); i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%len(polygon)]
		if p[1] > math.Min(p1[1], p2[1]) {
			if p[1] <= math.Max(p1[1], p2[1]) {
				if p[0] < math.Max(p1[0], p2[0]) {
					if p1[1] != p2[1] {
						latIntersection := (p[1]-p1[1])*(p2[0]-p1[0])/(p2[1]-p1[1]) + p1[0]
						if p1[0] == p2[0] || p[0] <= latIntersection {
							contains = !contains
						}
					}
				}
			}
		}
	}
	return contains
}
