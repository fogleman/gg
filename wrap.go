package gg

import (
	"strings"
	"unicode"
)

type measureStringer interface {
	MeasureString(s string) (w, h float64)
}

func splitOnSpace(x string) []string {
	var result []string
	pi := 0
	ps := false
	for i, c := range x {
		s := unicode.IsSpace(c)
		if s != ps && i > 0 {
			result = append(result, x[pi:i])
			pi = i
		}
		ps = s
	}
	result = append(result, x[pi:])
	return result
}

func wordWrap(m measureStringer, s string, width float64) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		fields := splitOnSpace(line)
		widths := make([]float64, len(fields))
		for i, field := range fields {
			widths[i], _ = m.MeasureString(field)
		}
		start := 0
		total := 0.0
		for i := 0; i < len(fields); i += 2 {
			if total+widths[i] > width {
				if i == start {
					end := i + 2
					result = append(result, strings.Join(fields[start:end], ""))
					start, total = end, 0
					continue
				} else {
					end := i
					result = append(result, strings.Join(fields[start:end], ""))
					start, total = end, 0
				}
			}
			total += widths[i] + widths[i+1]
		}
		if start < len(fields) {
			result = append(result, strings.Join(fields[start:], ""))
		}
	}
	for i, line := range result {
		result[i] = strings.TrimSpace(line)
	}
	return result
}
