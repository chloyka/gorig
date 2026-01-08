package keyboard

func init() {
	RegisterLayout(map[rune]rune{

		'z': 'y',
		'y': 'z',
		'Z': 'Y',
		'Y': 'Z',

		'ü': '[',
		'Ü': '{',
		'ö': ';',
		'Ö': ':',
		'ä': '\'',
		'Ä': '"',
		'ß': '-',
	})
}
