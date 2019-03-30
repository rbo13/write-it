// Package generate handles different types for generator.
package generate

import (
	"math/rand"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Color generates a color.
func Color() string {
	c := colorful.Hsv(rand.Float64()*360.0, 0.8, 0.8)
	return c.Hex()
}
