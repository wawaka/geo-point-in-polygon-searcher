package searcher

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// func TestPolygon(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		polygon := GeneratePolygon(10)
// 		fmt.Println(len(polygon))
// 	}
// 	// assert.Fail(t, "")
// }
