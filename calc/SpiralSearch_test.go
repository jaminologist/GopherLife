package calc

import (
	"fmt"
	"testing"
)

func TestSpiralCoordinates(t *testing.T) {

	spiral := NewSpiral(15, 15)

	for i := 0; i < 200; i++ {
		fmt.Println(spiral.Next())
	}

}
