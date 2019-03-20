package world

import (
	"reflect"
	"testing"
)

func TestBasic2DContainer_HasGopher(t *testing.T) {

	gridC := NewBasic2DContainer(0, 0, 10, 10)
	gopher1 := &Gopher{}
	gridC.InsertGopher(1, 1, gopher1)
	gridC.InsertGopher(1, 2, &Gopher{})
	gridC.RemoveGopher(1, 2)

	type args struct {
		x int
		y int
	}
	tests := []struct {
		name      string
		container *Basic2DContainer
		args      args
		want      *Gopher
		want1     bool
	}{
		{"Gopher Should Exist", &gridC, args{1, 1}, gopher1, true},
		{"Gopher Should Not Exist", &gridC, args{1, 2}, nil, false},
		{"Out of Bounds", &gridC, args{100, 200}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.container.HasGopher(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Basic2DContainer.HasGopher() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Basic2DContainer.HasGopher() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBasic2DContainer_HasFood(t *testing.T) {

	gridC := NewBasic2DContainer(0, 0, 10, 10)
	food1 := &Food{}
	gridC.InsertFood(1, 1, food1)
	gridC.InsertFood(1, 2, &Food{})
	gridC.RemoveFood(1, 2)

	type args struct {
		x int
		y int
	}
	tests := []struct {
		name      string
		container *Basic2DContainer
		args      args
		want      *Food
		want1     bool
	}{
		{"Food Should Exist", &gridC, args{1, 1}, food1, true},
		{"Food Should Not Exist", &gridC, args{1, 2}, nil, false},
		{"Out of Bounds", &gridC, args{100, 200}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.container.HasFood(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Basic2DContainer.HasFood() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Basic2DContainer.HasFood() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBasicGridContainer_HasGopher(t *testing.T) {

	gridC := NewBasicGridContainer(10, 10, 5, 5)
	gopher1 := &Gopher{}
	gridC.InsertGopher(1, 1, gopher1)
	gridC.InsertGopher(1, 2, &Gopher{})
	gridC.RemoveGopher(1, 2)

	type args struct {
		x int
		y int
	}
	tests := []struct {
		name      string
		container *BasicGridContainer
		args      args
		want      *Gopher
		want1     bool
	}{
		{"Gopher Should Exist", &gridC, args{1, 1}, gopher1, true},
		{"Gopher Should Not Exist", &gridC, args{1, 2}, nil, false},
		{"Out of Bounds", &gridC, args{100, 200}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.container.HasGopher(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicGridContainer.HasGopher() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BasicGridContainer.HasGopher() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBasicGridContainer_HasFood(t *testing.T) {

	gridC := NewBasicGridContainer(10, 10, 5, 5)
	food1 := &Food{}
	gridC.InsertFood(1, 1, food1)
	gridC.InsertFood(1, 2, &Food{})
	gridC.RemoveFood(1, 2)

	type args struct {
		x int
		y int
	}
	tests := []struct {
		name      string
		container *BasicGridContainer
		args      args
		want      *Food
		want1     bool
	}{
		{"Food Should Exist", &gridC, args{1, 1}, food1, true},
		{"Food Should Not Exist", &gridC, args{1, 2}, nil, false},
		{"Out of Bounds", &gridC, args{100, 200}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.container.HasFood(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicGridContainer.HasFood() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BasicGridContainer.HasFood() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
