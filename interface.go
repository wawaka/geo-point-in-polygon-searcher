package searcher

type PointInPolygonSearcher interface {
	AddMultiPolygon(id uint32, multipolygon [][][][]float64) error // multiple multipolygons can be assigned to the same id
	Remove(id uint32)                                              // removes all previously added multipolygons with given id
	SearchPoint(point []float64) []uint32                          // returns all ids of matched multipolygons, empty slice if not found
}
