package algorithm

// Search is a wrapper around the 'search' function.
// It converts string inputs to rune slices and passess them forward.
func Search(text, pattern string) []int {
	if len(pattern) > len(text) || len(pattern) == 0 || len(text) == 0 {
		return []int{}
	}

	t := make([]rune, 0, len(text))
	for _, ch := range text {
		t = append(t, ch)
	}
	p := make([]rune, 0, len(pattern))
	for _, ch := range pattern {
		p = append(p, ch)
	}

	return search(t, p)
}

// Returns a slice of start indicies where the pattern was found.
// It uses the Boyer-Moore-Horspool algorithm,
// which doesn't use much memory and computationaly more efficient
// than other algorithms like Knuth-Moris-Pratt or Rabin-Karp.
func search(text, pattern []rune) []int {
	n := len(text)
	m := len(pattern)

	found := make([]int, 0)
	shift := make(map[rune]int)
	for i, ch := range pattern[:m-1] {
		shift[ch] = m - i - 1
	}

	i := m - 1
	for i < n {
		j := 0
		for j < m && pattern[m-j-1] == text[i-j] {
			j++
		}
		if j >= m {
			found = append(found, i-m+1)
			i++
			continue
		}

		if v, ok := shift[text[i]]; ok {
			i += v
		} else {
			i += m
		}
	}

	return found
}
