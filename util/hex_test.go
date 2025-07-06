package util

import "testing"

func TestItoh(t *testing.T) {
	type args[I Intish] struct {
		i I
	}
	type testCase[I Intish] struct {
		name string
		args args[I]
		want string
	}
	tests := []testCase[int]{
		{
			name: "simple",
			args: args[int]{100},
			want: "64",
		},
		{
			name: "simple",
			args: args[int]{1000},
			want: "3e8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Itoh(tt.args.i); got != tt.want {
				t.Errorf("Itoh() = %v, want %v", got, tt.want)
			}
		})
	}
}
