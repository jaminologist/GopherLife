package geometry

import (
	"fmt"
	"testing"
)

func TestSpiralCoordinates(t *testing.T) {

	spiral := NewSpiral(3, 3)

	for i := 0; i < 20; i++ {
		fmt.Println(spiral.Next())
	}

}

func TestSpiral_Next(t *testing.T) {

	spiral := NewSpiral(3, 3)

	for {
		c, ok := spiral.Next()
		if !ok {
			break
		}
		fmt.Println(c)
	}

}

func BenchmarkSpiral(b *testing.B) {

	spiral := NewSpiral(1000, 1000)

	for n := 0; n < b.N; n++ {
		spiral.Next()
	}
}
