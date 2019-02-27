package world

import (
	"fmt"
	"reflect"
	"testing"
)

type AnythingContainer struct {
	id  int
	odj interface{}
}

func (ac *AnythingContainer) Get(id int) interface{} {
	if id == ac.id {
		return ac.odj
	}
	return nil
}

func TestMovingGophers(t *testing.T) {

	food := NewPotato()

	//fmt.Println(food.Name)

	ac := AnythingContainer{1, &food}

	newfood := ac.Get(1)

	switch t := newfood.(type) {
	case *Food:
		fmt.Println(t.Name)
	}
	//newfood.Name = "Happy"
	//fmt.Println(food.Name)
	//fmt.Println(newfood.Name)

	//world := CreateTileMap()
	//world.MoveGopher(world.SelectedGopher, 1, 1)
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
