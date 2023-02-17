package searcher

type BoundingBox struct {
	MinMax []float64
}

func (bb *BoundingBox) AddPoint(p []float64) {
	if len(bb.MinMax) == 4 {
		if p[0] < bb.MinMax[0] {
			bb.MinMax[0] = p[0]
		}
		if p[1] < bb.MinMax[1] {
			bb.MinMax[1] = p[1]
		}
		if p[0] > bb.MinMax[2] {
			bb.MinMax[2] = p[0]
		}
		if p[1] > bb.MinMax[3] {
			bb.MinMax[3] = p[1]
		}
	} else {
		bb.MinMax = []float64{p[0], p[1], p[0], p[1]}
	}
}

func (bb *BoundingBox) Contains(p []float64) bool {
	if len(bb.MinMax) == 4 {
		return bb.MinMax[0] <= p[0] && p[0] <= bb.MinMax[2] && bb.MinMax[1] <= p[1] && p[1] <= bb.MinMax[3]
	} else {
		return false
	}
}

func (bb *BoundingBox) Min() []float64 {
	return bb.MinMax[0:2]
}

func (bb *BoundingBox) Max() []float64 {
	return bb.MinMax[2:4]
}
