package calc

import (
	"reflect"
	"testing"
)

func TestStringToCoordinates(t *testing.T) {
	type args struct {
		coordString string
	}
	tests := []struct {
		name string
		args args
		want Coordinates
	}{
		{"Hey", args{"2,5"}, Coordinates{X: 2, Y: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToCoordinates(tt.args.coordString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateCoordinateArray(t *testing.T) {
	type args struct {
		startX int
		startY int
		endX   int
		endY   int
	}
	tests := []struct {
		name string
		args args
		want []Coordinates
	}{
		{
			name: "Test",
			args: args{
				0, 0, 2, 2,
			},
			want: []Coordinates{NewCoordinate(0, 0), NewCoordinate(0, 1), NewCoordinate(1, 0), NewCoordinate(1, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCoordinateArray(tt.args.startX, tt.args.startY, tt.args.endX, tt.args.endY); !reflect.DeepEqual(got, tt.want) {

				t.Errorf("GenerateCoordinateArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoordinates_IsInRange(t *testing.T) {
	type args struct {
		c2   Coordinates
		minX int
		minY int
	}
	tests := []struct {
		name string
		c    Coordinates
		args args
		want bool
	}{
		{
			"In range, size 1",
			Coordinates{1, 1},
			args{Coordinates{0, 0}, 1, 1},
			true,
		},
		{
			"Out of range, size 1",
			Coordinates{2, 1},
			args{Coordinates{0, 0}, 1, 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsInRange(tt.args.c2, tt.args.minX, tt.args.minY); got != tt.want {
				t.Errorf("Coordinates.IsInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortByNearestFromCoordinate(t *testing.T) {
	type args struct {
		coords Coordinates
		cs     []Coordinates
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortByNearestFromCoordinate(tt.args.coords, tt.args.cs)
		})
	}
}
