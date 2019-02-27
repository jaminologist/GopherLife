package world

import (
	"fmt"
	"testing"
)

func TestBasicGridContainer_Grid(t *testing.T) {

	bgc := NewBasicGridContainer(10, 10, 5, 5)

	_, ok := bgc.Grid(0, 0)

	fmt.Println(ok)

	_, ok = bgc.Grid(10, 10)

	fmt.Println(ok)

}
