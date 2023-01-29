package searcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testPointsMax = 1000000

var testPoints [][]float64

var _ PointInPolygonSearcher = NewBoundingBoxOptimizedSearcher() // assert interface contract

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < testPointsMax; i++ {
		point := []float64{rand.Float64()*4 - 2, rand.Float64()*4 - 2}
		testPoints = append(testPoints, point)
	}
}

func TestRecalcBounds(t *testing.T) {
	mpwb := MultiPolygonWithBounds{MultiPolygon: [][][][]float64{{{{0, 1}, {1, 2}, {2, 3}}}}}
	mpwb.RecalcBounds()
	assert.Equal(t, []float64{0, 1, 2, 3}, mpwb.Bounds[0].MinMax)
}

func makeSimpleSearcher() *BoundingBoxOptimizedSearcher {
	s := NewBoundingBoxOptimizedSearcher()
	s.AddMultiPolygon(1, [][][][]float64{{{
		{0, 1},
		{1, 0},
		{0, -1},
		{-1, 0},
	}}})
	return s
}

func TestEmpty(t *testing.T) {
	s := NewBoundingBoxOptimizedSearcher()
	assert.Equal(t, []uint32{}, s.SearchPoint([]float64{0, 0}))
}

func TestSimple(t *testing.T) {
	s := makeSimpleSearcher()
	assert.Equal(t, []uint32{1}, s.SearchPoint([]float64{0, 0}))
	assert.Equal(t, []uint32{}, s.SearchPoint([]float64{1, 1}))
}

func TestEquivalence(t *testing.T) {
	s := makeSimpleSearcher()
	for i := 0; i < testPointsMax; i++ {
		assert.Equal(t, s.SearchPointBruteforce(testPoints[i]), s.SearchPointFiltered(testPoints[i]))
	}
}

func TestSearchPoint(t *testing.T) {
	s := NewBoundingBoxOptimizedSearcher()

	m := map[uint32][][][]float64{}

	for i := 0; i < 100; i++ {
		id := uint32(rand.Intn(100))
		polygon := GeneratePolygon(10)
		m[id] = append(m[id], polygon)
		s.AddMultiPolygon(id, [][][][]float64{{polygon}})
	}

	for i := 0; i < 10000; i++ {
		point := GeneratePoint()

		found_expected := map[uint32]bool{}
		for id, polygons := range m {
			for _, polygon := range polygons {
				if PolygonContains(polygon, point) {
					found_expected[id] = true
				}
			}
		}

		found_actual := map[uint32]bool{}
		for _, id := range s.SearchPoint(point) {
			found_actual[id] = true
		}
		assert.Equal(t, found_expected, found_actual)
	}
}

func BenchmarkSearcherBruteforce(b *testing.B) {
	s := makeSimpleSearcher()
	for i := 0; i < b.N; i++ {
		s.SearchPointBruteforce(testPoints[i%testPointsMax])
	}
}

func BenchmarkSearcherFiltered(b *testing.B) {
	s := makeSimpleSearcher()
	for i := 0; i < b.N; i++ {
		s.SearchPointFiltered(testPoints[i%testPointsMax])
	}
}
