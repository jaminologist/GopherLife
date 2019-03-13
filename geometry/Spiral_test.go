package geometry

import (
	"reflect"
	"testing"
)

func BenchmarkSpiral(b *testing.B) {

	spiral := NewSpiral(1000, 1000)

	for n := 0; n < b.N; n++ {
		spiral.Next()
	}
}

func TestSpiral_Next(t *testing.T) {

	spiral := NewSpiral(3, 3)

	type spiralResult struct {
		coordinates Coordinates
		want        bool
	}

	tests := []struct {
		name          string
		s             *Spiral
		spiralResults []spiralResult
	}{
		{"Spiral 3x3 Test", &spiral,
			[]spiralResult{
				{Coordinates{0, 0}, true},
				{Coordinates{1, 0}, true},
				{Coordinates{1, 1}, true},
				{Coordinates{0, 1}, true},
				{Coordinates{-1, 1}, true},
				{Coordinates{-1, 0}, true},
				{Coordinates{-1, -1}, true},
				{Coordinates{0, -1}, true},
				{Coordinates{1, -1}, true},
				{Coordinates{0, 0}, false},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, result := range tt.spiralResults {
				got, got1 := tt.s.Next()
				if !reflect.DeepEqual(got, result.coordinates) {
					t.Errorf("Spiral.Next() got = %v, want %v", got, result.coordinates)
				}
				if got1 != result.want {
					t.Errorf("Spiral.Next() got1 = %v, want %v", got1, result.want)
				}
			}

		})
	}
}
