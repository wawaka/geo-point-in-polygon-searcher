package searcher

import (
	"fmt"
	"testing"

	geo "github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/assert"
)

const testContainsPolygonsMax = 1000
const testContainsPointsMax = 1000000

var testContainsPolygons [][][]float64
var testContainsPoints [][]float64
var testContainsConvertedPolygons []*geo.Polygon
var testContainsConvertedPoints []*geo.Point

func init() {
	for i := 0; i < testContainsPolygonsMax; i++ {
		polygon := GeneratePolygon(10)
		// if !reflect.DeepEqual(polygon[0], polygon[len(polygon)-1]) {
		// 	panic("oops")
		// }
		// polygon = polygon[0 : len(polygon)-1]
		// if reflect.DeepEqual(polygon[0], polygon[len(polygon)-1]) {
		// 	panic("oops2")
		// }
		testContainsPolygons = append(testContainsPolygons, polygon)
		testContainsConvertedPolygons = append(testContainsConvertedPolygons, ConvertPolygon(polygon))
	}
	for i := 0; i < testContainsPointsMax; i++ {
		point := GeneratePoint()
		testContainsPoints = append(testContainsPoints, point)
		testContainsConvertedPoints = append(testContainsConvertedPoints, ConvertPoint(point))
	}
}

func ConvertPolygon(polygon [][]float64) *geo.Polygon {
	converted_points := []*geo.Point{}
	for _, point := range polygon {
		converted_points = append(converted_points, ConvertPoint(point))
	}
	return geo.NewPolygon(converted_points)
}

func ConvertPoint(point []float64) *geo.Point {
	return geo.NewPoint(point[0], point[1])
}

func TestPolygonContainsEquivalence(t *testing.T) {
	for pol_idx := 0; pol_idx < 1000; pol_idx++ {
		for pt_idx := 0; pt_idx < 1000; pt_idx++ {
			contains := PolygonContains(testContainsPolygons[pol_idx%testContainsPolygonsMax], testContainsPoints[pt_idx%testContainsPointsMax])
			converted_contains := testContainsConvertedPolygons[pol_idx%testContainsPolygonsMax].Contains(testContainsConvertedPoints[pt_idx%testContainsPointsMax])
			assert.Equal(t, converted_contains, contains)
		}
	}
}

func TestPolygonContainsEquivalenceS2(t *testing.T) {
	C := map[string]int{}
	for pol_idx := 0; pol_idx < 1000; pol_idx++ {
		s2polygon := S2Polygon(testContainsPolygons[pol_idx%testContainsPolygonsMax])
		for pt_idx := 0; pt_idx < 1000; pt_idx++ {
			s2point := S2Point(testContainsPoints[pt_idx%testContainsPointsMax])
			contains := PolygonContains(testContainsPolygons[pol_idx%testContainsPolygonsMax], testContainsPoints[pt_idx%testContainsPointsMax])
			s2_contains := s2polygon.ContainsPoint(s2point)
			if contains == s2_contains {
				C["match"]++
				// fmt.Println(testContainsPoints[pt_idx%testContainsPointsMax])
			} else {
				C[fmt.Sprintf("mismatch-%v/%v", contains, s2_contains)]++
				// jsonPoints, _ := json.Marshal(testContainsPolygons[pol_idx%testContainsPolygonsMax])
				// jsonPoint, _ := json.Marshal(testContainsPoints[pt_idx%testContainsPointsMax])
				// fmt.Printf("var points = %v\n", string(jsonPoints))
				// fmt.Printf("var point = %v\n", string(jsonPoint))
				// fmt.Printf("own: %v, s2: %v\n", contains, s2_contains)
				// fmt.Println()
				// return
			}
			// assert.Equal(t, s2_contains, contains, "%v", testContainsPoints[pt_idx%testContainsPointsMax])
		}
	}
	// for k, v := range C {
	// 	fmt.Printf("[%v]=%v\n", k, v)
	// }
}

// func TestS2Basic(t *testing.T) {
// 	// s2polygon := S2Polygon([][]float64{{0, 0}, {0.5, 0}, {0.5, 0.5}, {0, 0.5}})
// 	s2polygon := S2Polygon([][]float64{{0, 0}, {0, 0.5}, {0.5, 0.5}, {0.5, 0}})
// 	// fmt.Printf(s2polygon.)
// 	for i := 0; i < 10; i++ {
// 		point := GeneratePoint()
// 		s2point := S2Point(point)
// 		s2contains := s2polygon.ContainsPoint(s2point)
// 		// fmt.Printf("%v %v\n", s2contains, point)
// 	}
// }

func BenchmarkPolygonContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for pol_idx := 0; pol_idx < 10; pol_idx++ {
			PolygonContains(testContainsPolygons[pol_idx%testContainsPolygonsMax], testContainsPoints[i%testContainsPointsMax])
		}
	}
}

func BenchmarkUpstreamPolygonContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for pol_idx := 0; pol_idx < 10; pol_idx++ {
			testContainsConvertedPolygons[pol_idx%testContainsPolygonsMax].Contains(testContainsConvertedPoints[i%testContainsPointsMax])
		}
	}
}
