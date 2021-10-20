package grand

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix() + 100))
}
func Int63() int64 {
	return r.Int63()
}
