package generator

import (
	"math/rand"
	"strconv"
	"time"
)

func Generate() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	generated := rng.Intn(99) + 1
	return strconv.Itoa(generated)
}
