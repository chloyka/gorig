package keyboard

var layoutMappings = make(map[rune]rune)

func RegisterLayout(mapping map[rune]rune) {
	for from, to := range mapping {
		layoutMappings[from] = to
	}
}

func Normalize(key string) string {
	if len(key) == 0 {
		return key
	}

	runes := []rune(key)
	if len(runes) == 1 {
		if latin, ok := layoutMappings[runes[0]]; ok {
			return string(latin)
		}
	}
	return key
}

func NormalizeRune(r rune) rune {
	if latin, ok := layoutMappings[r]; ok {
		return latin
	}
	return r
}

func IsRegistered(r rune) bool {
	_, ok := layoutMappings[r]
	return ok
}
