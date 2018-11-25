package world

import (
	"fmt"
	"testing"
)

func BenchmarkFibProcessWorld(b *testing.B) {
	// run the Fib function b.N times

	world := CreateWorldCustom(50000, 10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		world.ProcessWorld()
	}
}

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
		fmt.Sprintf("hello")
	}
}

func BenchmarkFibProcessWorld1(b *testing.B) {
	// run the Fib function b.N times

	world := CreateWorld()

	for n := 0; n < b.N; n++ {
		world.ProcessWorld()
	}
}

func TestWorld_ProcessWorld(t *testing.T) {
	tests := []struct {
		name  string
		world *World
		want  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.world.ProcessWorld(); got != tt.want {
				t.Errorf("World.ProcessWorld() = %v, want %v", got, tt.want)
			}
		})
	}
}
