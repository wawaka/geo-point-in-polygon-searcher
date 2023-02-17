package searcher

import (
	"github.com/engelsjk/polygol"
)

type MultiPolygonWithBounds struct {
	MultiPolygon [][][][]float64
	Bounds       []BoundingBox
}

func (mpwb *MultiPolygonWithBounds) RecalcBounds() {
	mpwb.Bounds = make([]BoundingBox, 0)
	for _, polygon := range mpwb.MultiPolygon {
		bbox := BoundingBox{}
		for _, p := range polygon[0] {
			bbox.AddPoint(p)
		}
		mpwb.Bounds = append(mpwb.Bounds, bbox)
	}
}

func (mpwb *MultiPolygonWithBounds) Contains(p []float64) bool {
	for i, bbox := range mpwb.Bounds {
		if bbox.Contains(p) {
			if PolygonWithHolesContains(mpwb.MultiPolygon[i], p) {
				return true
			}
		}
	}
	return false
}

type BoundingBoxOptimizedSearcher struct {
	MultiPolygonsMap map[uint32]*MultiPolygonWithBounds
}

func NewBoundingBoxOptimizedSearcher() *BoundingBoxOptimizedSearcher {
	s := &BoundingBoxOptimizedSearcher{}
	s.MultiPolygonsMap = make(map[uint32]*MultiPolygonWithBounds)
	return s
}

func (s *BoundingBoxOptimizedSearcher) AddRing(id uint32, ring [][]float64) error {
	return s.AddMultiPolygon(id, [][][][]float64{{ring}})
}

func (s *BoundingBoxOptimizedSearcher) AddMultiPolygon(id uint32, multipolygon [][][][]float64) error {
	mpwb, found := s.MultiPolygonsMap[id]
	if found {
		var err error
		mpwb.MultiPolygon, err = polygol.Union(mpwb.MultiPolygon, multipolygon)
		if err != nil {
			return err
		}
	} else {
		mpwb = &MultiPolygonWithBounds{MultiPolygon: multipolygon}
		s.MultiPolygonsMap[id] = mpwb
	}

	mpwb.RecalcBounds()
	return nil
}

func (s *BoundingBoxOptimizedSearcher) SearchPoint(p []float64) []uint32 {
	return s.SearchPointFiltered(p)
}

func (s *BoundingBoxOptimizedSearcher) SearchPointFiltered(p []float64) []uint32 {
	results := []uint32{}
	for id, mpwb := range s.MultiPolygonsMap {
		if mpwb.Contains(p) {
			results = append(results, id)
		}
	}
	return results
}

func (s *BoundingBoxOptimizedSearcher) SearchPointBruteforce(p []float64) []uint32 {
	results := []uint32{}
	for id, mpwb := range s.MultiPolygonsMap {
		if MultiPolygonContains(mpwb.MultiPolygon, p) {
			results = append(results, id)
		}
	}
	return results
}

func (s *BoundingBoxOptimizedSearcher) Remove(id uint32) {
	delete(s.MultiPolygonsMap, id)
}
