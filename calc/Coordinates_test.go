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
