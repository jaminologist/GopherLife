package world

import (
	"fmt"
	"gopherlife/calc"
	"testing"
)

func BenchmarkFibProcessWorld(b *testing.B) {
	// run the Fib function b.N times

	world := CreateWorldCustom(3000, 3000, 50000, 10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		world.ProcessWorld()
	}
}

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
	}
}

func BenchmarkFibProcessWorld1(b *testing.B) {
	// run the Fib function b.N times

	world := CreateWorld()

	for n := 0; n < b.N; n++ {
		world.ProcessWorld()
	}
}

func TestWorld_ProcessWorld(t *testing.T) {
	tests := []struct {
		name  string
		world *World
		want  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.world.ProcessWorld(); got != tt.want {
				t.Errorf("World.ProcessWorld() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorld_InsertGopher(t *testing.T) {

	worldA := CreateWorldCustom(10, 10, 0, 0)
	gopherA := NewGopher("Harry", calc.NewCoordinate(0, 0))
	worldB := CreateWorldCustom(10, 10, 0, 0)
	gopherB := NewGopher("Harry", calc.NewCoordinate(0, 0))

	type args struct {
		gopher *Gopher
		x      int
		y      int
	}
	tests := []struct {
		name  string
		world *World
		args  args
		want  bool
	}{
		{
			name:  "Insert Gopher into Empty Space",
			world: &worldA,
			args: args{
				gopher: &gopherA,
				x:      5,
				y:      5,
			},
			want: true,
		},

		{
			name:  "Insert Gopher into Occupied Space",
			world: &worldB,
			args: args{
				gopher: &gopherA,
				x:      5,
				y:      5,
			},
			want: false,
		},
	}

	worldB.InsertGopher(&gopherB, 5, 5)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.world.InsertGopher(tt.args.gopher, tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("World.InsertGopher() = %v, want %v", got, tt.want)
			}
		})
	}
}
