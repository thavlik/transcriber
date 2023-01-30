package base

import (
	"math/rand"
	"time"
)

func RandomizeSeed() {
	rand.Seed(time.Now().UTC().UnixNano())
}
