package searcher

import (
	"github.com/golang/geo/s2"
)

type ShapeIndexSearcher struct {
	ShapeIndex  *s2.ShapeIndex
	ShapeLookup map[s2.Shape]uint32
}

func NewShapeIndexSearcher() *ShapeIndexSearcher {
	return &ShapeIndexSearcher{
		ShapeIndex:  s2.NewShapeIndex(),
		ShapeLookup: make(map[s2.Shape]uint32),
	}
}

func (s *ShapeIndexSearcher) AddRing(id uint32, ring [][]float64) error {
	s2loop := S2Loop(ring)
	return s.AddShape(id, s2loop)
}

func (s *ShapeIndexSearcher) AddShape(id uint32, shape s2.Shape) error {
	s.ShapeIndex.Add(shape)
	s.ShapeLookup[shape] = id
	return nil
}

func (s *ShapeIndexSearcher) Query() *s2.ContainsPointQuery {
	return s2.NewContainsPointQuery(s.ShapeIndex, s2.VertexModelOpen)
}

func (s *ShapeIndexSearcher) SearchS2PointQuery(query *s2.ContainsPointQuery, point s2.Point) []uint32 {
	shapes := query.ContainingShapes(point)
	ids := []uint32{}
	for _, shape := range shapes {
		ids = append(ids, s.ShapeLookup[shape])
	}
	return ids
}

func (s *ShapeIndexSearcher) SearchPointQuery(q *s2.ContainsPointQuery, p []float64) []uint32 {
	return s.SearchS2PointQuery(q, S2Point(p))
}

func (s *ShapeIndexSearcher) SearchPoint(p []float64) []uint32 {
	return s.SearchPointQuery(s.Query(), p)
}

func (s *ShapeIndexSearcher) SearchS2Point(p s2.Point) []uint32 {
	return s.SearchS2PointQuery(s.Query(), p)
}
