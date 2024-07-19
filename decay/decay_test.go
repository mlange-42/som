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
		{
			name:  "Polynomial",
			start: 0.5,
			end:   0.1,
			d: &decay.Polynomial{
				Start: 0.5,
				End:   0.1,
				Exp:   2,
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
func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    decay.Decay
		wantErr bool
	}{
		{
			name:  "Linear with default args",
			input: "linear",
			want:  &decay.Linear{},
		},
		{
			name:  "Power with custom args",
			input: "power 0.8 0.2",
			want:  &decay.Power{Start: 0.8, End: 0.2},
		},
		{
			name:  "Polynomial with custom args",
			input: "polynomial 0.9 0.1 3",
			want:  &decay.Polynomial{Start: 0.9, End: 0.1, Exp: 3},
		},
		{
			name:    "Unknown decay function",
			input:   "unknown",
			wantErr: true,
		},
		{
			name:    "Invalid number of arguments",
			input:   "linear 0.5 0.1 1",
			wantErr: true,
		},
		{
			name:    "Invalid argument type",
			input:   "linear 0.5 invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decay.FromString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.want, got)
				if !tt.wantErr {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func TestFromString_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    decay.Decay
		wantErr bool
	}{
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Only spaces",
			input:   "   ",
			wantErr: true,
		},
		{
			name:  "Linear with zero values",
			input: "linear 0 0",
			want:  &decay.Linear{Start: 0, End: 0},
		},
		{
			name:  "Power with negative values",
			input: "power -0.5 -0.1",
			want:  &decay.Power{Start: -0.5, End: -0.1},
		},
		{
			name:    "Polynomial with too many arguments",
			input:   "polynomial 0.9 0.1 3 2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decay.FromString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.want, got)
				if !tt.wantErr {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}
