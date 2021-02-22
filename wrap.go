package gg

import (
	"strings"
)

type measureStringer interface {
	MeasureString(s string) (w, h float64)
}

func wordWrap(m measureStringer, s string, width float64) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		x := ""
		for _, c := range line {
			w, _ := m.MeasureString(x + string(c))
			if w > width {
				if x == "" {
					result = append(result, string(c))
					x = ""
					continue
				} else {
					result = append(result, x)
					x = ""
				}
			}
			x += string(c)
		}
		if x != "" {
			result = append(result, x)
		}
	}
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}
	return result
}
