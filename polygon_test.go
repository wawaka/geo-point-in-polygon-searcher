package searcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPolygon(t *testing.T) {
	for i := 0; i < 1000; i++ {
		polygon := GeneratePolygon(10)
		assert.GreaterOrEqual(t, len(polygon), 10)
		assert.NotEqual(t, polygon[0], polygon[len(polygon)-1], "%v %v", polygon[0], polygon[len(polygon)-1])
	}
}
