package vocabulary

type Language int

func GetLanguage(lang_abr string) Language {
	switch lang_abr {
	case "eng":
		return English
	case "dan":
		return Danish
	case "nob":
		return Norwegian
	case "kor":
		return Korean
	}
	return English
}

const (
	English Language = iota
	Danish
	Norwegian
	Korean
)

var NameEnglish = map[Language]string{
	English:   "English",
	Danish:    "Danish",
	Norwegian: "Norwegian",
	Korean:    "Korean",
}

var NameNative = map[Language]string{
	English:   "English",
	Danish:    "Dansk",
	Norwegian: "Norsk",
	Korean:    "하국어",
}

var Suffixes = map[Language][]string{
	English:   en_suffixes,
	Danish:    da_suffixes,
	Norwegian: no_suffixes,
	Korean:    kor_suffixes,
}

var SkipLists = map[Language][]string{
	English:   en_skiplist,
	Danish:    da_skiplist,
	Norwegian: {},
	Korean:    {},
}

var MinWordLength = map[Language]int{
	English:   3,
	Danish:    3,
	Norwegian: 3,
	Korean:    2,
}

var en_suffixes = []string{
	"'s",
	"’s",
}

var da_suffixes = []string{
	"ene",
	"erne",
}

var no_suffixes = []string{
	"ene",
}

var kor_suffixes = []string{
	"은",
	"는",
	"를",
	"을",
	"들",
	"에",
	"에서",
	"부터",
	"까지",
	"의",
}

var en_skiplist = []string{
	"monday",
	"tuesday",
	"wednesday",
	"thursday",
	"friday",
	"saturday",
	"sunday",
	"the",
	"and",
	"for",
	"with",
	"that",
	"has",
	"from",
	"after",
	"his",
	"are",
	"was",
	"its",
	"her",
	"you",
	"your",
	"their",
	"who",
	"what",
	"over",
	"this",
	"about",
	"have",
	"will",
}

var da_skiplist = []string{
	"mandag",
	"tirsdag",
	"onsdag",
	"torsdag",
	"fredag",
	"lørdag",
	"søndag",
	"med",
	"til",
	"fra",
}
