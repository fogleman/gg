package gg

import (
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		in         string
		r, g, b, a int
	}{
		{
			in: "garbage",
			r:  0, g: 0, b: 0, a: 255,
		},

		{
			in: "#000",
			r:  0, g: 0, b: 0, a: 255,
		}, {
			in: "#000000",
			r:  0, g: 0, b: 0, a: 255,
		}, {
			in: "#00000000",
			r:  0, g: 0, b: 0, a: 0,
		},

		{
			in: "#111",
			r:  0x11, g: 0x11, b: 0x11, a: 255,
		}, {
			in: "#111111",
			r:  0x11, g: 0x11, b: 0x11, a: 255,
		}, {
			in: "#11111111",
			r:  0x11, g: 0x11, b: 0x11, a: 0x11,
		},

		{
			in: "#fff",
			r:  0xff, g: 0xff, b: 0xff, a: 0xff,
		}, {
			in: "#ffffff",
			r:  0xff, g: 0xff, b: 0xff, a: 0xff,
		}, {
			in: "#ffffffff",
			r:  0xff, g: 0xff, b: 0xff, a: 0xff,
		},
	}

	for _, tt := range tests {
		r, g, b, a := parseHexColor(tt.in)
		if tt.r != r || tt.g != g || tt.b != b || tt.a != a {
			t.Errorf("parseHexColor(%s) failed\nwant: %d %d %d %d\n got: %d %d %d %d",
				tt.in,
				tt.r, tt.g, tt.b, tt.a,
				r, g, b, a)
		}
	}
}
