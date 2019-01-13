package world

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPartitionTileMap_Tile(t *testing.T) {

	tileMap := CreatePartitionTileMap()

	tile, ok := tileMap.Tile(12, 29)

	fmt.Println(len(tileMap.grid))
	fmt.Println(len(tileMap.grid[0]))
	fmt.Println("Hi?")
	fmt.Println(tileMap.convertToGridCoordinates(12, 29))

	fmt.Println(ok)
	fmt.Println(tile.isEmpty())
}

func TestGridSection_Tile(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name        string
		gridSection *GridSection
		args        args
		want        *Tile
		want1       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.gridSection.Tile(tt.args.x, tt.args.y)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GridSection.Tile() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GridSection.Tile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
