package world

import (
	"image/color"
	"reflect"
	"testing"
)

func TestCollisionWorld_InsertCollider(t *testing.T) {

	width := 100
	height := 100
	collisionMap := NewEmptyCollisionWorld(CollisionWorldSettings{
		Dimensions: Dimensions{Width: width, Height: height},
		IsDiagonal: false},
	)

	type args struct {
		x int
		y int
		c *Collider
	}
	type expected struct {
		wantX    int
		wantY    int
		inserted bool
	}
	tests := []struct {
		name         string
		collisionMap *CollisionWorld
		args         args
		want         expected
	}{
		{"Insert Within Map", &collisionMap, args{5, 5, &Collider{}}, expected{5, 5, true}},
		{"Insert Into Same Position As Above", &collisionMap, args{5, 5, &Collider{}}, expected{0, 0, false}},
		{"Insert Outside Map", &collisionMap, args{-50, -50, &Collider{}}, expected{0, 0, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collisionMap.InsertCollider(tt.args.x, tt.args.y, tt.args.c); got != tt.want.inserted {
				t.Errorf("CollisionWorld.InsertCollider() = %v, want %v", got, tt.want.inserted)
			}

			if tt.args.c.GetX() != tt.want.wantX || tt.args.c.GetY() != tt.want.wantY {
				t.Errorf("Collider Position = (%d, %d), want (%d, %d)", tt.args.c.GetX(), tt.args.c.GetY(), tt.want.wantX, tt.want.wantY)
			}
		})
	}
}

func TestCollisionWorld_HasCollider(t *testing.T) {

	width := 5
	height := 5
	collisionMap := NewEmptyCollisionWorld(CollisionWorldSettings{
		Dimensions: Dimensions{Width: width, Height: height},
		IsDiagonal: false},
	)

	insertX, insertY := 1, 1
	collider := Collider{}
	collisionMap.InsertCollider(insertX, insertY, &collider)

	type args struct {
		x int
		y int
	}
	tests := []struct {
		name         string
		collisionMap *CollisionWorld
		args         args
		want         *Collider
		want1        bool
	}{
		{"Has Collider", &collisionMap, args{insertX, insertY}, &collider, true},
		{"Does Not Have Collider", &collisionMap, args{2, 2}, nil, false},
		{"Out of Bounds", &collisionMap, args{-2, -2}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.collisionMap.HasCollider(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollisionWorld.HasCollider() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CollisionWorld.HasCollider() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCollisionWorld_MoveCollider(t *testing.T) {
	type args struct {
		moveX int
		moveY int
		c     *Collider
	}
	type expected struct {
		x          int
		y          int
		funcReturn bool
	}

	width := 5
	height := 5
	collisionMap := NewEmptyCollisionWorld(CollisionWorldSettings{
		Dimensions: Dimensions{Width: width, Height: height},
		Population: Population{InitialPopulation: 100, MaxPopulation: 100},
		IsDiagonal: false},
	)

	insertX, insertY := 1, 1
	collider := &Collider{}
	collisionMap.InsertCollider(insertX, insertY, collider)

	tests := []struct {
		name         string
		collisionMap *CollisionWorld
		args         args
		want         expected
	}{
		{"Move into Empty Space", &collisionMap, args{0, 1, collider}, expected{insertX, insertY + 1, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.collisionMap.MoveCollider(tt.args.moveX, tt.args.moveY, tt.args.c); got != tt.want.funcReturn {
				t.Errorf("CollisionWorld.MoveCollider() = %v, want %v", got, tt.want.funcReturn)
			}

			tt.collisionMap.Process()

			if tt.args.c.GetX() != tt.want.x || tt.args.c.GetY() != tt.want.y {
				t.Errorf("Collider Position = (%d, %d), want (%d, %d)", tt.args.c.GetX(), tt.args.c.GetY(), tt.want.x, tt.want.y)
			}
		})
	}
}

func TestCollider_ChangeColor(t *testing.T) {

	type expected struct {
		color color.RGBA
	}

	tests := []struct {
		name     string
		collider *Collider
		want     expected
	}{
		{"Negative ColorSelection", &Collider{colorSelection: -1}, expected{colliderColors[0]}},
		{"Large than len ColorSelection", &Collider{colorSelection: len(colliderColors) - 1}, expected{colliderColors[0]}},
		{"Normal Collider", &Collider{colorSelection: 0}, expected{colliderColors[1]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.collider.ChangeColor()
			if tt.collider.Color != tt.want.color {
				t.Errorf("Collider.Color = %v, want %v", tt.collider.Color, tt.want)
			}
		})
	}
}
