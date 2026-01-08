package keyboard

func init() {
	RegisterLayout(map[rune]rune{

		'a': 'q',
		'z': 'w',

		'q': 'a',

		'm': ';',

		'w': 'z',

		',': 'm',
		';': ',',
		':': '.',
		'!': '/',

		'A': 'Q',
		'Z': 'W',
		'Q': 'A',
		'M': ':',
		'W': 'Z',
		'?': 'M',

		'é': '2',
		'è': '7',
		'ç': '9',
		'à': '0',
		'ù': '\'',
	})
}
