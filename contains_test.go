package searcher

import (
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
