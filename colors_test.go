package easygif

import (
	"image/color"
	"reflect"
	"testing"
)

func TestLerpColor(t *testing.T) {
	var (
		black color.RGBA = color.RGBA{0x00, 0x00, 0x00, 0xFF} // #000000
		white color.RGBA = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // #FFFFFF
	)

	type args struct {
		a     color.RGBA
		b     color.RGBA
		ratio float64
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		// TODO: Add test cases.
		{name: "ratio 0", args: args{a: black, b: white, ratio: 0}, want: black},
		{name: "ratio 1", args: args{a: black, b: white, ratio: 1}, want: white},
		{
			name: "ratio 1",
			args: args{
				a:     color.RGBA{0x00, 0x00, 0x00, 0xFF},
				b:     color.RGBA{0x10, 0x10, 0x10, 0x00},
				ratio: .5,
			},
			want: color.RGBA{0x08, 0x08, 0x08, 0x7F},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LerpColor(tt.args.a, tt.args.b, tt.args.ratio); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LerpColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
