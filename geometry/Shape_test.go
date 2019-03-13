package geometry

import "testing"

func TestRectangle_Contains(t *testing.T) {
	type args struct {
		x    int
		y    int
		want bool
	}
	tests := []struct {
		name string
		r    *Rectangle
		args []args
	}{
		{"Rectangle Starting at Origin", &Rectangle{0, 0, 10, 10},
			[]args{
				{1, 1, true},
				{0, 0, true},
				{10, 10, false},
				{-1, -1, false},
			},
		},

		{"Rectangle Starting at Negative Point", &Rectangle{-50, -50, 10, 10},
			[]args{
				{-50, -50, true},
				{-40, -40, false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, args := range tt.args {
				if got := tt.r.Contains(args.x, args.y); got != args.want {
					t.Errorf("Rectangle.Contains() = %v, want %v", got, args.want)
				}
			}
		})
	}
}
