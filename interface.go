package searcher

type PointInPolygonSearcher interface {
	AddRing(id uint32, ring [][]float64) error
	SearchPoint(point []float64) []uint32 // returns all ids of matched multipolygons, empty slice if not found
}
