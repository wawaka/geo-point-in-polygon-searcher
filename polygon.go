package searcher

import (
	"math/rand"
	"time"

	"github.com/engelsjk/polygol"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GeneratePoint() []float64 {
	return []float64{rand.Float64(), rand.Float64()}
}

func GenerateTriangle() [][]float64 {
	polygon := [][]float64{}
	for i := 0; i < 3; i++ {
		point := GeneratePoint()
		polygon = append(polygon, point)
	}
	return polygon
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
				if len(ring) >= n {
					return ring
				}
			}
		}
		// fmt.Printf("%d polygons, %d rings, %d points\n", len(multipolygon), len(multipolygon[0]), len(multipolygon[0][0]))
	}
	return multipolygon[0][0]
}
