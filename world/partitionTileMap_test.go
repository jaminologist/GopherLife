package world

import (
	"gopherlife/geometry"
	"testing"
)

func TestPartitionTileMap_Tile(t *testing.T) {

}

func TestGridSection_Tile(t *testing.T) {
}

func TestPartitionTileMap_MoveGopher(t *testing.T) {

	stats := Statistics{
		Width:                  10,
		Height:                 10,
		NumberOfGophers:        0,
		NumberOfFood:           20,
		MaximumNumberOfGophers: 100000,
		GopherBirthRate:        7,
	}

	var tileMap = CreatePartitionTileMapCustom(stats)

	gopher := NewGopher("a", geometry.NewCoordinate(1, 2))
	tileMap.InsertGopher(1, 2, &gopher)

	tileMap.ActionQueuer.Add(func() {
		tileMap.MoveGopher(&gopher, 0, 1)
	})
	tileMap.Update()

	des, _ := tileMap.Tile(1, 3)

	if des.Gopher == nil {
		t.Errorf("Destination is empty")
	}

	prev, _ := tileMap.Tile(1, 2)

	if prev.Gopher != nil {
		t.Errorf("Previous Destination is not empty")
	}

}

func TestPartitionTileMap_RemoveGopher(t *testing.T) {

	stats := Statistics{
		Width:                  10,
		Height:                 10,
		NumberOfGophers:        0,
		NumberOfFood:           20,
		MaximumNumberOfGophers: 100000,
		GopherBirthRate:        7,
	}

	var tileMap = CreatePartitionTileMapCustom(stats)

	gopher := NewGopher("a", geometry.Coordinates{1, 2})
	tileMap.InsertGopher(1, 2, &gopher)

	bool := tileMap.RemoveGopher(1, 2)

	if !bool {
		t.Errorf("Gopher is not removed")
	}

	tile, _ := tileMap.Tile(1, 2)

	if tile.Gopher != nil {
		t.Errorf("Gopher is not removed")
	}
}
