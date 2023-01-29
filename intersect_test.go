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

func TestIntersectCommutative(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		l1 := Line{Point{rand.Float64(), rand.Float64()}, Point{rand.Float64(), rand.Float64()}}
		l2 := Line{Point{rand.Float64(), rand.Float64()}, Point{rand.Float64(), rand.Float64()}}
		assert.Equal(t, Intersect(l1, l2), Intersect(l2, l1))
	}
}

func TestIntersect(t *testing.T) {
	results := map[bool]int{}
	for i := 0; i < 1000000; i++ {
		l1 := Line{Point{rand.Float64(), rand.Float64()}, Point{rand.Float64(), rand.Float64()}}
		l2 := Line{Point{rand.Float64(), rand.Float64()}, Point{rand.Float64(), rand.Float64()}}
		// assert.Equal(t, Intersect(l1, l2), Intersect(l2, l1))
		results[Intersect(l1, l2)]++
	}
	// fmt.Println(results)
	// assert.Fail(t, "")
}
