package searcher

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/h3-go/v4"
)

const BENCHMARK_POLYGONS_NUMBER = 100
const BENCHMARK_POINTS_NUMBER = 1000000

type IdPolygon struct {
	Id      uint32
	Polygon [][]float64
}

var idpolygons []IdPolygon
var points [][]float64

func init() {
	var err error

	jsonFile, err := os.Open("idpolygons.json")
	if err != nil {
		// fmt.Println(err)
		panic(err)
	} else {
		defer jsonFile.Close()

		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, &idpolygons)
		if err != nil {
			panic(err)
		}
		// fmt.Println(len(idpolygons))
	}

	var bboxes []BoundingBox
	for _, idpolygon := range idpolygons {
		bbox := BoundingBox{}
		for _, p := range idpolygon.Polygon {
			bbox.AddPoint(p)
		}
		bboxes = append(bboxes, bbox)
	}

	for i := 0; i < BENCHMARK_POINTS_NUMBER; i++ {
		idpolygoni := rand.Intn(len(idpolygons))
		bbox := bboxes[idpolygoni]
		min := bbox.Min()
		max := bbox.Max()
		p := []float64{
			RandRange(min[0], max[0]),
			RandRange(min[1], max[1]),
		}
		points = append(points, p)
	}
}

func TestShapeIndexSearcherCorrectness(t *testing.T) {
	return
	bbs := NewBoundingBoxOptimizedSearcher()
	sis := NewShapeIndexSearcher()
	for i := 0; i < 100; i++ {
		polygon := GeneratePolygon(10)
		s2loop := S2Loop(polygon)
		sis.AddShape(uint32(i), s2loop)
		bbs.AddMultiPolygon(uint32(i), [][][][]float64{{polygon}})
	}

	for i := 0; i < 1000; i++ {
		point := GeneratePoint()
		// s2point := S2Point(point)
		ids_sis := sis.SearchPoint(point)
		sort.Slice(ids_sis, func(i, j int) bool { return ids_sis[i] < ids_sis[j] })
		// fmt.Printf("sis %v\n", ids_sis)

		ids_bbs := bbs.SearchPoint(point)
		sort.Slice(ids_bbs, func(i, j int) bool { return ids_bbs[i] < ids_bbs[j] })
		// fmt.Printf("bbs %v\n", ids_bbs)

		// fmt.Println()

		assert.Equal(t, ids_bbs, ids_sis)
		// bbs.AddMultiPolygon(uint32(i), [][][][]float64{{polygon}})
	}
}

func BenchmarkShapeIndexSearcher(b *testing.B) {
	s := NewShapeIndexSearcher()
	for i := 0; i < BENCHMARK_POLYGONS_NUMBER; i++ {
		polygon := GeneratePolygon(10)
		s2loop := S2Loop(polygon)
		s.AddShape(uint32(i), s2loop)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := GeneratePoint()
		s.SearchPoint(p)
	}
}
func BenchmarkShapeIndexSearcher_real(b *testing.B) {
	s := NewShapeIndexSearcher()
	for _, idpolygon := range idpolygons {
		s2loop := S2Loop(idpolygon.Polygon)
		s.AddShape(uint32(idpolygon.Id), s2loop)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := points[i%len(points)]
		s.SearchPoint(p)
	}
}

func BenchmarkShapeIndexSearcher_real_empty(b *testing.B) {
	s := NewShapeIndexSearcher()
	for _, idpolygon := range idpolygons {
		s2loop := S2Loop(idpolygon.Polygon)
		s.AddShape(uint32(idpolygon.Id), s2loop)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SearchPoint([]float64{0, 0})
	}
}

func BenchmarkShapeIndexSearcher_real_query(b *testing.B) {
	s := NewShapeIndexSearcher()
	for _, idpolygon := range idpolygons {
		s2loop := S2Loop(idpolygon.Polygon)
		s.AddShape(uint32(idpolygon.Id), s2loop)
	}

	q := s.Query()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := points[i%len(points)]
		s.SearchPointQuery(q, p)
	}
}

func BenchmarkShapeIndexSearcher_real_query_empty(b *testing.B) {
	s := NewShapeIndexSearcher()
	for _, idpolygon := range idpolygons {
		s2loop := S2Loop(idpolygon.Polygon)
		s.AddShape(uint32(idpolygon.Id), s2loop)
	}

	q := s.Query()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SearchPointQuery(q, []float64{0, 0})
	}
}

func BenchmarkShapeIndexSearcher_real_query_empty_point(b *testing.B) {
	s := NewShapeIndexSearcher()
	for _, idpolygon := range idpolygons {
		s2loop := S2Loop(idpolygon.Polygon)
		s.AddShape(uint32(idpolygon.Id), s2loop)
	}

	q := s.Query()
	p := S2Point([]float64{0, 0})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SearchS2PointQuery(q, p)
	}
}

func BenchmarkShapeIndexSearcherQuery(b *testing.B) {
	s := NewShapeIndexSearcher()
	for i := 0; i < BENCHMARK_POLYGONS_NUMBER; i++ {
		polygon := GeneratePolygon(10)
		s2loop := S2Loop(polygon)
		s.AddShape(uint32(i), s2loop)
	}

	q := s.Query()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := GeneratePoint()
		s.SearchPointQuery(q, p)
	}
}

func BenchmarkBBoxSearcher(b *testing.B) {
	s := NewBoundingBoxOptimizedSearcher()
	for i := 0; i < BENCHMARK_POLYGONS_NUMBER; i++ {
		polygon := GeneratePolygon(10)
		s.AddMultiPolygon(uint32(i), [][][][]float64{{polygon}})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := GeneratePoint()
		s.SearchPoint(p)
	}
}

func BenchmarkBBoxSearcher_real(b *testing.B) {
	s := NewBoundingBoxOptimizedSearcher()
	for _, idpolygon := range idpolygons {
		s.AddRing(uint32(idpolygon.Id), idpolygon.Polygon)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := points[i%len(points)]
		s.SearchPoint(p)
	}
}

func BenchmarkBBoxSearcher_real_empty(b *testing.B) {
	s := NewBoundingBoxOptimizedSearcher()
	for _, idpolygon := range idpolygons {
		s.AddRing(uint32(idpolygon.Id), idpolygon.Polygon)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SearchPoint([]float64{0, 0})
	}
}

func empty() {

}
func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		empty()
	}
}

func BenchmarkS2Point(b *testing.B) {
	for i := 0; i < b.N; i++ {
		S2Point(points[i%len(points)])
	}
}

func BenchmarkQuery(b *testing.B) {
	s := NewShapeIndexSearcher()
	for i := 0; i < b.N; i++ {
		s.Query()
	}
}

func BenchmarkH3Lookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		point := points[i%len(points)]
		latLng := h3.NewLatLng(point[0], point[1])
		resolution := 9 // between 0 (biggest cell) and 15 (smallest cell)
		h3.LatLngToCell(latLng, resolution)
	}
}
