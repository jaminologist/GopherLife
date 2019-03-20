package world

import (
	"gopherlife/geometry"
	"testing"
)

func TestSnakePart_AttachToBack(t *testing.T) {
	type args struct {
		partToAttach *SnakePart
	}
	tests := []struct {
		name string
		sp   *SnakePart
		args args
	}{
		{"Attach To Back", &SnakePart{}, args{&SnakePart{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sp.AttachToBack(tt.args.partToAttach)
			if tt.sp.snakePartBehind == nil {
				t.Errorf("Error SnakePart not attached to back")
			}

			if tt.args.partToAttach.snakePartInFront == nil {
				t.Errorf("Error Attached SnakePart does not have SnakePart attached to Front")
			}

		})
	}
}

func TestSnakeWorld_ChangeDirection(t *testing.T) {
	type args struct {
		d geometry.Direction
	}
	tests := []struct {
		name string
		smt  *SnakeWorld
		args args
		want geometry.Direction
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.smt.ChangeDirection(tt.args.d)
		})
	}
}

func TestSnakeWorld_InsertSnakePart(t *testing.T) {

	sw := NewEmptySnakeWorld(SnakeWorldSettings{Dimensions{10, 10}, 5})

	snakeFoodX, snakeFoodY := 5, 5
	snakeWallX, snakeWallY := 6, 6

	sw.InsertSnakeFood(snakeFoodX, snakeFoodY, &SnakeFood{})
	sw.InsertSnakeWall(snakeWallX, snakeWallY, &SnakeWall{})

	type args struct {
		x int
		y int
		p *SnakePart
	}
	type expected struct {
		funcReturn bool
		x          int
		y          int
	}
	tests := []struct {
		name string
		smt  *SnakeWorld
		args args
		want expected
	}{
		{"Insert Snake Part Into Empty Space", &sw, args{1, 1, &SnakePart{}}, expected{true, 1, 1}},
		{"Insert Snake Part Into Occupied Space (SnakePart)", &sw, args{1, 1, &SnakePart{}}, expected{false, 0, 0}},
		{"Insert Snake Part Into Occupied Space (SnakeFood)", &sw, args{snakeFoodX, snakeFoodY, &SnakePart{}}, expected{true, snakeFoodX, snakeFoodY}},
		{"Insert Snake Part Into Occupied Space (SnakeWall)", &sw, args{snakeWallX, snakeWallY, &SnakePart{}}, expected{false, 0, 0}},
		{"Insert Snake Part Out of Bounds", &sw, args{-100, -100, &SnakePart{}}, expected{false, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.smt.InsertSnakePart(tt.args.x, tt.args.y, tt.args.p); got != tt.want.funcReturn {
				t.Errorf("SnakeWorld.InsertSnakePart() = %v, want %v", got, tt.want)
			}

			if tt.args.p.GetX() != tt.want.x || tt.args.p.GetY() != tt.want.y {
				t.Errorf("SnakeWorld.InsertSnakePart() SnakePart Position = (%d, %d) want (%d, %d)", tt.args.p.GetX(), tt.args.p.GetY(), tt.want.x, tt.want.y)
			}

		})
	}
}
