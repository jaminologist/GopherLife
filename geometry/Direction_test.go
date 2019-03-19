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

func TestDirection_AddToPoint(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name  string
		d     Direction
		args  args
		want  int
		want1 int
	}{
		{"Up", Up, args{0, 0}, 0, 1},
		{"Left", Left, args{0, 0}, -1, 0},
		{"Down", Down, args{0, 0}, 0, -1},
		{"Right", Right, args{0, 0}, 1, 0},
		{"Up Negative", Up, args{-20, -20}, -20, -19},
		{"Left Negative", Left, args{-20, -20}, -21, -20},
		{"Down Negative", Down, args{-20, -20}, -20, -21},
		{"Right Negative", Right, args{-20, -20}, -19, -20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.d.AddToPoint(tt.args.x, tt.args.y)
			if got != tt.want {
				t.Errorf("Direction.AddToPoint() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Direction.AddToPoint() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
