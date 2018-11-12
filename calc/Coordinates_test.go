package calc

import (
	"reflect"
	"testing"
)

func TestStringToCoordinates(t *testing.T) {
	type args struct {
		coordString string
	}
	tests := []struct {
		name string
		args args
		want Coordinates
	}{
		{"Hey", args{"2,5"}, Coordinates{X: 2, Y: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToCoordinates(tt.args.coordString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}
