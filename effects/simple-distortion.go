//go:build ignore

package effects

import "math"

var Name = "simple distortion"
var Enabled = false

var Gain float32 = 4.0
var Level float32 = 1
var loopLevel int = 1

func Process(samples []float32) {
	for n := 0; n <= loopLevel; n++ {
		for i := range samples {
			sample := samples[i] * Gain
			sample = float32(math.Tanh(float64(sample)))
			samples[i] = sample * Level
		}
	}
}
