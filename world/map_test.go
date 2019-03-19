package world

import (
	"testing"
)

func BenchmarkHello(b *testing.B) {

	a := Basic2DContainer{}
	c := interface{}(a)

	for i := 0; i < b.N; i++ {

		switch v := c.(type) {
		case Gopher:
			v.Gender = Male
		case Basic2DContainer:

		}
	}
}

func BenchmarkFibProcessWorld1(b *testing.B) {
	// run the Fib function b.N times

	//	world := CreateTileMap()

	//	for n := 0; n < b.N; n++ {
	//		world.Update()
	//	}
}

/*func TestWorld_ProcessWorld(t *testing.T) {
	tests := []struct {
		name  string
		world *World
		want  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.world.Update(); got != tt.want {
				t.Errorf("World.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

/*func TestWorld_InsertGopher(t *testing.T) {

	worldA := CreateGopherWorldSpiralSearch(10, 10, 0, 0)
	gopherA := NewGopher("Harry", geometry.NewCoordinate(0, 0))
	worldB := CreateGopherWorldSpiralSearch(10, 10, 0, 0)
	gopherB := NewGopher("Harry", geometry.NewCoordinate(0, 0))

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
} */
