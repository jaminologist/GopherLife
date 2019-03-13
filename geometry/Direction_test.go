package geometry

import (
	"reflect"
	"testing"
)

func TestDirection_TurnClockWise90(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want Direction
	}{
		{"Up", Up, Right},
		{"Right", Right, Down},
		{"Down", Down, Left},
		{"Left", Left, Up},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.TurnClockWise90(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Direction.TurnClockWise90() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirection_TurnAntiClockWise90(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want Direction
	}{
		{"Up", Up, Left},
		{"Left", Left, Down},
		{"Down", Down, Right},
		{"Right", Right, Up},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.TurnAntiClockWise90(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Direction.TurnAntiClockWise90() = %v, want %v", got, tt.want)
			}
		})
	}
}
