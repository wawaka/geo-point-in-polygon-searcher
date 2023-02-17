package searcher

import (
	"math/rand"
	"time"

	"github.com/engelsjk/polygol"
	"github.com/golang/geo/s2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GeneratePoint() []float64 {
	return []float64{
		rand.Float64(),
		rand.Float64(),
	}
}

func RandRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func GenerateTriangle() [][]float64 {
	return [][]float64{
		GeneratePoint(),
		GeneratePoint(),
		GeneratePoint(),
	}
}

func GeneratePolygon(n int) [][]float64 {
	multipolygon := [][][][]float64{{GenerateTriangle()}}
	// fmt.Printf("%d polygons, %d rings, %d points\n", len(multipolygon), len(multipolygon[0]), len(multipolygon[0][0]))
	for i := 0; i < 10; i++ {
		// triangle := GenerateTriangle()
		var err error
		multipolygon, err = polygol.Union(multipolygon, [][][][]float64{{GenerateTriangle()}})
		if err != nil {
			panic(err)
		}
		for _, polygon := range multipolygon {
			for _, ring := range polygon {
				if len(ring) > n {
					return ring[:len(ring)-1]
				}
			}
		}
		// fmt.Printf("%d polygons, %d rings, %d points\n", len(multipolygon), len(multipolygon[0]), len(multipolygon[0][0]))
	}
	ring := multipolygon[0][0]
	// fmt.Printf("%d %v\n", len(polygon), polygon)
	ring = ring[:len(ring)-1]
	// fmt.Printf("%d %v\n", len(polygon), polygon)
	return ring
}

func S2Point(ll []float64) s2.Point {
	s2ll := s2.LatLngFromDegrees(ll[0], ll[1])
	s2point := s2.PointFromLatLng(s2ll)
	return s2point
}

func S2Loop(ring [][]float64) *s2.Loop {
	var s2points []s2.Point
	for _, ll := range ring {
		s2point := S2Point(ll)
		s2points = append(s2points, s2point)
	}
	s2loop := s2.LoopFromPoints(s2points)
	if s2loop.TurningAngle() < 0 {
		s2loop.Invert()
	}
	return s2loop
}

func S2Polygon(ring [][]float64) *s2.Polygon {
	s2loop := S2Loop(ring)
	var s2loops []*s2.Loop
	s2loops = append(s2loops, s2loop)
	s2polygon := s2.PolygonFromLoops(s2loops)
	return s2polygon
}
