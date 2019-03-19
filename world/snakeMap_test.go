package world

import "testing"

func TestSnakePart_AttachToBack(t *testing.T) {
	type args struct {
		partToAttach *SnakePart
	}
	tests := []struct {
		name string
		sp   *SnakePart
		args args
	}{
		{"Attach To Back", &SnakePart{}, args{&SnakePart{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sp.AttachToBack(tt.args.partToAttach)
			if tt.sp.snakePartBehind == nil {
				t.Errorf("Error SnakePart not attached to back")
			}

			if tt.args.partToAttach.snakePartInFront == nil {
				t.Errorf("Error Attached SnakePart does not have SnakePart attached to Front")
			}

		})
	}
}
