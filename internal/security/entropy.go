package security

import (
	"math"
	"strings"
)

func shannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	freq := make(map[rune]float64)
	for _, c := range s {
		freq[c]++
	}

	length := float64(len([]rune(s)))
	entropy := 0.0
	for _, count := range freq {
		p := count / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func isHighEntropySecret(s string, threshold float64) bool {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 8 {
		return false
	}
	return shannonEntropy(trimmed) >= threshold
}

func extractSecretValue(line string, patterns []string) string {
	for _, p := range patterns {
		idx := strings.Index(line, p)
		if idx == -1 {
			continue
		}
		rest := line[idx+len(p):]
		rest = strings.TrimSpace(rest)

		if len(rest) == 0 {
			continue
		}

		if rest[0] == '"' {
			end := strings.Index(rest[1:], "\"")
			if end != -1 {
				return rest[1 : end+1]
			}
		}
		if rest[0] == '\'' {
			end := strings.Index(rest[1:], "'")
			if end != -1 {
				return rest[1 : end+1]
			}
		}

		fields := strings.Fields(rest)
		if len(fields) > 0 {
			return fields[0]
		}
	}
	return ""
}
