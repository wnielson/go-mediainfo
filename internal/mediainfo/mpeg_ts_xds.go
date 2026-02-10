package mediainfo

import (
	"strings"
	"unicode"
)

func normalizeXDSTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	lastSpace := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c > 0x7E {
			continue
		}
		if c == ' ' || c == '\t' {
			if b.Len() == 0 || lastSpace {
				continue
			}
			b.WriteByte(' ')
			lastSpace = true
			continue
		}
		b.WriteByte(c)
		lastSpace = false
	}
	out := strings.TrimSpace(b.String())
	if len(out) < 4 {
		return ""
	}
	letters := 0
	for i := 0; i < len(out); i++ {
		c := out[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			letters++
		}
	}
	if letters < 4 || letters*2 < len(out) {
		return ""
	}
	return out
}

func bestXDSTitle(counts map[string]int) string {
	best := ""
	bestScore := 0
	for k, v := range counts {
		score := scoreXDSTitle(k, v)
		if score > bestScore {
			best = k
			bestScore = score
		}
	}
	return best
}

func scoreXDSTitle(s string, count int) int {
	words := strings.Fields(s)
	if len(words) == 0 {
		return 0
	}
	good := 0
	bad := 0
	for _, w := range words {
		w = strings.Trim(w, " \t\r\n\"'()[]{}.,;:!?")
		if w == "" {
			continue
		}
		if isSmallTitleWord(w) {
			if w == strings.ToLower(w) {
				good++
			} else {
				bad++
			}
			continue
		}
		if isTitleWord(w) {
			good++
		} else {
			bad++
		}
	}
	// Favor plausibility over repetition (XDS packets are often noisy/partial).
	score := count*200 + len(words)*20 + good*80 - bad*120
	if bad == 0 {
		score += 40
	}
	// Bonus for common title openers.
	first := strings.Trim(words[0], "\"'()[]{}.,;:!?")
	switch strings.ToLower(first) {
	case "a", "an", "my", "the":
		score += 30
	}
	return score
}

func isSmallTitleWord(w string) bool {
	switch strings.ToLower(w) {
	case "a", "an", "and", "as", "at", "by", "for", "from", "in", "of", "on", "or", "the", "to", "with":
		return true
	default:
		return false
	}
}

func isTitleWord(w string) bool {
	// Accept patterns like "Babysitter's" or "Spider-Man" with TitleCase.
	seenLetter := false
	for i, r := range w {
		if r == '\'' || r == '-' {
			continue
		}
		if !unicode.IsLetter(r) {
			continue
		}
		if !seenLetter {
			seenLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
			if i+1 >= len(w) {
				return true
			}
			continue
		}
		// No internal capitals.
		if unicode.IsUpper(r) {
			return false
		}
	}
	return seenLetter
}
