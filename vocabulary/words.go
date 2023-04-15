package vocabulary

type Language int

func GetLanguage(lang_abr string) Language {
	switch lang_abr {
	case "eng":
		return English
	case "da":
		return Danish
	case "no":
		return Norwegian
	}
	return English
}

const (
	English Language = iota
	Danish
	Norwegian
)

var Suffixes = map[Language][]string{
	English:   en_suffixes,
	Danish:    {},
	Norwegian: {},
}

var SkipLists = map[Language][]string{
	English:   en_skiplist,
	Danish:    da_skiplist,
	Norwegian: {},
}

var en_suffixes = []string{
	"'s",
	"’s",
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
	"new",
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
}
