package decay_test

import (
	"testing"

	"github.com/mlange-42/som/decay"
	"github.com/stretchr/testify/assert"
)

func TestDecay(t *testing.T) {
	tMax := 25

	tests := []struct {
		name  string
		start float64
		end   float64
		d     decay.Decay
	}{
		{
			name:  "Linear",
			start: 0.5,
			end:   0.1,
			d: &decay.Linear{
				Start: 0.5,
				End:   0.1,
			},
		},
		{
			name:  "Power",
			start: 0.5,
			end:   0.1,
			d: &decay.Power{
				Start: 0.5,
				End:   0.1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.d.Decay(0, tMax)
			assert.Equal(t, tt.start, v)
			for i := 1; i < tMax; i++ {
				v2 := tt.d.Decay(i, tMax)
				assert.Less(t, v2, v)
				v = v2
			}
			v = tt.d.Decay(tMax, tMax)
			assert.Equal(t, tt.end, v)
		})
	}
}
