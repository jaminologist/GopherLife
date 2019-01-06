package world

import (
	"reflect"
	"testing"
)

func TestMovingGophers(t *testing.T) {
	world := CreateTileMap()
	world.MoveGopher(world.SelectedGopher, 1, 1)
}

func TestGender_Opposite(t *testing.T) {
	tests := []struct {
		name   string
		gender Gender
		want   Gender
	}{
		// TODO: Add test cases.
		{
			name:   "first",
			gender: Male,
			want:   Female,
		},

		{
			name:   "second",
			gender: Female,
			want:   Male,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gender.Opposite(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gender.Opposite() = %v, want %v", got, tt.want)
			}
		})
	}
}
